package scan

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mcoops/deplist/internal/utils"
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

// NodeJSGather dependencies found, name and version
type NodeJSGather struct {
	Name    string
	Version string
}

// TODO: remove this global
var gatheredNode map[string]NodeJSGather

func recordPackage(packageName, version string) {
	// opposite now, we don't care if its specifying version ranges like 5.x.x,
	// or 5.* etc. Just get the versions.
	if len(version) > 0 {
		if !utils.CharIsDigit(version) {
			return
		}

		if version[len(version)-1] == 'x' {
			return
		}
	}

	if _, ok := gatheredNode[packageName+version]; !ok {
		gatheredNode[packageName+version] = NodeJSGather{
			Name:    packageName,
			Version: version,
		}
	}
}

func gatherYarnNode(dep yarnDependency) {
	// incase package starts with @
	splitIdx := strings.LastIndex(dep.Name, "@")

	var name string
	var version string

	if splitIdx != -1 {
		name = dep.Name[:splitIdx]
		version = dep.Name[splitIdx+1:]
	} else {
		name = dep.Name
		version = ""
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

// GetNodeJSDeps scans the path for either `yarn.lock` or `package-lock.json`,
// then use the appropriate pkg managers to produce depencies lists of type
// `NodeJSGather`
func GetNodeJSDeps(path string) (map[string]NodeJSGather, error) {
	switch filepath.Base(path) {
	case "yarn.lock":
		return getYarnDeps(path)
	case "package-lock.json":
		return getNPMDeps(path)
	}
	return nil, fmt.Errorf("unknown NodeJS dependency file %q", path)
}

func getYarnDeps(path string) (map[string]NodeJSGather, error) {
	var yarnOutput yarnOutput
	gatheredNode = make(map[string]NodeJSGather)

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

	return gatheredNode, nil
}

func getNPMDeps(path string) (map[string]NodeJSGather, error) {
	var npmOutput npmListOutput
	gatheredNode = make(map[string]NodeJSGather)

	cmd := exec.Command("npm", "list", "--prod", "--json", "--depth=99")
	cmd.Dir = filepath.Dir(path)

	data, err := cmd.Output()

	// npm has a nasty habbit of not returning cleanly so if there is data
	// just attempt to unmarshal
	if data == nil && err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &npmOutput)
	if err != nil {
		return nil, err
	}

	for depName, dep := range npmOutput.Dependencies {
		gatherNPMNode(depName, dep)
	}

	return gatheredNode, nil
}
