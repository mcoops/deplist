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
		"internal/testlog",
		"os",
		"fmt",
		"golang.org/x/text/unicode",
		"github.com/mcoops/deplist",
	}

	nodejsPaths := []string{
		"testrepo",
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
			depType: 1,
			path:    n,
		}

		deps = append(deps, d)
	}

	end := len(deps) - 2 // get the unicode ver
	deps[end].version = "v0.3.3"

	for _, n := range nodejsPaths {
		d := Dependency{
			depType: 2,
			path:    n,
		}
		deps = append(deps, d)
	}

	return deps
}

func TestGetDeps(t *testing.T) {
	want := BuildWant()

	got, gotBitmask, _ := GetDeps("test/testRepo/")

	if gotBitmask != 3 {
		t.Errorf("GotBitmask() != 3")
	}

	// iterate thru and compare
	if len(want) != len(got) {
		t.Errorf("GetDeps() = %d; want %d", len(got), len(want))
	}

	for i, pkg := range want {
		if pkg.path != got[i].path {
			t.Errorf("GetDeps() got %s; want %s", got[i].path, pkg.path)
		}

		if pkg.version != "" && pkg.version != got[i].version {
			t.Errorf("GetDeps() got %s %s; want %s %s", got[i].path, got[i].version, pkg.path, pkg.version)
		}
	}
}
