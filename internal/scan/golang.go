package scan

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

type GoListDeps struct {
	ImportPath string `json:"ImportPath"`
	Module     struct {
		Version string `json:"Version"`
		Replace struct {
			Version string `json:"Version"`
		} `json:"Replace"`
	} `json:"Module"`
	GoFiles []string `json:"GoFiles"`
}

type GoPkg struct {
	Version string
	Gofiles []string
}

func getVersion(deps GoListDeps) string {
	/* if replace is specified, then use that version
	* not seen when version and replace.version are differnt
	* but just in case
	 */
	if deps.Module.Replace.Version != "" {
		// due to the way we're loading the json this time, this just works
		return deps.Module.Replace.Version
	}
	return deps.Module.Version
}

func runCmd(path string, mod bool) ([]byte, error) {
	// go list -f '{{if not .Standard}}{{.Module}}{{end}}' -json -deps ./...
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// go list -f '{{if not .Standard}}{{.Module}}{{end}}' -json -deps ./...
	var cmd *exec.Cmd

	if !mod {
		cmd = exec.CommandContext(ctx, "go", "list", "-f", "'{{if not .Standard}}{{.Module}}{{end}}'", "-json", "-deps", "./...")
	} else {
		vendorDir := filepath.Join(filepath.Dir(path), "vendor")
		if _, err := os.Stat(vendorDir); err != nil {
			if os.IsNotExist(err) {
				return nil, errors.New("no 'vendor' directory, can't use '-mod=vendor'")
			}
		}
		cmd = exec.CommandContext(ctx, "go", "list", "-mod=vendor", "-f", "'{{if not .Standard}}{{.Module}}{{end}}'", "-json", "-deps", "./...")
	}

	cmd.Dir = filepath.Dir(path) // // force directory
	out, err := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		return nil, ctx.Err()
	}

	// mod=vendor sometimes still returns results but returns an error. In
	// that case ignore the error and return what we can
	if err != nil {
		if !mod {
			// assume some retrival error, we have to redo the cmd with mod=vendor
			return nil, err
		} else {
			if len(out) == 0 {
				return nil, err
			}
		}
	}

	return out, nil
}

/*
* Need to support re-running the go list with and without -mod=vendor
* First run defaults to without, if any kind of error we'll just retry the run
 */
func runGoList(path string) ([]byte, error) {
	out, err := runCmd(path, false)
	if err != nil {
		// rerun
		out, err = runCmd(path, true)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func GetGolangDeps(path string) (map[string]GoPkg, error) {
	// need to use a map as we'll get lots of duplicate entries
	gathered := make(map[string]GoPkg)

	out, err := runGoList(path)

	if err != nil {
		return nil, err
	}

	/* we cann't just marshall the json as go list returns multiple json
	 * documents not an array of json - which is annoying
	 */
	decoder := json.NewDecoder(strings.NewReader(string(out)))

	for {
		var goListDeps GoListDeps
		err := decoder.Decode(&goListDeps)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		importPath := goListDeps.ImportPath

		if _, ok := gathered[importPath]; ok {
			if gathered[importPath].Version != semver.Max(gathered[importPath].Version, getVersion(goListDeps)) {
				gathered[importPath] = GoPkg{
					Version: getVersion(goListDeps),
					Gofiles: goListDeps.GoFiles,
				}
			}
		} else {
			gathered[importPath] = GoPkg{
				Version: getVersion(goListDeps),
				Gofiles: goListDeps.GoFiles,
			}
		}
	}
	return gathered, nil
}
