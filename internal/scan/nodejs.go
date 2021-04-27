package scan

import (
	"encoding/json"
	"fmt"
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

type npmDependency struct {
	Version      string                   `json:"version"`
	Dependencies map[string]npmDependency `json:"dependencies"`
}

type npmListOutput struct {
	Dependencies map[string]npmDependency `json:"dependencies"`
}

var gathered map[string]string

func recordPackage(packageName, version string) {
	// compare everything
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}
	version = strings.Replace(version, "^", "", 1)
	version = strings.Replace(version, "~", "", 1)

	version = strings.Replace(version, "x", "0", 1)
	version = strings.Replace(version, "*", "0.0.0", 1)

	if oldVersion, ok := gathered[packageName]; ok {
		gathered[packageName] = semver.Max(oldVersion, version)
	} else {
		gathered[packageName] = version
	}
}

func gatherYarnNode(dep yarnDependency) {
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

	recordPackage(name, version)

	if len(dep.Children) > 0 {
		for _, child := range dep.Children {
			gatherYarnNode(child)
		}
	}
}

func gatherNPMNode(name string, dependency npmDependency) {
	recordPackage(name, dependency.Version)
	for childName, childDep := range dependency.Dependencies {
		gatherNPMNode(childName, childDep)
	}
}

func GetNodeJSDeps(path string) (map[string]string, error) {
	switch filepath.Base(path) {
	case "yarn.lock":
		return getYarnDeps(path)
	case "package-lock.json":
		return getNPMDeps(path)
	}
	return nil, fmt.Errorf("unknown NodeJS dependency file %q", path)
}

func getYarnDeps(path string) (map[string]string, error) {
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
		gatherYarnNode(deps)
	}

	return gathered, nil
}

func getNPMDeps(path string) (map[string]string, error) {
	var npmOutput npmListOutput
	gathered = make(map[string]string)

	cmd := exec.Command("npm", "list", "--prod", "--json", "--depth=99")
	cmd.Dir = filepath.Dir(path)

	data, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &npmOutput)
	if err != nil {
		return nil, err
	}

	for depName, dep := range npmOutput.Dependencies {
		gatherNPMNode(depName, dep)
	}

	return gathered, nil
}
