package main

import (
	"flag"
	"fmt"
	"github.com/chain-exporter/exporter"
	"github.com/chain-exporter/version"
)

var v bool

func main() {
	flag.BoolVar(&v, "version", false, "show version and exit")
	flag.Parse()

	if v {
		//fmt.Printf("version:\t%s \nbuild time:\t%s\ngit branch:\t%s\ngit commit:\t%s\ngo version:\t%s\n", version.VERSION, version.BUILD_TIME, version.GIT_BRANCH, version.COMMIT_SHA1, version.GO_VERSION)
		fmt.Printf("version:\t%s\n", version.Version)
		return
	}

	// Start exporting chain data
	exporter := exporter.NewExporter()
	exporter.Start()
}
