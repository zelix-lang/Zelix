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

package function

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/format"
	"fluent/analyzer/rule/value"
	"fluent/filecode"
	"fluent/filecode/function"
)

// AnalyzeParameter analyzes a function parameter for various conditions.
// It checks if the parameter name is in snake_case, if the parameter type is "nothing",
// and if there are any undefined references.
//
// Parameters:
// - name: A pointer to the parameter name string.
// - param: A pointer to the function parameter structure.
// - trace: A pointer to the file code trace structure.
// - generics: A pointer to a map of generics.
//
// Returns:
// - An error if the parameter type is "nothing".
// - A warning if the parameter name is not in snake_case.
func AnalyzeParameter(
	name *string,
	param *function.Param,
	trace *filecode.FileCode,
	generics *map[string]bool,
) (*error3.Error, *error3.Error) {
	var warning *error3.Error

	// Check that the case matches snake_case
	if !format.CheckCase(name, format.SnakeCase) {
		warning = &error3.Error{
			Code:       error3.NameShouldBeSnakeCase,
			Line:       param.Trace.Line,
			Column:     param.Trace.Column,
			Additional: []string{*name},
		}
	}

	// Check that the type is not "nothing"
	if param.Type.BaseType == "nothing" {
		return &error3.Error{
			Code:   error3.ParamTypeNothing,
			Line:   param.Trace.Line,
			Column: param.Trace.Column,
		}, warning
	}

	// Check for undefined references
	return value.AnalyzeUndefinedReference(trace, param.Type, generics), warning
}
