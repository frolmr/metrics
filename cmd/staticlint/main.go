// Package staticlint implements a custom static analysis multi-checker that combines:
// - Standard Go analyzers from golang.org/x/tools/go/analysis/passes
// - All SA-class analyzers from staticcheck.io
// - Selected analyzers from other classes of staticcheck.io
// - Additional public analyzers
// - Custom analyzers (e.g., forbidding os.Exit in main function)
//
// The included analyzers are:
//
// Standard Go analyzers:
//   - asmdecl: report mismatches between assembly files and Go declarations
//   - assign: check for useless assignments
//   - atomic: check for common mistakes using the sync/atomic package
//   - bools: check for common mistakes involving boolean operators
//   - buildtag: check that +build tags are well-formed and correctly located
//   - cgocall: detect some violations of the cgo pointer passing rules
//   - composite: check for unkeyed composite literals
//   - copylock: check for locks erroneously passed by value
//   - errorsas: report passing non-pointer or non-error values to errors.As
//   - httpresponse: check for mistakes using HTTP responses
//   - ifaceassert: detect impossible interface-to-interface type assertions
//   - loopclosure: check references to loop variables from within nested functions
//   - lostcancel: check for failure to call a context cancellation function
//   - nilfunc: check for useless comparisons between functions and nil
//   - nilness: check for redundant or impossible nil comparisons
//   - printf: check consistency of Printf format strings and arguments
//   - shadow: check for possible unintended shadowing of variables
//   - shift: check for shifts that equal or exceed the width of the integer
//   - sigchanyzer: check for unbuffered channel of os.Signal
//   - sortslice: check the argument type of sort.Slice
//   - stdmethods: check signature of methods of well-known interfaces
//   - stringintconv: check for string(int) conversions
//   - structtag: check that struct field tags conform to reflect.StructTag's rules
//   - testinggoroutine: report calls to (*testing.T).Fatal from goroutines started by a test
//   - tests: check for common mistaken usages of tests and examples
//   - unmarshal: report passing non-pointer or non-interface values to unmarshal
//   - unreachable: check for unreachable code
//   - unsafeptr: check for invalid conversions of uintptr to unsafe.Pointer
//   - unusedresult: check for unused results of calls to certain pure functions
//   - fieldalignment:disabled: find structs that would use less memory if their fields were sorted
//
// Staticcheck analyzers:
//   - All SA* analyzers: catch a wide range of bugs and suspicious constructs
//   - ST1001: Dot imports are discouraged
//   - QF1001: Apply De Morgan's law
//   - S1001: Replace for loop with call to copy
//
// Additional public analyzers:
//   - errchkjson: Checks for missing error checking around json marshaling/unmarshaling
//   - bodyclose: Checks for unclosed HTTP response bodies
//
// Custom analyzers:
//   - noosexit: Reports usage of os.Exit in main function of main package
//
// # Usage
//
// To run the staticlint tool:
//
//	go install yourmodule/cmd/staticlint
//	staticlint ./...
//
// Alternatively, run directly without installation:
//
//	go run yourmodule/cmd/staticlint ./...
//
// You can also run specific analyzers by name:
//
//	staticlint -analyzer=analyzer_name ./...
package main

import (
	"go/ast"
	"strings"

	"github.com/breml/errchkjson"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
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
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {
	// Standard Go analyzers from golang.org/x/tools/go/analysis/passes
	stdAnalyzers := []*analysis.Analyzer{
		asmdecl.Analyzer,          // Report mismatches between assembly files and Go declarations
		assign.Analyzer,           // Check for useless assignments
		atomic.Analyzer,           // Check for common mistakes using the sync/atomic package
		bools.Analyzer,            // Check for common mistakes involving boolean operators
		buildtag.Analyzer,         // Check that +build tags are well-formed and correctly located
		cgocall.Analyzer,          // Detect some violations of the cgo pointer passing rules
		composite.Analyzer,        // Check for unkeyed composite literals
		copylock.Analyzer,         // Check for locks erroneously passed by value
		errorsas.Analyzer,         // Report passing non-pointer or non-error values to errors.As
		httpresponse.Analyzer,     // Check for mistakes using HTTP responses
		ifaceassert.Analyzer,      // Detect impossible interface-to-interface type assertions
		loopclosure.Analyzer,      // Check references to loop variables from within nested functions
		lostcancel.Analyzer,       // Check for failure to call a context cancellation function
		nilfunc.Analyzer,          // Check for useless comparisons between functions and nil
		nilness.Analyzer,          // Check for redundant or impossible nil comparisons
		printf.Analyzer,           // Check consistency of Printf format strings and arguments
		shadow.Analyzer,           // Check for possible unintended shadowing of variables
		shift.Analyzer,            // Check for shifts that equal or exceed the width of the integer
		sigchanyzer.Analyzer,      // Check for unbuffered channel of os.Signal
		sortslice.Analyzer,        // Check the argument type of sort.Slice
		stdmethods.Analyzer,       // Check signature of methods of well-known interfaces
		stringintconv.Analyzer,    // Check for string(int) conversions
		structtag.Analyzer,        // Check that struct field tags conform to reflect.StructTag's rules
		testinggoroutine.Analyzer, // Report calls to (*testing.T).Fatal from goroutines started by a test
		tests.Analyzer,            // Check for common mistaken usages of tests and examples
		unmarshal.Analyzer,        // Report passing non-pointer or non-interface values to unmarshal
		unreachable.Analyzer,      // Check for unreachable code
		unsafeptr.Analyzer,        // Check for invalid conversions of uintptr to unsafe.Pointer
		unusedresult.Analyzer,     // Check for unused results of calls to certain pure functions
		//fieldalignment.Analyzer,   // Find structs that would use less memory if their fields were sorted
	}

	selectedChecks := map[string]bool{
		"ST1001": true,
		"QF1001": true,
		"S1001":  true,
	}

	// Collect all analyzers from staticcheck.io packages
	var staticCheckAnalyzers []*analysis.Analyzer

	for _, analyzer := range staticcheck.Analyzers {
		if strings.HasPrefix(analyzer.Analyzer.Name, "SA") {
			staticCheckAnalyzers = append(staticCheckAnalyzers, analyzer.Analyzer)
		}
	}

	for _, analyzer := range stylecheck.Analyzers {
		if selectedChecks[analyzer.Analyzer.Name] {
			staticCheckAnalyzers = append(staticCheckAnalyzers, analyzer.Analyzer)
		}
	}

	for _, analyzer := range quickfix.Analyzers {
		if selectedChecks[analyzer.Analyzer.Name] {
			staticCheckAnalyzers = append(staticCheckAnalyzers, analyzer.Analyzer)
		}
	}

	for _, analyzer := range simple.Analyzers {
		if selectedChecks[analyzer.Analyzer.Name] {
			staticCheckAnalyzers = append(staticCheckAnalyzers, analyzer.Analyzer)
		}
	}

	// Additional public analyzers
	otherPublicAnalyzers := []*analysis.Analyzer{
		errchkjson.NewAnalyzer(), // Checks for missing error checking around json marshaling/unmarshaling
		bodyclose.Analyzer,       // Checks whether res.Body is correctly closed
	}

	var analyzers []*analysis.Analyzer
	analyzers = append(analyzers, stdAnalyzers...)
	analyzers = append(analyzers, staticCheckAnalyzers...)
	analyzers = append(analyzers, otherPublicAnalyzers...)
	analyzers = append(analyzers, noOsExitAnalyzer)

	multichecker.Main(analyzers...)
}

// noOsExitAnalyzer is a custom analyzer that reports usage of os.Exit in main function of main package.
// It helps ensure proper program shutdown by preventing direct os.Exit calls in main functions,
// which can skip deferred cleanup operations and leave resources in an undefined state.
var noOsExitAnalyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "reports usage of os.Exit in main function of main package",
	Run:  runNoOsExit,
}

// runNoOsExit implements the analysis logic for the noosexit analyzer.
// It scans for os.Exit calls in the main function of the main package.
func runNoOsExit(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if isInCacheDir(pass.Fset.Position(file.Pos()).Filename) {
			continue
		}
		if pass.Pkg.Name() != "main" {
			continue
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "main" {
				continue
			}

			ast.Inspect(fn, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				ident, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}

				if ident.Name == "os" && sel.Sel.Name == "Exit" {
					pass.Reportf(call.Pos(), "os.Exit call forbidden in main function of main package")
				}
				return true
			})
		}
	}
	return nil, nil
}

func isInCacheDir(filename string) bool {
	cachePatterns := []string{
		"/tmp/",
		"/var/folders/",
		"/.cache/",
		"/go-build",
	}

	for _, pattern := range cachePatterns {
		if strings.Contains(filename, pattern) {
			return true
		}
	}
	return false
}
