package checker

import (
	"strconv"
	"surf/ast"
	"surf/code"
	"surf/core/stack"
	"surf/logger"
	"surf/object"
	"surf/tokenUtil"
	"time"
)

// AnalyzeFun analyzes the given function
func AnalyzeFun(
	function *ast.Function,
	functions *map[string]map[string]*ast.Function,
	trace code.Token,
	args ...object.SurfObjectType,
) object.SurfObjectType {
	// Standard functions are not evaluated
	if function.IsStd() {
		return object.NothingType
	}

	function.SetTimesCalled(function.GetTimesCalled() + 1)
	function.SetLastCalled(time.Now())

	// Check for stack overflows
	if function.GetTimesCalled() > 1000 && time.Since(function.GetLastCalled()).Seconds() < 1 {
		logger.TokenError(
			trace,
			"Stack overflow",
			"This function has overflown its stack",
			"Check for infinite loops",
		)
	}

	if len(args) != len(function.GetParameters()) {
		logger.TokenError(
			trace,
			"Invalid number of arguments",
			"This function expects "+strconv.Itoa(len(function.GetParameters()))+" arguments",
			"Add the missing arguments",
		)
	}

	// Create a new StaticStack
	variables := stack.NewStaticStack()

	// Used to skip tokens
	skipToIndex := 0

	for i, token := range function.GetBody() {
		if i < skipToIndex {
			continue
		}

		tokenType := token.GetType()

		if tokenType == code.Identifier {
			// Extract the statement
			statement := tokenUtil.ExtractTokensBefore(
				function.GetBody()[i:],
				code.Semicolon,
				// Don't handle nested statements here
				false,
				code.Unknown,
				code.Unknown,
			)

			skipToIndex = i + len(statement) + 1

			// Analyze the statement
			AnalyzeStatement(statement, variables, functions)
			continue
		}

		logger.TokenError(
			token,
			"Unexpected token",
			"Expected an identifier or a statement",
			"Check the function body",
		)
	}

	// TODO! Parse return statements
	return object.NothingType
}
