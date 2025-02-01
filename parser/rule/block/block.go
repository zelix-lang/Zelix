/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package block

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/queue"
	"fluent/parser/rule/conditional"
	declaration2 "fluent/parser/rule/declaration"
	"fluent/parser/rule/expression"
	"fluent/parser/rule/loop"
	"fluent/parser/rule/ret"
	"fluent/parser/util"
	"fluent/token"
)

// ProcessBlock processes a block of tokens and returns the corresponding AST and any parsing errors.
// Parameters:
// - input: a slice of tokens to be processed.
// Returns:
// - ast.AST: the abstract syntax tree representing the block.
// - error.Error: any error encountered during parsing.
func ProcessBlock(input []token.Token) (ast.AST, error.Error) {
	if len(input) == 0 {
		return ast.AST{
			Rule:     ast.Block,
			Children: &[]*ast.AST{},
		}, error.Error{}
	}

	// All tokens have been checked before they reach this function
	// Therefore it is not necessary to check the input length here
	result := ast.AST{
		Rule:     ast.Block,
		Line:     input[0].Line,
		Column:   input[0].Column,
		File:     &input[0].File,
		Children: &[]*ast.AST{},
	}

	// Use a queue for nested blocks
	var processQueue []queue.Element

	// Add the first element to the queue
	processQueue = append(processQueue, queue.Element{
		Tokens: input,
		Parent: &result,
	})

	for len(processQueue) > 0 {
		// Pop the first element from the queue
		element := processQueue[0]
		processQueue = processQueue[1:]

		// Subtract the data from the element
		input := element.Tokens
		result := element.Parent

		// Used to skip tokens
		skip := 0

		// Used for parsing ifs
		var lastIf *ast.AST

		for i, el := range input {
			if i < skip {
				continue
			}

			tokenType := el.TokenType

			// Handle conditionals
			switch tokenType {
			case token.If, token.OpenCurly, token.While, token.ElseIf, token.Else, token.For:
				// Extract the tokens in the nested block
				block := util.ExtractTokensBefore(
					input[i:],
					[]token.Type{token.CloseCurly},
					true,
					token.OpenCurly,
					token.CloseCurly,
					true,
				)

				if block == nil {
					return ast.AST{}, error.Error{
						Line:     el.Line,
						Column:   el.Column,
						File:     &el.File,
						Expected: []ast.Rule{ast.Block},
					}
				}

				skip = i + len(block) + 1

				if tokenType == token.OpenCurly {
					// Create a nested block node
					blockNode := ast.AST{
						Rule:     ast.Block,
						Line:     el.Line,
						Column:   el.Column,
						File:     &el.File,
						Children: &[]*ast.AST{},
					}

					// Add the nested block to the result
					*result.Children = append(*result.Children, &blockNode)

					// Schedule the nested block
					processQueue = append(processQueue, queue.Element{
						Tokens: block[1:], // Strip the open curly
						Parent: &blockNode,
					})
				} else {
					isIf := tokenType == token.If
					isElseIf := tokenType == token.ElseIf
					isElse := tokenType == token.Else
					isFor := tokenType == token.For

					if isElseIf || isElse {
						if lastIf == nil {
							return ast.AST{}, error.Error{
								Line:     el.Line,
								Column:   el.Column,
								File:     &el.File,
								Expected: []ast.Rule{ast.If},
							}
						}
					}

					// Parse the condition
					var condition ast.AST
					var parsingError error.Error

					if isFor {
						condition, parsingError = loop.ProcessForLoop(block, &processQueue)
					} else if isElse {
						condition, parsingError = conditional.ProcessElse(block, &processQueue)
					} else {
						condition, parsingError = conditional.ProcessLoopOrConditional(
							block,
							&processQueue,
							isIf,
							isElseIf,
						)
					}

					if parsingError.IsError() {
						return ast.AST{}, parsingError
					}

					if isIf {
						lastIf = &condition
					} else if isElseIf || isElse {
						*lastIf.Children = append(*lastIf.Children, &condition)
					}

					if !isElseIf && !isElse {
						*result.Children = append(*result.Children, &condition)
					}
				}
			default:
				// Extract a statement
				statement := util.ExtractTokensBefore(
					input[i:],
					[]token.Type{token.Semicolon},
					false,
					token.Unknown,
					token.Unknown,
					true,
				)

				if statement == nil {
					return ast.AST{}, error.Error{
						Line:     el.Line,
						Column:   el.Column,
						File:     &el.File,
						Expected: []ast.Rule{ast.Statement},
					}
				}

				skip = i + len(statement) + 1

				var specialCaseNode ast.AST
				var specialCaseError error.Error

				// Process special cases
				if tokenType == token.Let || tokenType == token.Const {
					// Process the declaration
					specialCaseNode, specialCaseError = declaration2.ProcessDeclaration(statement)
				} else if tokenType == token.Continue || tokenType == token.Break {
					// Process the continue or break statement
					specialCaseNode, specialCaseError = loop.ProcessContinueOrBreak(statement)
				} else if tokenType == token.Return {
					// Process the return statement
					specialCaseNode, specialCaseError = ret.ProcessReturn(statement)
				}

				// Check for special cases
				if specialCaseNode.Rule != ast.Program {
					if specialCaseError.IsError() {
						return ast.AST{}, specialCaseError
					}

					// Append the node to the result
					*result.Children = append(*result.Children, &specialCaseNode)

					// Avoid further processing
					continue
				}

				// Parse as an expression
				expr, parsingError := expression.ProcessExpression(statement)

				if parsingError.IsError() {
					return ast.AST{}, parsingError
				}

				*result.Children = append(*result.Children, &expr)
			}
		}
	}

	return result, error.Error{}
}
