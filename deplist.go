package deplist

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mcoops/deplist/internal/scan"
	"github.com/mcoops/deplist/internal/utils"
)

// enums start at 1 to allow us to specify found languages 0 = nil
const (
	langGolang = 1 << iota
	langNodeJS
	langPython
	langRuby
)

// GetDeps scans a given repository and returns all dependencies found in a DependencyList struct.
func GetDeps(fullPath string) ([]Dependency, Bitmask, error) {
	// var deps DependencyList
	var deps []Dependency
	var foundTypes Bitmask = 0

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, 0, os.ErrNotExist
	}

	// point at the parent repo, but can't assume where the indicators will be
	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// prevent walking down the vendors, docs, etc
			if utils.BelongsToIgnoreList(info.Name()) {
				return filepath.SkipDir
			}
		} else {
			switch filename := info.Name(); filename {
			case "go.mod": // just support the top level go.mod for now
				pkgs, err := scan.GetGolangDeps(path)
				if err != nil {
					return err
				}

				foundTypes.DepFoundAddFlag(langGolang)

				for _, p := range pkgs {
					d := Dependency{
						depType: langGolang,
						path:    p.PkgPath,
						files:   p.GoFiles,
					}
					if p.Module != nil {
						d.version = p.Module.Version
					}
					deps = append(deps, d)
				}

			case "package.json":
				pkgs, err := scan.GetNodeJSDeps(path)
				if err != nil {
					return err
				}

				foundTypes.DepFoundAddFlag(langNodeJS)

				for _, p := range pkgs {
					splitIdx := strings.LastIndex(p, "@")

					d := Dependency{
						depType: langNodeJS,
						path:    p[:splitIdx],
						version: p[splitIdx+1:],
					}
					deps = append(deps, d)
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
