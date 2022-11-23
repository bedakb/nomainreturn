// Package nomainreturn defines an Analyzer that reports use of return keyword in the main.
package nomainreturn

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// DefaultAllowPackages is an allow-list of packages to include when running the linter.
var DefaultAllowPackages = []string{"main"}

// NoMainReturnConfig is a set of configuration values which configure the linter behavior.
type NoMainReturnConfig struct {
	// AllowPackages defines a list of packages that linter should check.
	//
	// Generally, package main will be only allowed package, but for testing cases
	// you may include other packages.
	AllowPackages []string `mapstructure:"allowPackages" yaml:"allowPackages"`
}

// NewDefaultConfig returns default linter config.
func NewDefaultConfig() NoMainReturnConfig {
	return NoMainReturnConfig{AllowPackages: DefaultAllowPackages}
}

// NewAnalyzer creates new nomainreturn analyzer.
func NewAnalyzer(cfg NoMainReturnConfig) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "nomainreturn",
		Doc:      "reports use of return keyword in the main",
		Run:      run(cfg),
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

type analysisRunner func(*analysis.Pass) (interface{}, error)

func run(cfg NoMainReturnConfig) analysisRunner {
	return func(pass *analysis.Pass) (interface{}, error) {
		for _, pkg := range cfg.AllowPackages {
			if pass.Pkg != nil && pass.Pkg.Name() != pkg {
				return nil, nil
			}
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
}
