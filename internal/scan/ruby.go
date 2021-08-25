package scan

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetRubyDeps(path string) (map[string]string, error) {
	gathered := make(map[string]string)

	dirPath := filepath.Dir(path)

	// override the gem path otherwise might hit perm issues and it's annoying
	gem_path, err := ioutil.TempDir("", "gem_vendor")
	if err != nil {
		return nil, err
	}

	// cleanup after ourselves
	defer os.RemoveAll(gem_path)

	//Make sure that the Gemfile we are loading is supported by the version of bundle currently installed.
	cmd := exec.Command("bundle", "update", "--bundler")
	cmd.Dir = dirPath
	cmd.Env = append(cmd.Env, "GEM_HOME="+gem_path)
	_, err = cmd.Output()
	if err != nil {
		return nil, err
	}

	cmd = exec.Command("bundle", "list")

	cmd.Dir = dirPath
	cmd.Env = append(cmd.Env, "GEM_HOME="+gem_path)

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
