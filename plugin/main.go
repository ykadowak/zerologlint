// This file can build as a plugin for golangci-lint by below command.
//
//	go build -buildmode=plugin -o path_to_plugin_dir github.com/ykadowak/zerologlint/plugin/zerologlint
//
// See: https://golangci-lint.run/contributing/new-linters/#how-to-add-a-private-linter-to-golangci-lint
package main

import (
	"github.com/ykadowak/zerologlint"
	"golang.org/x/tools/go/analysis"
)

// AnalyzerPlugin provides analyzers as a plugin.
// It follows golangci-lint style plugin.
var AnalyzerPlugin analyzerPlugin

type analyzerPlugin struct{}

func (analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		zerologlint.Analyzer,
	}
}
