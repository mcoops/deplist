package scan

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type GoPkgLockDeps struct {
	Name    string
	Version string
	Gofiles []string
}

type goPkg struct {
	Name     string   `toml:"name"`
	Packages []string `toml:"packages"`
	Revision string   `toml:"revision"`
	Version  string   `toml:"version"`
}

type goPkgDeps struct {
	Deps []goPkg `toml:"projects"`
}

func GetGoPkgDeps(path string) ([]GoPkgLockDeps, error) {
	log.Debugf("GetGoPkgDeps %s", path)

	var deps goPkgDeps
	var gathered []GoPkgLockDeps

	_, err := toml.DecodeFile(path, &deps)
	if err != nil {
		return nil, err
	}

	for _, d := range deps.Deps {
		ver := d.Version
		if ver == "" {
			ver = d.Revision
		}

		gathered = append(gathered,
			GoPkgLockDeps{
				Name:    d.Name,
				Version: ver,
			},
		)

		for _, subpackage := range d.Packages {
			if subpackage != "." {
				gathered = append(gathered,
					GoPkgLockDeps{
						Name:    filepath.Join(d.Name, subpackage),
						Version: ver,
					},
				)
			}
		}
	}

	return gathered, nil
}
