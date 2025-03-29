/*
   The Fluent Programming Language
   -----------------------------------------------------
   This code is released under the GNU GPL v3 license.
   For more information, please visit:
   https://www.gnu.org/licenses/gpl-3.0.html
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package rule

import (
	"fluent/analyzer/pool"
	function2 "fluent/analyzer/rule/function"
	"fluent/analyzer/rule/module"
	"fluent/analyzer/stack"
	"fluent/filecode"
)

// AnalyzeFileCode analyzes the given FileCode and returns two ErrorPools:
// one for errors and one for warnings. It iterates over all functions and
// modules in the FileCode, analyzing each and collecting errors and warnings.
//
// Parameters:
//   - code: The FileCode to be analyzed.
//
// Returns:
//   - *pool.ErrorPool: A pool of errors found during the analysis.
//   - *pool.ErrorPool: A pool of warnings found during the analysis.
func AnalyzeFileCode(code *filecode.FileCode) (*pool.ErrorPool, *pool.ErrorPool) {
	globalErrors := pool.NewErrorPool()
	globalWarnings := pool.NewErrorPool()

	// Iterate over all the functions
	for _, function := range code.Functions {
		// Skip functions that are not in the current file
		if function.Path != code.Path {
			continue
		}

		// Analyze the function
		errors, warnings, _ := function2.AnalyzeFunction(*function, code, "", &function.Templates, stack.ScopedStack{
			Scopes: make(map[int]stack.Stack),
		}, false)
		globalErrors.Extend(errors.Errors)
		globalWarnings.Extend(warnings.Errors)
	}

	// Iterate over all the modules
	for _, mod := range code.Modules {
		// Skip functions that are not in the current file
		if mod.Path != code.Path {
			continue
		}

		// Analyze the function
		errors, warnings := module.AnalyzeModule(*mod, code)
		globalErrors.Extend(errors.Errors)
		globalWarnings.Extend(warnings.Errors)
	}

	return globalErrors, globalWarnings
}
