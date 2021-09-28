package jargo

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/vifraa/gopom"
)

type Manifest map[string]string

type JarInfo struct {
	*Manifest
	Files []string
}

var JarFilter = [...]string{".jar", ".war", ".hpi", ".rar", ".ear", ".zip"}

var jarVersionRegex *regexp.Regexp

const MANIFEST_FULL_NAME = "META-INF/MANIFEST.MF"

func init() {
	jarVersionRegex = regexp.MustCompile(".*?-([0-9\\.][\\w0-9\\.-]*)\\.jar")
}

// https://codereview.stackexchange.com/questions/191238/return-unique-items-in-a-go-slice/192954#192954
func unique(slice []string) []string {
	// create a map with all the values as key
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	// turn the map keys into a slice
	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}

// GetManifest extracts the manifest info from a Java JAR file
// It takes as parameter the path to the jar file of interest
// It returns a pointer to a Manifest (map[string]string) which is the key:values pairs from the META-INF/MANIFEST.MF file
func GetManifest(filename string) (*Manifest, error) {
	jar, err := readFromFile(filename, false)

	if err != nil {
		return nil, err
	}
	return jar.Manifest, nil
}

// GetJarInfo extracts various info from a Java JAR file
// It takes as parameter the path to the jar file of interest
// It extracts the Manifest (like GetManifest)
// It extracts an array of the filenames in the JAR file
// It returns a pointer to a JarInfo struct
func GetJarInfo(filename string) (*JarInfo, error) {
	jar, err := readFromFile(filename, true)

	if err != nil {
		return nil, errors.New("Error processing " + filename + " " + err.Error())
	}

	first := jar.Files[0]
	// not great, but given it's recursive we're generally going to get the top level jar accidentally.
	if strings.Contains(first, strings.Replace(filename, filepath.Ext(filename), "", 1)) {
		jar.Files = jar.Files[1:]
	}

	// unique it
	jar.Files = unique(jar.Files)

	return jar, err
}

var gJar JarInfo

func readFromFile(filename string, fullJar bool) (*JarInfo, error) {
	var err error
	var file *os.File
	var fi os.FileInfo
	var r *zip.Reader

	if file, err = os.Open(filename); err != nil {
		return nil, err
	}
	defer file.Close()

	if fi, err = file.Stat(); err != nil {
		return nil, err
	}

	if r, err = zip.NewReader(file, fi.Size()); err != nil {
		return nil, err
	}
	return readFromReader(r, fullJar, filename)
}

func inExtFilter(name string) bool {
	for _, ext := range JarFilter {
		if strings.HasSuffix(name, ext) {
			return true
		}
	}
	return false
}

func cleanJarName(name string) string {
	out := strings.TrimPrefix(name, "WEB-INF/lib/")
	out = strings.TrimPrefix(out, "WEB-INF/")
	out = strings.TrimPrefix(out, "detached-plugins/")
	out = strings.Replace(out, filepath.Ext(out), "", 1)

	return out
}

func readFromReader(r *zip.Reader, fullJar bool, jarName string) (*JarInfo, error) {
	var (
		part   []byte
		prefix bool
		lines  []string
	)

	jar := new(JarInfo)
	if fullJar {
		lines = make([]string, 0)
	}
	lineNumber := -1
	for _, f := range r.File {
		if fullJar {

			if inExtFilter(f.Name) {
				// jar.Files = append(jar.Files, f.Name)
				rf, err := f.Open()
				if err != nil { // just in case
					continue
				}

				body, _ := ioutil.ReadAll(rf)
				zr, err := zip.NewReader(bytes.NewReader(body), f.FileInfo().Size())
				if err == nil { // pretty much ignore errors
					rJar, err := readFromReader(zr, true, f.Name)
					if err == nil {
						jar.Files = append(jar.Files, rJar.Files...)
					}
				}
			} else if strings.HasPrefix(f.Name, "META-INF") && strings.HasSuffix(f.Name, "pom.xml") {
				// get the version from MANIFEST.MF

				name := strings.Replace(f.Name, "META-INF/maven/", "", 1)

				rf, err := f.Open()
				if err != nil {
					continue
				}
				body, _ := ioutil.ReadAll(rf)
				var parsedPom gopom.Project

				err = xml.Unmarshal(body, &parsedPom)
				if err != nil {
					continue
				}

				// we could potentially add jarname to the path, but makes searching a pain later
				if parsedPom.Version == "" {
					jar.Files = append(jar.Files, strings.Replace(name, "/pom.xml", "@"+parsedPom.Parent.Version, 1))
				} else {
					jar.Files = append(jar.Files, strings.Replace(name, "/pom.xml", "@"+parsedPom.Version, 1))
				}
			}
		}

		if f.Name == MANIFEST_FULL_NAME {
			rc, err := f.Open()
			if err != nil {
				log.Println(err)
				return nil, err
			}
			reader := bufio.NewReader(rc)
			buffer := bytes.NewBuffer(make([]byte, 0))

			for {
				if part, prefix, err = reader.ReadLine(); err != nil {
					break
				}
				if len(part) == 0 {
					continue
				}
				buffer.Write(part)
				if !prefix {
					//lines = append(lines, buffer.String())
					line := buffer.String()
					if line[0] == ' ' {
						lines[lineNumber] = lines[lineNumber] + line
					} else {
						lines = append(lines, line)
						lineNumber = lineNumber + 1
					}
					buffer.Reset()
				}
			}
			if err == io.EOF {
				err = nil
			}
			rc.Close()
			jar.Manifest = makeManifestMap(lines)
		}
	}

	// if we're here and jar.Files is empty, means no pom.xml, so log something
	if jar.Files == nil {
		// use the jarfilename
		if jarName != "" {
			v := jarVersionRegex.FindStringSubmatch(jarName)
			n := cleanJarName(jarName)
			var ver string = ""
			if v == nil { // just go off last "-"
				idx := strings.LastIndex(n, "-")
				if idx != -1 {
					ver = n[idx+1:]
				}
			} else {
				ver = v[1]
			}

			n = strings.Replace(n, "-"+ver, "@"+ver, 1)
			jar.Files = append(jar.Files, n)
		}
	}

	return jar, nil
}

func makeManifestMap(lines []string) *Manifest {
	manifestMap := make(Manifest)

	for _, line := range lines {
		i := strings.Index(line, ":")
		if i == -1 || i == 0 {
			// log.Println("Not properties file?? This line missing colon (:): " + line)
			continue
		}
		key := strings.TrimSpace(line[0:i])
		value := strings.TrimSpace(line[i+1:])
		manifestMap[key] = value
	}
	return &manifestMap
}
