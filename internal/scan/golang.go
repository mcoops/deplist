package scan

import (
	"os"
	"path/filepath"
	"sort"

	"golang.org/x/tools/go/packages"
)

const defaultOptions = packages.NeedDeps |
	packages.NeedImports |
	packages.NeedModule |
	packages.NeedFiles |
	packages.NeedName

func GetGolangDeps(path string) ([]*packages.Package, error) {
	dirPath := filepath.Dir(path) // force directory

	cfg := packages.Config{Mode: defaultOptions, Dir: dirPath}

	pkgs, err := packages.Load(&cfg, "./...")

	if err != nil {
		cfg.Env = append(os.Environ(), "GOFLAGS=-mod=vendor")
		pkgs, err = packages.Load(&cfg, "./...")
		if err != nil {
			return nil, err
		}
	}

	// based off the original https://github.com/golang/tools/blob/e140590b16906206021525faa5a48c7314806569/go/packages/gopackages/main.go#L99
	// todo: should just get this put into the go/packages repo instead
	var all []*packages.Package // postorder
	seen := make(map[*packages.Package]bool)

	var visit func(*packages.Package)
	visit = func(lpkg *packages.Package) {
		if !seen[lpkg] {
			seen[lpkg] = true

			// visit imports
			var importPaths []string
			for path := range lpkg.Imports {
				importPaths = append(importPaths, path)
			}
			sort.Strings(importPaths) // for determinism
			for _, path := range importPaths {
				visit(lpkg.Imports[path])
			}

			all = append(all, lpkg)
		}
	}

	for _, pkg := range pkgs {
		visit(pkg)
	}
	pkgs = all

	return pkgs, nil
}
