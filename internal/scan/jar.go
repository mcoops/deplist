package scan

import (
	"strings"

	"github.com/mcoops/jargo"
)

// GetJarDeps uses github.com/mcoops/jargo retrieve the java dependencies
func GetJarDeps(path string) (map[string]string, error) {
	gathered := make(map[string]string)

	jar, err := jargo.GetJarInfo(path)

	if err != nil {
		return nil, err
	}

	for _, j := range jar.Files {
		idx := strings.LastIndex(j, "@")

		if idx == -1 || idx+1 > len(j) {
			gathered[j] = ""
		} else {
			gathered[j[:idx]] = j[idx+1:]
		}
	}
	return gathered, nil
}
