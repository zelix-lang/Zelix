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
	AnalyzeFun(mainFunction, code.GetFunctions(), mainFunction.GetTrace())

	// Analyze all other functions
	for file, functions := range *code.GetFunctions() {
		for name, function := range functions {
			if file == source && name == "main" {
				continue
			}

			AnalyzeFun(function, code.GetFunctions(), function.GetTrace())
		}
	}
}
