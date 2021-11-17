package scan

import (
	"io/ioutil"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type GlideDeps struct {
	Name    string
	Version string
}

type glideDep struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	Subpackages []string `yaml:"subpackages"`
}

type glideDeps struct {
	Imports []glideDep
}

func GetGlideDeps(path string) ([]GlideDeps, error) {
	log.Debugf("GetGlideDeps %s", path)

	var deps glideDeps
	var gathered []GlideDeps

	yaml_file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yaml_file, &deps)

	if err != nil {
		return nil, err
	}

	for _, d := range deps.Imports {
		gathered = append(gathered,
			GlideDeps{
				Name:    d.Name,
				Version: d.Version,
			},
		)

		for _, subpackage := range d.Subpackages {
			gathered = append(gathered,
				GlideDeps{
					Name:    filepath.Join(d.Name, subpackage),
					Version: d.Version,
				},
			)
		}
	}

	return gathered, nil
}
