# deplist

Scans a given repo for Golang, NodeJS (more comming) for dependencies.

The api functions as follows:

```
func GetDeps(fullPath string) ([]Dependency, Bitmask, error) {
```

### Parameters

* **fullPath:**

  To the repository to scan. Expects it to be present locally first.

### Returns

* **Depenency:**
  
  Array of Dependency structs from [dependencies.go](dependencies.go)


* **Bitmask:**

  A bitmask of found languages:

```
const (
	LangGolang = 1 << iota
	LangNodeJS
	LangPython
	LangRuby
)
```

* **error:**

  Standard Go error handling

## Command Line

```bash
$ make 
$ ./deplist path/to/repo # Go deps
golang.org/x/tools/go/gcexportdata v0.0.0-20201223010750-3fa0e8f87c1a
golang.org/x/tools/internal/gocommand v0.0.0-20201223010750-3fa0e8f87c1a
flag
...
$ ./deplist -deptype 4 path/to/repo # NodeJs deps
whatwg-fetch 3.1.0
less-loader 5.0.0
pseudomap 1.0.2
...
```
