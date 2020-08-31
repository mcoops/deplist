package deplist

// Bitmask type allows easy tagging of what langs there are
type Bitmask uint32

// Dependency per dependency info
type Dependency struct {
	depType int      // golang, nodejs, python etc
	path    string   // the module path, github.com/teris-io/shortid
	version string   // v0.0.0-20171029131806-771a37caa5cf
	files   []string // if available, list of all files for a package
	// /usr/lib/go-1.13/src/regexp/syntax/compile.go
	// /usr/lib/go-1.13/src/regexp/syntax/doc.go
}

// DepFoundAddFlag add a lang type to the bitmask
func (f *Bitmask) DepFoundAddFlag(flag Bitmask) { *f |= flag }

// DepFoundHasFlag deteremine if bitmask has a lang type
func (f Bitmask) DepFoundHasFlag(flag Bitmask) bool { return f&flag != 0 }
