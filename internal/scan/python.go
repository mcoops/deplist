package scan

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Account for >, <, >=, <=, ==, !=, ~= and *
var /* const */ re = regexp.MustCompile(`[<>!~*]+`)

// GetPythonDeps scans path for python deps using the `requirements.txt` file
func GetPythonDeps(path string) (map[string]string, error) {
	log.Debugf("GetPythonDeps %s", path)
	gathered := make(map[string]string)

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// skip comments
		// and editable? https://github.com/pypa/pip/issues/4812
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "-e") || line == "" {
			continue
		}

		// easy case, elasticsearch-curator==5.8.1
		// record name and version, only for ==
		idx := strings.Index(line, "==")
		if idx > 0 {
			// test if there's a ';', i.e. unittest2==0.5.1; python_version == '2.6'
			colIdx := strings.Index(line, ";")
			if colIdx > 0 {
				if idx+2 >= colIdx {
					continue
				}
				// truncate line
				line = line[:colIdx]
			}
			gathered[line[:idx]] = line[idx+2:]
			continue
		}

		// every other permitation just use the name as we can't guarantee
		// the version, just grab the name using first occurrence
		match := re.FindStringIndex(line)

		if match != nil {
			gathered[line[:match[0]]] = ""
		}
	}

	return gathered, nil
}
