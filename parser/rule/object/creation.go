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

package object

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/queue"
	"fluent/parser/rule/call"
	"fluent/parser/rule/generic"
	"fluent/parser/rule/identifier"
	"fluent/parser/util"
	"fluent/token"
)

// ProcessObjectCreation processes the creation of an object in the Fluent programming language.
// It takes a slice of tokens and a pointer to an expression queue, and returns an AST and an error.
//
// Parameters:
//   - input: a slice of tokens representing the object creation statement.
//   - expressionQueue: a pointer to a queue of elements used for processing expressions.
//
// Returns:
//   - ast.AST: the abstract syntax tree representing the object creation.
//   - error.Error: an error object containing details if an error occurs during processing.
func ProcessObjectCreation(input []token.Token, expressionQueue *[]queue.Element) (ast.AST, error.Error) {
	// Check the input's length
	if len(input) < 4 {
		return ast.AST{}, error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.ObjectCreation},
		}
	}

	result := ast.AST{
		Rule:     ast.ObjectCreation,
		File:     &input[0].File,
		Line:     input[0].Line,
		Column:   input[0].Column,
		Children: &[]*ast.AST{},
	}

	// The first token (new) is omitted, it has already been validated
	// The second token is the object's name, it must be an identifier
	id, parsingError := identifier.ProcessIdentifier(&input[1])

	if parsingError.IsError() {
		return ast.AST{}, parsingError
	}

	// Append the identifier to the result's children
	*result.Children = append(*result.Children, &id)

	// Used to skip some tokens if needed
	startAt := 2

	// Check for generics
	if input[startAt].TokenType == token.LessThan {
		// Extract all generics
		genericsRaw := util.ExtractTokensBefore(
			input[startAt:],
			[]token.Type{token.GreaterThan},
			true,
			token.LessThan,
			token.GreaterThan,
			true,
		)

		if genericsRaw == nil {
			return ast.AST{}, error.Error{
				Line:     input[startAt].Line,
				Column:   input[startAt].Column,
				File:     &input[startAt].File,
				Expected: []ast.Rule{ast.Generics},
			}
		}

		// Pass the tokens to the generics parser
		generics, parsingError := generic.ProcessGenerics(genericsRaw[1:])

		if parsingError.IsError() {
			return ast.AST{}, parsingError
		}

		// Append the generics to the result's children
		*result.Children = append(*result.Children, &generics)

		// Skip the generics
		startAt += len(genericsRaw) + 1
	}

	// Parse arguments
	arguments, parsingError := call.ProcessCallArguments(input[startAt:], expressionQueue)

	if parsingError.IsError() {
		return ast.AST{}, parsingError
	}

	if arguments.Rule != ast.Program {
		// Append the arguments to the result's children
		*result.Children = append(*result.Children, &arguments)
	}

	return result, error.Error{}
}
