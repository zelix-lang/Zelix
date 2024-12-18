package main

import (
	"fmt"
	"os"
	"surf/ast"
	"surf/checker"
	"surf/core"
	"surf/lexer"
	"time"
)

func main() {
	file := os.Getenv("SURF_DUMMY_FILE")

	// Read the file
	input, err := os.ReadFile(file)

	if err != nil {
		panic(err)
	}

	start := time.Now()

	// Lex the input
	tokens := lexer.Lex(string(input), file)

	fmt.Println("Lexed in:", time.Since(start))

	start = time.Now()
	// Parse the tokens into a FileCode
	fileCode := ast.Parse(tokens)
	fmt.Println("Wrapped to AST in:", time.Since(start))

	start = time.Now()
	// Analyze the file code
	checker.AnalyzeFileCode(*fileCode, file)
	fmt.Println("Analyzed in:", time.Since(start))

	start = time.Now()
	// Interpret the file code
	core.Interpret(fileCode, file)
	fmt.Println("Interpreted in:", time.Since(start))

}
