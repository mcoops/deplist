package deplist

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/RedHatProductSecurity/deplist/internal/scan"
	"github.com/RedHatProductSecurity/deplist/internal/utils"
)

// enums start at 1 to allow us to specify found languages 0 = nil
const (
	LangGolang = 1 << iota
	LangJava
	LangNodeJS
	LangPython
	LangRuby
)

func init() {
	// check for the library required binaries
	if _, err := exec.LookPath("yarn"); err != nil {
		log.Fatal("yarn is required in PATH")
	}

	if _, err := exec.LookPath("npm"); err != nil {
		log.Fatal("npm is required in PATH")
	}

	if _, err := exec.LookPath("go"); err != nil {
		log.Fatal("go is required")
	}

	if _, err := exec.LookPath("mvn"); err != nil {
		log.Fatal("maven is required")
	}

	if _, err := exec.LookPath("bundle"); err != nil {
		log.Fatal("bundler gem is required")
	}
}

// GetLanguageStr returns from a bitmask return the ecosystem name
func GetLanguageStr(bm Bitmask) string {
	if bm&LangGolang != 0 {
		return "go"
	} else if bm&LangJava != 0 {
		return "mvn"
	} else if bm&LangNodeJS != 0 {
		return "npm"
	} else if bm&LangPython != 0 {
		return "pypi"
	} else if bm&LangRuby != 0 {
		return "gem"
	}
	return "unknown"
}

// GetDeps scans a given repository and returns all dependencies found in a DependencyList struct.
func GetDeps(fullPath string) ([]Dependency, Bitmask, error) {
	// var deps DependencyList
	var deps []Dependency
	var foundTypes Bitmask = 0

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, 0, os.ErrNotExist
	}

	pomPath := filepath.Join(fullPath, "pom.xml")
	goPath := filepath.Join(fullPath, "go.mod")
	goPkgPath := filepath.Join(fullPath, "Gopkg.lock")
	glidePath := filepath.Join(fullPath, "glide.lock")
	rubyPath := filepath.Join(fullPath, "Gemfile.lock")
	pythonPath := filepath.Join(fullPath, "requirements.txt")

	// point at the parent repo, but can't assume where the indicators will be
	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// prevent panic by handling failure https://golang.org/pkg/path/filepath/#Walk
			return err
		}

		if info.IsDir() {
			// prevent walking down the vendors, docs, etc
			if utils.BelongsToIgnoreList(info.Name()) {
				return filepath.SkipDir
			}
		} else {
			// Two checks, one for filenames and the second switch for full
			// paths. Useful if we're looking for top of repo

			switch filename := info.Name(); filename {
			// for now only go for yarn and npm
			case "package-lock.json":
				// if theres not a yarn.lock fall thru
				if _, err := os.Stat(
					filepath.Join(
						filepath.Dir(path),
						"yarn.lock")); err == nil {
					return nil
				}
				fallthrough

			case "yarn.lock":
				pkgs, err := scan.GetNodeJSDeps(path)
				if err != nil {
					return err
				}

				if len(pkgs) > 0 {
					foundTypes.DepFoundAddFlag(LangNodeJS)
				}

				for _, p := range pkgs {
					deps = append(deps,
						Dependency{
							DepType: LangNodeJS,
							Path:    p.Name,
							Version: p.Version,
							Files:   []string{},
						})
				}
			default:

				ext := filepath.Ext(filename)

				// java
				if ext == ".jar" || ext == ".war" || ext == ".ear" || ext == ".adm" || ext == ".hpi" || ext == ".zip" {
					file := strings.Replace(filepath.Base(path), ext, "", 1) // get filename, check if we can ignore
					if strings.HasSuffix(file, "-sources") || strings.HasSuffix(file, "-javadoc") {
						return nil
					}
					pkgs, err := scan.GetJarDeps(path)
					if err == nil {

						if len(pkgs) > 0 {
							foundTypes.DepFoundAddFlag(LangJava)
						}

						for name, version := range pkgs {
							// just in case we report the full path to the dep
							name = strings.Replace(name, fullPath, "", 1)

							// if the dep ends with -javadoc or -sources, not really interested
							if !strings.HasSuffix(version, "-javadoc") && !strings.HasSuffix(version, "-sources") {
								deps = append(deps,
									Dependency{
										DepType: LangJava,
										Path:    name,
										Version: version,
										Files:   []string{},
									})
							}
						}
					}
				}
			}

			switch path {
			case goPath: // just support the top level go.mod for now
				pkgs, err := scan.GetGolangDeps(path)
				if err != nil {
					return err
				}

				if len(pkgs) > 0 {
					foundTypes.DepFoundAddFlag(LangGolang)
				}

				for path, goPkg := range pkgs {
					d := Dependency{
						DepType: LangGolang,
						Path:    path,
						Files:   goPkg.Gofiles,
						Version: goPkg.Version,
					}
					deps = append(deps, d)
				}
			case goPkgPath:
				pkgs, err := scan.GetGoPkgDeps(path)
				if err != nil {
					return err
				}

				if len(pkgs) > 0 {
					foundTypes.DepFoundAddFlag(LangGolang)
				}
				for _, goPkg := range pkgs {
					d := Dependency{
						DepType: LangGolang,
						Path:    goPkg.Name,
						Version: goPkg.Version,
					}
					deps = append(deps, d)
				}
			case glidePath:
				pkgs, err := scan.GetGlideDeps(path)
				if err != nil {
					return err
				}

				if len(pkgs) > 0 {
					foundTypes.DepFoundAddFlag(LangGolang)
				}
				for _, goPkg := range pkgs {
					d := Dependency{
						DepType: LangGolang,
						Path:    goPkg.Name,
						Version: goPkg.Version,
					}
					deps = append(deps, d)
				}
			case pomPath:
				pkgs, err := scan.GetMvnDeps(path)
				if err != nil {
					return err
				}

				if len(pkgs) > 0 {
					foundTypes.DepFoundAddFlag(LangJava)
				}

				for name, version := range pkgs {
					deps = append(deps,
						Dependency{
							DepType: LangJava,
							Path:    name,
							Version: strings.Replace(version, "v", "", 1),
							Files:   []string{},
						})
				}
			case rubyPath:
				pkgs, err := scan.GetRubyDeps(path)
				if err != nil {
					return err
				}

				if len(pkgs) > 0 {
					foundTypes.DepFoundAddFlag(LangRuby)
				}

				for name, version := range pkgs {
					deps = append(deps,
						Dependency{
							DepType: LangRuby,
							Path:    strings.TrimSuffix(name, "\n"),
							Version: strings.Replace(version, "v", "", 1),
							Files:   []string{},
						})
				}
			case pythonPath:
				pkgs, err := scan.GetPythonDeps(path)
				if err != nil {
					return err
				}

				if len(pkgs) > 0 {
					foundTypes.DepFoundAddFlag(LangPython)
				}

				for name, version := range pkgs {
					deps = append(deps,
						Dependency{
							DepType: LangPython,
							Path:    name,
							Version: version,
							Files:   []string{},
						})
				}
			}

		}
		return nil
	})

	if err != nil {
		return nil, 0, err // should't matter
	}

	return deps, foundTypes, nil
}
