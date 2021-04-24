package scan

import (
	"os/exec"
	"path/filepath"
	"strings"
)

func GetMvnDeps(path string) (map[string]string, error) {
	var gathered map[string]string
	var found map[string]bool

	dirPath := filepath.Dir(path)

	// cmd := exec.Command("mvn",
	// 	"--no-transfer-progress",
	// 	"dependency:tree",
	// 	"-DoutputType=dot",
	// )

	// Opposed to mvn dependency:tree which fails if there's issues with
	// finding build deps dependency:collect does not fail to continue
	cmd := exec.Command(
		"mvn",
		"--no-transfer-progress",
		"dependency:collect",
		"-DincludeScope=runtime")
	cmd.Dir = dirPath

	// supress error, it always returns errors
	data, _ := cmd.Output()

	res := strings.Split(string(data), "\n")

	gathered = make(map[string]string)

	for _, s := range res {
		if len(s) < 5 && !strings.HasPrefix(s, "[INFO]") {
			continue
		}

		if !strings.HasSuffix(s, "compile") && !strings.HasSuffix(s, "runtime") {
			continue
		}

		// remove the :compile or :runtime off the end
		lastColon := strings.LastIndex(s, ":")
		if lastColon == -1 {
			continue
		}
		s = s[:lastColon]

		verIdx := strings.LastIndex(s, ":")
		if verIdx == -1 || len(s) < (verIdx+1) {
			continue
		}
		ver := s[verIdx+1:]

		name := strings.Replace(s, ":"+ver, "", 1)

		startIdx := strings.Index(name, "    ")
		if startIdx == -1 || len(name) < (startIdx+4) {
			continue
		}
		name = name[startIdx+4:]

		if _, ok := found[name+ver]; ok == false {
			gathered[name] = ver
		}
	}
	return gathered, nil
}
