package zerologlint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "zerologlint",
	Doc:  "check that zerolog log methods have a final Msg call",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		has, name := checkZerologImport(f)
		if !has {
			continue
		}

		inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		exprFilter := []ast.Node{
			(*ast.ExprStmt)(nil),
		}
		inspector.Preorder(exprFilter, func(n ast.Node) {
			es, _ := n.(*ast.ExprStmt)
			if has := hasLogIdent(es, name); has {
				if has := hasMsg(es); !has {
					pass.Reportf(n.Pos(), "missing Msg or Send call for zerolog log method")
				}
			}
		})
	}
	return nil, nil
}

func checkZerologImport(f *ast.File) (bool, string) {
	for _, imp := range f.Imports {
		if imp.Path.Value == `"github.com/rs/zerolog/log"` {
			if imp.Name == nil {
				// no alias
				return true, "log"
			}
			return true, imp.Name.Name
		}
	}
	return false, ""
}

func hasLogIdent(es *ast.ExprStmt, name string) bool {
	has := false
	ast.Inspect(es, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			if ident.Name == name {
				has = true
				return false
			}
		}
		return true
	})
	return has
}

func hasMsg(es *ast.ExprStmt) bool {
	if ce, ok := es.X.(*ast.CallExpr); ok {
		if se, ok := ce.Fun.(*ast.SelectorExpr); ok {
			switch se.Sel.Name {
			case "Msg", "Msgf", "Send":
				return true
			}
		}
	}
	return false
}
