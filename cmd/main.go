package main

import (
	"flag"
	"os"
)

var options struct {
	version   bool
}

func main() {
	flag.BoolVar(&options.version,"version", false, "View the version of maxo lang")
	flag.Parse()
	if options.version {
		print("MAXOv0.0.1")
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		print("type -h or --help for help on usage")
		os.Exit(1)
	}
}