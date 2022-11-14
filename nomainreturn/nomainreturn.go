// Package nomainreturn defines an Analyzer that reports use of return keyword in the main.
package nomainreturn

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer is a nomainreturn analyzer.
var Analyzer = &analysis.Analyzer{
	Name:     "nomainreturn",
	Doc:      "reports use of return keyword in the main",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}
	inspector.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return
		}

		funcName := funcDecl.Name.Name
		if funcName != "main" {
			return
		}

		for _, stmt := range funcDecl.Body.List {
			switch stmt.(type) {
			case *ast.AssignStmt, *ast.DeclStmt:
				continue
			}

			ast.Inspect(stmt, func(n ast.Node) bool {
				switch n.(type) {
				case *ast.ReturnStmt:
					pass.Reportf(n.Pos(), "return found in main")
				}
				return true
			})
		}
	})

	return nil, nil
}
