package deplist

import (
	"testing"
)

func BuildWant() []Dependency {
	var deps []Dependency

	golangPaths := []string{
		"internal/cpu",
		"unsafe",
		"internal/bytealg",
		"runtime/internal/atomic",
		"runtime/internal/sys",
		"runtime/internal/math",
		"runtime",
		"internal/reflectlite",
		"errors",
		"math/bits",
		"math",
		"unicode/utf8",
		"strconv",
		"internal/race",
		"sync/atomic",
		"sync",
		"unicode",
		"reflect",
		"sort",
		"internal/fmtsort",
		"io",
		"internal/oserror",
		"syscall",
		"time",
		"internal/poll",
		"internal/syscall/unix",
		"internal/syscall/execenv",
		"internal/testlog",
		"os",
		"fmt",
		"github.com/openshift/api/config/v1",
		"golang.org/x/text/unicode",
		"github.com/mcoops/deplist",
	}

	npmSet1 := []string{
		"angular",
		"d3",
		"d3-array",
		"d3-axis",
		"d3-brush",
		"d3-chord",
		"d3-color",
		"d3-contour",
		"d3-delaunay",
		"delaunator",
		"d3-dispatch",
		"d3-drag",
		"d3-dsv",
		"commander",
		"iconv-lite",
	}

	npmSet2 := []string{
		"safer-buffer",
		"rw",
		"d3-ease",
		"d3-fetch",
		"d3-force",
		"d3-format",
		"d3-geo",
		"d3-hierarchy",
		"d3-interpolate",
		"d3-path",
		"d3-polygon",
		"d3-quadtree",
		"d3-random",
		"d3-scale",
		"d3-scale-chromatic",
		"d3-selection",
		"d3-shape",
		"d3-time",
		"d3-time-format",
		"d3-timer",
		"d3-transition",
		"d3-zoom",
		"prismjs",
		"clipboard",
		"good-listener",
		"delegate",
		"select",
		"tiny-emitter",
		"react",
		"loose-envify",
		"js-tokens",
		"object-assign",
		"prop-types",
		"react-is",
		"rxjs",
		"tslib",
		"slate",
		"@types/esrever",
		"esrever",
		"immer",
		"is-plain-object",
		"tiny-warning",
		"tether",
	}

	for _, n := range golangPaths {
		d := Dependency{
			DepType: 1,
			Path:    n,
		}

		deps = append(deps, d)
	}

	end := len(deps) - 2 // get the unicode ver
	deps[end].Version = "v0.3.3"

	for _, n := range npmSet1 {
		d := Dependency{
			DepType: 3,
			Path:    n,
		}
		deps = append(deps, d)
	}
	deps = append(deps, Dependency{DepType: 2, Path: "com.amazonaws:aws-lambda-java-core"}) // java

	for _, n := range npmSet2 {
		d := Dependency{
			DepType: 3,
			Path:    n,
		}
		deps = append(deps, d)
	}
	return deps
}

func TestGetDeps(t *testing.T) {
	want := BuildWant()

	got, gotBitmask, _ := GetDeps("test/testRepo")

	if gotBitmask != 7 {
		t.Errorf("GotBitmask() != 7")
	}

	// iterate thru and compare
	if len(want) != len(got) {
		t.Errorf("GetDeps() = %d; want %d", len(got), len(want))
	}

	for _, pkg := range want {
		// because the maps are random...
		flag := false
		for _, g := range got {
			if pkg.Path == g.Path {
				if pkg.Version == "" || pkg.Version == g.Version {
					flag = true
					break
				}
			}
		}
		if !flag {
			t.Errorf("GetDeps() wanted: %s - not found", pkg.Path)
		}
	}
}
