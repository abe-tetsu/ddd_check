package main

import (
	"github.com/abe-tetsu/ddd_check"

	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(ddd_check.AnalyzerID) }
