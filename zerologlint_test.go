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

// TestAnalyzerWithCheckMsgf tests the analyzer with the checkmsgf flag enabled.
func TestAnalyzerWithCheckMsgf(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)

	// Enable the checkmsgf flag
	analyzer := zerologlint.Analyzer
	if err := analyzer.Flags.Set("checkmsgf", "true"); err != nil {
		t.Fatalf("failed to set checkmsgf flag: %v", err)
	}

	analysistest.Run(t, testdata, analyzer, "msgf")
}
