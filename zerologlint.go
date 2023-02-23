package zerologlint

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "logmsg",
	Doc:  "check that zerolog log methods have a final Msg call",
	Run:  run,
}


func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		if !hasZerologImport(f) {
			continue
		}
		ast.Inspect(f, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.CallExpr:
				if isLogMethod(n.Fun, pass.TypesInfo) && !hasMsgCall(n.Args, pass.TypesInfo) {
					pass.Reportf(n.Pos(), "missing Msg call for zerolog log method")
				}
			}
			return true
		})
	}
	return nil, nil
}

func hasZerologImport(f *ast.File) bool {
	for _, imp := range f.Imports {
		if imp.Path.Value == `"github.com/rs/zerolog/log"` {
			return true
		}
	}
	return false
}

func isLogMethod(expr ast.Expr, info *types.Info) bool {
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		if id, ok := sel.X.(*ast.Ident); ok {
			if obj := info.ObjectOf(id); obj != nil {
				switch sel.Sel.Name {
				case "Debug", "Info", "Warn", "Error", "Fatal", "Panic":
					return true
				}
			}
		}
	}
	return false
}

func hasMsgCall(args []ast.Expr, info *types.Info) bool {
	for _, arg := range args {
		if call, ok := arg.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Msg" && len(call.Args) == 1 {
				return true
			}
		}
	}
	return false
}
