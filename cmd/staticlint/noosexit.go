// Package main запрещает прямой вызов os.Exit в функции main пакета main.
package main

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "Запрещает прямой вызов os.Exit в функции main пакета main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			funcDecl, ok := node.(*ast.FuncDecl)
			if !ok || funcDecl.Name.Name != "main" {
				return true
			}

			ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
				callExpr, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}

				selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				if obj, ok := pass.TypesInfo.Uses[selectorExpr.Sel]; ok {
					if fn, ok := obj.(*types.Func); ok {
						if fn.FullName() == "os.Exit" {
							pass.Reportf(selectorExpr.Pos(), "запрещён прямой вызов os.Exit в main.main")
						}
					}
				}
				return true
			})

			return false
		})
	}

	return nil, nil
}
