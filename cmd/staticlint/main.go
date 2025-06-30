// Package main implements a multichecker for static analysis.
//
// Запуск: go run ./cmd/staticlint [пакеты]
//
// Анализаторы:
//   - Стандартные анализаторы из golang.org/x/tools/go/analysis/passes
//   - Все SA-анализаторы из staticcheck.io
//   - Один анализатор другого класса staticcheck (например, ST1000)
//   - go-critic (различные code-style и error-prone проверки)
//   - nilnil (проверка двойных nil в return)
//   - noosexit (запрет os.Exit в main.main)
package main

import (
	"flag"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"

	gc_analyzer "github.com/go-critic/go-critic/checkers/analyzer"
	bodyclose_analyzer "github.com/timakin/bodyclose/passes/bodyclose"

	"os"

	"honnef.co/go/tools/staticcheck"
)

func main() {
	var analyzers []*analysis.Analyzer

	// Стандартные анализаторы
	analyzers = append(analyzers,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	)

	// SA анализаторы staticcheck
	for _, v := range staticcheck.Analyzers {
		if len(v.Analyzer.Name) >= 2 && v.Analyzer.Name[:2] == "SA" {
			analyzers = append(analyzers, v.Analyzer)
		}
	}
	// Один анализатор другого класса staticcheck (например, ST1000)
	for _, v := range staticcheck.Analyzers {
		if v.Analyzer.Name == "ST1000" {
			analyzers = append(analyzers, v.Analyzer)
			break
		}
	}

	// Публичные анализаторы
	analyzers = append(analyzers, gc_analyzer.Analyzer)        // go-critic
	analyzers = append(analyzers, bodyclose_analyzer.Analyzer) // bodyclose

	// Ваш собственный анализатор
	analyzers = append(analyzers, Analyzer)

	multichecker.Main(analyzers...)
}

func init() {
	var skipTest bool
	flag.BoolVar(&skipTest, "skip-test", false, "skip test files")
	flag.Parse()
	if skipTest {
		// This is a hack to exclude test files.
		// A better way would be to properly configure the packages to be checked.
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}
}
