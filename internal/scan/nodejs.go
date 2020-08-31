package scan

import (
	"os/exec"
	"path/filepath"
	"strings"
)

type nodejs struct {
	path    string
	version string
}

func GetNodeJSDeps(path string) ([]string, error) {
	dirPath := filepath.Dir(path)

	cmd := exec.Command("npm", "ls", "--prod", "--parseable", "--long", "--silent")
	cmd.Dir = dirPath

	results, err := cmd.Output()

	if err != nil {

	}

	var out []string

	res := strings.Split(string(results), "\n")

	if len(res) <= 1 {

	}

	for _, s := range res {
		idx := strings.Index(s, "node_modules")

		var formatted string

		if idx == -1 {
			formatted = filepath.Base(s)
		} else {
			formatted = s[idx+len("node_modules/"):]
		}

		formatted = strings.TrimRight(formatted, ":undefined")

		if len(formatted) > 2 {
			dupIdx := strings.Index(formatted, ":")
			formatted = formatted[dupIdx+1:]

			out = append(out, formatted)
		}
	}
	return out, nil
}
