package ddd_check_test

import (
	"testing"

	"github.com/abe-tetsu/ddd_check"

	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for AnalyzerIDConstructor.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, ddd_check.AnalyzerID, "a/...")
}
