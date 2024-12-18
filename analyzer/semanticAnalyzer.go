package analyzer

import (
	"surf/ast"
	"surf/logger"
)

// AnalyzeFileCode analyzes the given file code
func AnalyzeFileCode(code *ast.FileCode, source string) {
	fileFunctions, found := (*code.GetFunctions())[source]

	if !found {
		logger.Error("No functions in the file")
		logger.Log("The file does not contain any functions")
		logger.Help("Add a function to the file")
	}

	mainFunction, found := fileFunctions["main"]

	if !found {
		logger.Error("No main function")
		logger.Log("The file does not contain a main function")
		logger.Help("Add a function named 'main' to the file")
	}

	// Analyze the main function
	AnalyzeMainFunc(mainFunction)

	// Analyze all other functions
	for _, functions := range *code.GetFunctions() {
		for _, function := range functions {
			// During start phase, argument checking is not necessary
			AnalyzeFun(function, code.GetFunctions(), function.GetTrace(), false)
		}
	}
}
