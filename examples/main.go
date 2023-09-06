package main

import (
	"fmt"
	"log"

	"jdel.org/go-syno"
)

func main() {
	o := syno.GetOptions()
	o.CacheDir = "./cache"
	o.PackagesDir = "../tests/packages"

	fmt.Println("Options:")
	fmt.Println(o)
	fmt.Println("--------------------")

	m, err := syno.GetModels(false)
	if err != nil {
		log.Fatalln("Exiting due to fatal error", err)
	}
	fmt.Println("Models:")
	fmt.Println(m)
	fmt.Println("--------------------")

	ps := syno.Packages{}

	// Load a package from disk
	p, err := syno.NewPackage("real-package.spk")
	if err != nil {
		fmt.Println("Skipping real-package.spk", err)
	} else {
		ps = append(ps, p)
	}

	// Attempt to load a bad package from disk
	p, err = syno.NewPackage("bad-package.spk")
	if err != nil {
		fmt.Println("Skipping bad-package.spk", err)
	} else {
		ps = append(ps, p)
	}

	// Create a debug package
	ps = append(ps, syno.NewDebugPackage(o.String()))

	fmt.Println("Packages:")
	fmt.Println(ps)
	fmt.Println("--------------------")
}
