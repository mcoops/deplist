package main

import (
	"flag"
	"fmt"

	"github.com/mcoops/deplist"
)

func main() {
	deptypePtr := flag.Int("deptype", -1, "golang, nodejs, python etc")

	flag.Parse()

	if flag.Args() == nil || len(flag.Args()) == 0 {
		fmt.Println("Not path to scan was specified, i.e. deplist /tmp/files/")
		return
	}

	path := flag.Args()[0]

	deps, _, err := deplist.GetDeps(path)
	if err != nil {
		fmt.Println(err.Error())
	}

	if *deptypePtr == -1 {
		for _, dep := range deps {
			fmt.Printf("pkg:%s/%s@%s\n", deplist.GetLanguageStr(dep.DepType), dep.Path, dep.Version)
		}
	} else {
		deptype := deplist.Bitmask(*deptypePtr)
		for _, dep := range deps {
			if (dep.DepType & deptype) == deptype {
				fmt.Printf("%s@%s\n", dep.Path, dep.Version)
			}
		}
	}
}
