package main

import (
	"flag"
	"fmt"

	"github.com/mcoops/deplist"
)

func main() {
	deptypePtr := flag.Int("deptype", deplist.LangGolang, "golang, nodejs, python etc")

	flag.Parse()
	path := flag.Args()[0]
	deptype := deplist.Bitmask(*deptypePtr)

	deps, _, err := deplist.GetDeps(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, dep := range deps {
		if (dep.DepType & deptype) == deptype {
			fmt.Printf("%s %s\n", dep.Path, dep.Version)
		}
	}
}
