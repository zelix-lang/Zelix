/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package function

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/format"
	"fluent/analyzer/rule/value"
	"fluent/filecode"
	"fluent/filecode/function"
	"fluent/message/warn"
	"fluent/state"
	"fluent/util"
)

func AnalyzeParameter(
	name *string,
	param *function.Param,
	trace *filecode.FileCode,
	generics *map[string]bool,
) error3.Error {
	// Check that the case matches snake_case
	if !format.CheckCase(name, format.SnakeCase) {
		state.WarnAllSpinners()
		warn.SnakeCase(*name)
		util.BuildAndPrintDetails(
			&trace.Contents,
			&trace.Path,
			param.Trace.Line,
			param.Trace.Column,
			false,
		)
	}

	// Check that the type is not "nothing"
	if param.Type.BaseType == "nothing" {
		return error3.Error{
			Code:   error3.ParamTypeNothing,
			Line:   param.Trace.Line,
			Column: param.Trace.Column,
		}
	}

	// Check for undefined references
	return value.AnalyzeUndefinedReference(trace, param.Type, generics)
}
