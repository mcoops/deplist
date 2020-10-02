package scan

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/mod/semver"
)

type yarnDependencies []yarnDependency

type yarnDependency struct {
	Name     string           `json:"name"`
	Children yarnDependencies `json:"children"`
}

type yarnOutput struct {
	Type string `json:"type"`
	Data struct {
		Deps yarnDependencies `json:"trees"`
	}
}

var gathered map[string]string

func gatherNode(dep yarnDependency) {
	// incase package starts with @
	splitIdx := strings.LastIndex(dep.Name, "@")

	var name string
	var version string

	if splitIdx != -1 {
		name = dep.Name[:splitIdx]
		version = "v" + dep.Name[splitIdx+1:]
	} else {
		name = dep.Name
		version = "v0.0.0"
	}

	// compare everything
	version = strings.Replace(version, "^", "", 1)
	version = strings.Replace(version, "~", "", 1)

	version = strings.Replace(version, "x", "0", 1)
	version = strings.Replace(version, "*", "0.0.0", 1)

	if _, ok := gathered[name]; ok {
		gathered[name] = semver.Max(gathered[name], version)
	} else {
		gathered[name] = version
	}

	if len(dep.Children) > 0 {
		for _, child := range dep.Children {
			gatherNode(child)
		}
	}
}

func GetNodeJSDeps(path string) (map[string]string, error) {

	var yarnOutput yarnOutput
	gathered = make(map[string]string)

	dirPath := filepath.Dir(path)

	data, err := exec.Command("yarn", "--cwd", dirPath, "list", "--prod", "--json", "--no-progress").Output()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &yarnOutput)
	if err != nil {
		return nil, err
	}

	for _, deps := range yarnOutput.Data.Deps {
		gatherNode(deps)
	}

	return gathered, nil
}
