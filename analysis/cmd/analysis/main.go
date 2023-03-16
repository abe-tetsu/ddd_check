package main

import (
	"analysis"

	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(analysis.AnalyzerID) }
