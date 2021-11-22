package main

import (
	"flag"
	"fmt"

	"github.com/RedHatProductSecurity/deplist"
	purl "github.com/mcoops/packageurl-go"
	log "github.com/sirupsen/logrus"
)

func main() {
	deptypePtr := flag.Int("deptype", -1, "golang, nodejs, python etc")
	debugPtr := flag.Bool("debug", false, "debug logging (default false)")

	flag.Parse()

	if *debugPtr == true {
		log.SetLevel(log.DebugLevel)
	}

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
			version := dep.Version

			inst, _ := purl.FromString(fmt.Sprintf("pkg:%s/%s@%s", deplist.GetLanguageStr(dep.DepType), dep.Path, version))
			fmt.Println(inst)
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
