/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package rule

import (
	"fluent/analyzer/pool"
	function2 "fluent/analyzer/rule/function"
	"fluent/filecode"
)

func AnalyzeFileCode(code filecode.FileCode) *pool.ErrorPool {
	globalErrors := pool.NewErrorPool()

	// Iterate over all the functions
	for _, function := range code.Functions {
		// Skip functions that are not in the current file
		if function.Path != code.Path {
			continue
		}

		// Analyze the function
		errors := function2.AnalyzeFunction(function, &code)
		globalErrors.Extend(errors.Errors)
	}

	return globalErrors
}
