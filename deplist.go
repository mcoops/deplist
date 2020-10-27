package deplist

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mcoops/deplist/internal/scan"
	"github.com/mcoops/deplist/internal/utils"
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

	if _, err := exec.LookPath("go"); err != nil {
		log.Fatal("go is required")
	}

	if _, err := exec.LookPath("mvn"); err != nil {
		log.Fatal("maven is required")
	}
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
			// for now only go for yarn
			case "yarn.lock":
				pkgs, err := scan.GetNodeJSDeps(path)
				if err != nil {
					return err
				}

				if len(pkgs) > 0 {
					foundTypes.DepFoundAddFlag(LangNodeJS)
				}

				for name, version := range pkgs {
					deps = append(deps,
						Dependency{
							DepType: LangNodeJS,
							Path:    name,
							Version: strings.Replace(version, "v", "", 1),
						})
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

				for _, p := range pkgs {
					d := Dependency{
						DepType: LangGolang,
						Path:    p.PkgPath,
						Files:   p.GoFiles,
					}
					if p.Module != nil {
						d.Version = p.Module.Version
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
