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

package pkg

import (
	"fluent/ast"
	"fluent/lexer"
	"fluent/parser/rule/block"
	"fluent/util"
	"os"
)

var invalidStructureMsg = "Invalid package structure"

// ParsePackage reads a file from the given location, lexes its contents,
// parses the tokens into an AST, and extracts package information from the AST.
// It returns a pointer to a Package struct containing the parsed information.
// If any errors occur during reading, lexing, or parsing, the function will
// print the error and exit the program.
func ParsePackage(location string) *Package {
	result := Package{}

	// Read the file at the given location
	data, err := os.ReadFile(location)
	if err != nil {
		panic(err)
	}

	// Lex the file using the fluent lexer
	contents := string(data)
	tokens, lexerErr := lexer.Lex(contents, location)
	if lexerErr != nil {
		// Build and print the error
		util.PrintError(&contents, &location, &lexerErr.Message, lexerErr.Line, lexerErr.Column)
		os.Exit(1)
	}

	// Parse the tokens
	tree, parserErr := block.ProcessBlock(tokens)

	if parserErr != nil {
		errorMessage := util.BuildMessageFromParsingError(*parserErr)

		// Build and print the error
		util.PrintError(&contents, &location, &errorMessage, parserErr.Line, parserErr.Column)
		os.Exit(1)
	}

	// Transverse the AST
	children := *tree.Children

	for _, child := range children {
		switch child.Rule {
		case ast.Assignment:
			// Get the left and right hand side of the assignment
			leftExpr := (*child.Children)[0]
			rightExpr := (*child.Children)[1]
			left := (*leftExpr.Children)[0]
			right := (*rightExpr.Children)[0]

			// Verify rules
			if left.Rule != ast.Identifier {
				// Build and print the error
				util.PrintError(&contents, &location, &invalidStructureMsg, left.Line, left.Column)
				os.Exit(1)
			}

			if right.Rule != ast.StringLiteral {
				// Build and print the error
				util.PrintError(&contents, &location, &invalidStructureMsg, right.Line, right.Column)
				os.Exit(1)
			}

			// Get the value of the left and right hand side
			leftValue := *left.Value
			rightValue := *right.Value

			// Assign the appropriate key
			switch leftValue {
			case "Name":
				result.Name = rightValue
			case "Version":
				result.Version = rightValue
			case "Description":
				result.Description = rightValue
			case "Author":
				result.Author = rightValue
			case "License":
				result.License = rightValue
			case "Entry":
				result.Entry = rightValue
			default:
				// Build and print the error
				util.PrintError(&contents, &location, &invalidStructureMsg, left.Line, left.Column)
				os.Exit(1)
			}
		default:
			// Build and print the error
			util.PrintError(&contents, &location, &invalidStructureMsg, child.Line, child.Column)
			os.Exit(1)
		}
	}

	return &result
}
