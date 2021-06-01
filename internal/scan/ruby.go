package scan

import (
	"os/exec"
	"path/filepath"
	"strings"
)

func GetRubyDeps(path string) (map[string]string, error) {
	gathered := make(map[string]string)

	dirPath := filepath.Dir(path)

	//Make sure that the Gemfile we are loading is supported by the version of bundle currently installed.
	cmd := exec.Command("bundle", "update", "--bundler")
	cmd.Dir = dirPath
	_, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	cmd = exec.Command("bundle", "list")

	cmd.Dir = dirPath

	data, err := cmd.Output()

	splitOutput := strings.Split(string(data), "\n")

	for _, line := range splitOutput {
		if !strings.HasPrefix(line, "  *") {
			continue
		}
		rawDep := strings.TrimPrefix(line, "  * ")
		dep := strings.Split(rawDep, " ")
		dep[1] = dep[1][1 : len(dep[1])-1]
		gathered[dep[0]] = dep[1]
	}

	if err != nil {
		return nil, err
	}

	return gathered, nil
}
