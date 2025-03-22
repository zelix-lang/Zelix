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

package call

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/queue"
	"fluent/parser/rule/identifier"
	"fluent/parser/util"
	"fluent/token"
)

// ProcessCallArguments processes the arguments of a function call from the given input tokens and updates the expression queue.
// It returns an AST node representing the parameters and an error if the parsing fails.
//
// Parameters:
// - input: A slice of tokens representing the function call arguments.
// - expressionQueue: A pointer to a slice of queue elements to be updated with parsed parameters.
//
// Returns:
// - ast.AST: The AST node representing the parameters.
// - error.Error: An error object if the parsing fails.
func ProcessCallArguments(
	input []token.Token,
	expressionQueue *[]queue.Element,
) (*ast.AST, *error.Error) {
	// Check the 2nd token
	if input[0].TokenType != token.OpenParen {
		return nil, &error.Error{
			Line:     input[1].Line,
			Column:   input[1].Column,
			File:     &input[1].File,
			Expected: []ast.Rule{ast.FunctionCall},
		}
	}

	// Check the last token
	if input[len(input)-1].TokenType != token.CloseParen {
		return nil, &error.Error{
			Line:     input[len(input)-1].Line,
			Column:   input[len(input)-1].Column,
			File:     &input[len(input)-1].File,
			Expected: []ast.Rule{ast.FunctionCall},
		}
	}

	// Check if there are parameters
	if len(input) < 3 {
		return &ast.AST{}, nil
	}

	// Strip the first 1 token and the last token to get the parameters
	parametersRaw := input[1 : len(input)-1]
	// Split the parameters to get the individual ones
	parameters := util.SplitTokens(
		parametersRaw,
		token.Comma,
		[]token.Type{token.OpenParen, token.LessThan, token.OpenBracket},
		[]token.Type{token.CloseParen, token.GreaterThan, token.CloseBracket},
	)

	// Create a parameters node
	parametersNode := ast.AST{
		Rule:     ast.Parameters,
		Line:     input[1].Line,
		Column:   input[1].Column,
		File:     &input[1].File,
		Children: &[]*ast.AST{},
	}

	// Iterate over all parameters to parse them
	for _, parameter := range parameters {
		// Create a new expression node for the parameter
		expressionNode := ast.AST{
			Rule:     ast.Expression,
			Line:     parameter[0].Line,
			Column:   parameter[0].Column,
			File:     &parameter[0].File,
			Children: &[]*ast.AST{},
		}

		// Create a new parameter node
		parameterNode := ast.AST{
			Rule:     ast.Parameter,
			Line:     parameter[0].Line,
			Column:   parameter[0].Column,
			File:     &parameter[0].File,
			Children: &[]*ast.AST{&expressionNode},
		}

		// Add the parameter to the expression queue
		*expressionQueue = append(*expressionQueue, queue.Element{Tokens: parameter, Parent: &expressionNode})

		// Append the parameter to the parameters node
		*parametersNode.Children = append(*parametersNode.Children, &parameterNode)
	}

	// Append the parameters to the result
	return &parametersNode, nil
}

// ProcessFunctionCall processes a function call from the given input tokens and updates the expression queue.
// It returns an AST node representing the function call and an error if the parsing fails.
//
// Parameters:
// - input: A slice of tokens representing the function call.
// - expressionQueue: A pointer to a slice of queue elements to be updated with parsed parameters.
//
// Returns:
// - ast.AST: The AST node representing the function call.
// - error.Error: An error object if the parsing fails.
func ProcessFunctionCall(input []token.Token, expressionQueue *[]queue.Element) (*ast.AST, *error.Error) {
	// Check the input length
	if len(input) < 3 {
		return nil, &error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.FunctionCall},
		}
	}

	// Check the 1st token
	functionName, parsingError := identifier.ProcessIdentifier(&input[0])

	if parsingError != nil {
		return nil, parsingError
	}

	// Create a result node
	result := ast.AST{
		Rule:     ast.FunctionCall,
		Line:     input[0].Line,
		Column:   input[0].Column,
		File:     &input[0].File,
		Children: &[]*ast.AST{},
	}

	// Append the function name to the result
	*result.Children = append(*result.Children, functionName)

	// Parse the call arguments
	callArguments, parsingError := ProcessCallArguments(input[1:], expressionQueue)

	if parsingError != nil {
		return nil, parsingError
	}

	if callArguments.Rule != ast.Program {
		// Append the call arguments to the function call
		*result.Children = append(*result.Children, callArguments)
	}

	return &result, nil
}
