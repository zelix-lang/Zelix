package core

import (
	"strings"
	"surf/ast"
	"surf/code"
	"surf/core/engine/fun"
	"surf/logger"
	"surf/object"
	runtime2 "surf/runtime"
)

// loadRuntime returns the runtime built-in functions
func loadRuntime() map[string]func(...object.SurfObject) {
	runtime := make(map[string]func(...object.SurfObject))

	// Println
	runtime["impl_write"] = runtime2.Write
	runtime["impl_writeln"] = runtime2.Writeln

	return runtime
}

// Interpret interprets the given file code
func Interpret(fileCode *ast.FileCode, sourceFile string) {
	// Locate the main function
	main, exists, _ := ast.LocateFunction(*fileCode.GetFunctions(), sourceFile, "main")

	if !exists {
		logger.Error(
			"Function 'main' not found",
			"The function 'main' was not found in the file",
			"Add a function named 'main' to the file",
		)
	}

	// Load the runtime functions
	runtime := loadRuntime()

	// Interpret the main function
	fun.CallFun(
		&main,
		runtime,
		fileCode.GetFunctions(),
		*code.NewToken(
			code.Unknown,
			"",
			sourceFile,
			0,
			0,
			"",
			0,
			*new(strings.Builder),
		),
	)
}
