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

// Config is a set of configuration values which configure the linter behavior.
type Config struct {
	// AllowPackages defines a list of packages that linter should check.
	//
	// Generally, package main will be only allowed package, but for testing cases
	// you may include other packages.
	AllowPackages []string `mapstructure:"allowPackages" yaml:"allowPackages"`
}

// DefaultConfig is a default linter config.
var DefaultConfig = Config{
	AllowPackages: DefaultAllowPackages,
}

// NewAnalyzer creates new nomainreturn analyzer.
func NewAnalyzer(cfg Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "nomainreturn",
		Doc:      "reports use of return keyword in the main",
		Run:      run(cfg),
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

type analysisRunner func(*analysis.Pass) (interface{}, error)

func run(cfg Config) analysisRunner {
	return func(pass *analysis.Pass) (interface{}, error) {
		if pass.Pkg == nil {
			return nil, nil
		}

		if pkgAllowed := isPkgAllowed(pass.Pkg.Name(), cfg.AllowPackages); !pkgAllowed {
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
}

func isPkgAllowed(pkg string, allowedPkgs []string) bool {
	var allowed bool
	for _, p := range allowedPkgs {
		if p == pkg {
			allowed = true
			break
		}
	}
	return allowed
}
