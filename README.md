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
	langGolang = 1 << iota
	langNodeJS
	langPython
	langRuby
)
```

* **error:**

  Standard Go error handling

## Command Line

```bash
$ go build cmd/deplist.go
$ ./deplist path/to/repo
$ ./deplist -deptype 2 path/to/repo # NodeJs deps
```
