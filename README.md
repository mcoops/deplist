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
