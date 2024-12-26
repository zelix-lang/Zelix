package analyzer

import (
	"fluent/ansi"
	"fluent/ast"
	"fluent/code/types"
	"fluent/code/wrapper"
	"fluent/logger"
	"fluent/stack"
	"os"
	"regexp"
)

// A regex to match camelCase variable names
var snakeCaseRegex = regexp.MustCompile("^[a-z]+(_[a-z0-9]+)*$")

// A dummy "nothing" type wrapper
var dummyNothingType = wrapper.ForceNewTypeWrapper(
	"dummy-nothing",
	[]wrapper.TypeWrapper{},
	types.NothingType,
)

// AnalyzeFileCode analyzes the given file code
func AnalyzeFileCode(code *ast.FileCode, source string) {
	fileFunctions, found := (*code.GetFunctions())[source]

	if !found {
		logger.Error("No functions in the file")
		logger.Log("The file does not contain any functions")
		logger.Help("Add a function to the file")

		os.Exit(1)
	}

	mainFunction, found := fileFunctions["main"]

	if !found {
		logger.Error("No main function")
		logger.Log("The file does not contain a main function")
		logger.Help("Add a function named 'main' to the file")

		os.Exit(1)
	}

	// Analyze the main function
	AnalyzeMainFunc(mainFunction)

	for _, mods := range *code.GetModules() {
		for _, mod := range mods {
			AnalyzeMod(
				*mod,
				code.GetFunctions(),
				code.GetModules(),
			)
		}
	}

	// Analyze all other functions
	for _, functions := range *code.GetFunctions() {
		for name, function := range functions {
			if !snakeCaseRegex.MatchString(name) {
				logger.TokenWarning(
					function.GetTrace(),
					"Function name is not in snake_case",
					"Fluent uses snake_case for functions' names",
					"Check "+ansi.Colorize("yellow", "[U-002]")+" in the style guide",
				)
			}

			// During start phase, argument checking is not necessary
			AnalyzeFun(
				function,
				code.GetFunctions(),
				code.GetModules(),
				function.GetTrace(),
				false,
				stack.NewStack(),
			)
		}
	}
}
