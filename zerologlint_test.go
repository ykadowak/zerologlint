package zerologlint_test

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	"github.com/ykadowak/zerologlint"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, zerologlint.Analyzer, "a")
}

// TestAnalyzerWithAdditionalPrefix tests that the -prefix flag correctly extends
// the analyzer to handle a zerolog fork under a custom module path.
func TestAnalyzerWithAdditionalPrefix(t *testing.T) {
	a := *zerologlint.Analyzer
	a.Flags.Init("zerologlint", 0)
	a.Flags.Set("prefix", "myfork/myzerolog")
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, &a, "b")
}
