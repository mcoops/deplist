package deplist

import (
	"testing"
)

func BuildWant() []Dependency {
	var deps []Dependency

	golangPaths := []string{
		"errors",
		"fmt",
		"github.com/mcoops/deplist",
		"github.com/openshift/api/config/v1",
		"golang.org/x/text/unicode",
		"internal/abi",
		"internal/bytealg",
		"internal/cpu",
		"internal/fmtsort",
		"internal/goexperiment",
		"internal/itoa",
		"internal/oserror",
		"internal/poll",
		"internal/race",
		"internal/reflectlite",
		"internal/syscall/execenv",
		"internal/syscall/unix",
		"internal/testlog",
		"internal/unsafeheader",
		"io",
		"io/fs",
		"math",
		"math/bits",
		"os",
		"path",
		"reflect",
		"runtime",
		"runtime/internal/atomic",
		"runtime/internal/math",
		"runtime/internal/sys",
		"sort",
		"strconv",
		"sync",
		"sync/atomic",
		"syscall",
		"time",
		"unicode",
		"unicode/utf8",
		"unsafe",
	}

	npmSet1 := []string{
		"loose-envify",
		"iconv-lite",
		"d3-brush",
		"d3-zoom",
		"rw",
		"d3-ease",
		"object-assign",
		"commander",
		"d3-dsv",
		"d3-scale",
		"is-plain-object",
		"d3-quadtree",
		"tiny-warning",
		"d3-hierarchy",
		"d3-scale-chromatic",
		"d3-axis",
		"d3-color",
		"prismjs",
		"iconv-lite",
		"angular",
		"d3-delaunay",
		"rxjs",
		"d3-path",
		"d3-array",
		"js-tokens",
		"d3-contour",
		"safer-buffer",
		"react-is",
		"d3-dispatch",
		"d3-force",
		"prop-types",
		"tiny-emitter",
		"d3-polygon",
		"d3-chord",
		"d3-fetch",
		"tslib",
		"good-listener",
		"d3",
		"delegate",
		"d3-drag",
		"delaunator",
		"d3-timer",
		"d3-geo",
		"slate",
		"select",
		"esrever",
		"d3-transition",
		"clipboard",
		"d3-format",
		"d3-random",
		"d3-shape",
		"d3-time",
		"immer",
		"@types/esrever",
		"d3-time-format",
		"d3-selection",
		"react",
		"tether",
		"d3-interpolate",
	}

	rubySet := []string{
		"fluent-plugin-kafka",
		"fluent-plugin-rewrite-tag-filter",
		"faraday",
		"concurrent-ruby",
		"elasticsearch",
		"sigdump",
		"syslog_protocol",
		"uuidtools",
		"aws-partitions",
		"http-cookie",
		"ltsv",
		"quantile",
		"connection_pool",
		"tzinfo-data",
		"unf",
		"aws-sdk-core",
		"fluent-plugin-cloudwatch-logs",
		"fluent-plugin-kubernetes_metadata_filter",
		"http-form_data",
		"jmespath",
		"kubeclient",
		"msgpack",
		"rest-client",
		"dig_rb",
		"unf_ext",
		"to_regexp",
		"cool.io",
		"ethon",
		"ffi",
		"fluent-plugin-remote-syslog",
		"lru_redux",
		"prometheus-client",
		"typhoeus",
		"aws-sdk-cloudwatchlogs",
		"fluent-plugin-prometheus",
		"mime-types",
		"mime-types-data",
		"public_suffix",
		"domain_name",
		"aws-sigv4",
		"elasticsearch-api",
		"excon",
		"fluentd",
		"http_parser.rb",
		"netrc",
		"recursive-open-struct",
		"aws-eventstream",
		"systemd-journal",
		"net-http-persistent",
		"elasticsearch-transport",
		"ffi-compiler",
		"fluent-plugin-record-modifier",
		"http-accept",
		"http-parser",
		"rake",
		"digest-crc",
		"fluent-plugin-splunk-hec",
		"fluent-plugin-systemd",
		"multi_json",
		"multipart-post",
		"ruby-kafka",
		"strptime",
		"fluent-plugin-concat",
		"serverengine",
		"fluent-plugin-multi-format-parser",
		"tzinfo",
		"fluent-mixin-config-placeholders",
		"jsonpath",
		"fluent-config-regexp-type",
		"fluent-plugin-elasticsearch",
		"http",
		"yajl-ruby",
		"addressable",
	}

	pythonSet := []string{
		"cotyledon",
		"Flask",
		"kuryr-lib",
		"docutils",
		"python-dateutil",
		"unittest2",
		"cryptography",
	}

	for _, n := range golangPaths {
		d := Dependency{
			DepType: 1,
			Path:    n,
		}

		deps = append(deps, d)
	}

	deps[4].Version = "v0.3.3" // test golang.org/x/text/unicode version

	for _, n := range npmSet1 {
		d := Dependency{
			DepType: 3,
			Path:    n,
		}
		deps = append(deps, d)
	}
	deps = append(deps, Dependency{DepType: 2, Path: "com.amazonaws:aws-lambda-java-core:jar", Version: "1.0.0"}) // java

	for _, n := range rubySet {
		d := Dependency{
			DepType: LangRuby,
			Path:    n,
		}
		deps = append(deps, d)
	}

	for _, n := range pythonSet {
		d := Dependency{
			DepType: LangPython,
			Path:    n,
		}
		deps = append(deps, d)
	}

	end := len(deps) - 2 // get the cryptography ver
	deps[end].Version = "0.5.1"

	end = len(deps) - 1 // get the cryptography ver
	deps[end].Version = "2.3.0"

	return deps
}

func TestGetDeps(t *testing.T) {
	want := BuildWant()

	got, gotBitmask, _ := GetDeps("test/testRepo")

	if gotBitmask != 31 {
		t.Errorf("GotBitmask() != 31; got: %d", gotBitmask)
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
