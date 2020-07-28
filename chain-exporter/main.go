package main

import (
	"github.com/chain-exporter/exporter"
)

func main() {
	// Start exporting chain data
	exporter := exporter.NewExporter()
	exporter.Start()
}
