package analysis

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const docStruct = "analysis is ..."

// AnalyzerIDConstructor is ...
var AnalyzerStruct = &analysis.Analyzer{
	Name: "analysis",
	Doc:  docStruct,
	Run:  runStruct,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runStruct(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.Ident)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.Ident:
			if n.Name == "gopher" {
				pass.Reportf(n.Pos(), "identifier is aaa")
			}
		}
	})

	return nil, nil
}
