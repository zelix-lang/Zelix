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

package mod

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/rule/declaration"
	"fluent/parser/rule/function"
	"fluent/parser/rule/identifier"
	"fluent/parser/rule/template"
	"fluent/parser/util"
	"fluent/token"
)

// ProcessMod processes a module from the given tokens.
// It handles public modules, module names, generics, and the module's block content.
//
// Parameters:
//   - input: a slice of tokens representing the module.
//   - public: a boolean indicating if the module is public.
//
// Returns:
//   - ast.AST: the abstract syntax tree representing the module.
//   - error.Error: an error object if any parsing error occurs.
func ProcessMod(input []token.Token, public bool) (*ast.AST, *error.Error) {
	// Check the input's length
	if len(input) < 3 {
		return nil, &error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.Module},
		}
	}

	result := ast.AST{
		Rule:     ast.Module,
		Line:     input[0].Line,
		Column:   input[0].Column,
		File:     &input[0].File,
		Children: &[]*ast.AST{},
	}

	// Append the 'pub' keyword if it is public
	if public {
		pub := ast.AST{
			Rule:     ast.Public,
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Children: &[]*ast.AST{},
		}

		*result.Children = append(*result.Children, &pub)
	}

	// Extract the module name
	moduleName, parsingError := identifier.ProcessIdentifier(&input[1])

	if parsingError != nil {
		return nil, parsingError
	}

	*result.Children = append(*result.Children, moduleName)

	// Used to skip some tokens if needed
	startAt := 2

	if input[startAt].TokenType == token.LessThan {
		// Extract the generics
		genericsRaw := util.ExtractTokensBefore(
			input[startAt:],
			[]token.Type{token.GreaterThan},
			true,
			token.LessThan,
			token.GreaterThan,
			true,
		)

		if genericsRaw == nil {
			return nil, &error.Error{
				Line:     input[startAt].Line,
				Column:   input[startAt].Column,
				File:     &input[startAt].File,
				Expected: []ast.Rule{ast.Generics},
			}
		}

		// Skip all the generics
		startAt += len(genericsRaw) + 1

		// Pass the input to the template parser
		generics, parsingError := template.ProcessTemplates(genericsRaw)

		if parsingError != nil {
			return nil, parsingError
		}

		// Append the generics to the result's children
		*result.Children = append(*result.Children, generics)
	}

	if input[startAt].TokenType != token.OpenCurly {
		return nil, &error.Error{
			Line:     input[startAt].Line,
			Column:   input[startAt].Column,
			File:     &input[startAt].File,
			Expected: []ast.Rule{ast.Module},
		}
	}

	// Extract the whole block
	block := input[startAt+1:]

	// Create a new block AST
	blockAST := ast.AST{
		Rule:     ast.Block,
		Line:     input[0].Line,
		Column:   input[0].Column,
		File:     &input[0].File,
		Children: &[]*ast.AST{},
	}

	// Append the block to the result's children
	*result.Children = append(*result.Children, &blockAST)

	// Unlike traditional blocks, a module's block allows solely
	// incomplete declarations and functions
	// therefore, the block parser implementation is different
	skip := 0

	for i, unit := range block {
		if i < skip {
			continue
		}

		tokenType := unit.TokenType

		switch tokenType {
		case token.Let, token.Const:
			// Extract the tokens before the next semicolon
			extracted := util.ExtractTokensBefore(
				block[i:],
				[]token.Type{token.Semicolon},
				false,
				token.Unknown,
				token.Unknown,
				true,
			)

			if extracted == nil {
				return nil, &error.Error{
					Line:     unit.Line,
					Column:   unit.Column,
					File:     &unit.File,
					Expected: []ast.Rule{ast.Declaration},
				}
			}

			skip = i + len(extracted) + 1

			// Pass the input to the incomplete declaration parser
			var child *ast.AST
			var parsingError *error.Error

			if util.TokenSliceContains(extracted, map[token.Type]struct{}{token.Assign: {}}) {
				child, parsingError = declaration.ProcessDeclaration(extracted)
			} else {
				child, parsingError = declaration.ProcessIncompleteDeclaration(extracted)
			}

			if parsingError != nil {
				return nil, parsingError
			}

			// Append the declaration to the block's children
			*blockAST.Children = append(*blockAST.Children, child)
		case token.Pub, token.Function:
			// Extract the whole block
			extracted := util.ExtractTokensBefore(
				block[i:],
				[]token.Type{token.CloseCurly},
				true,
				token.OpenCurly,
				token.CloseCurly,
				true,
			)

			if extracted == nil {
				return nil, &error.Error{
					Line:     unit.Line,
					Column:   unit.Column,
					File:     &unit.File,
					Expected: []ast.Rule{ast.Function},
				}
			}

			skip = i + len(extracted) + 1
			isPublic := tokenType == token.Pub

			if isPublic {
				// Strip the 'pub' keyword
				extracted = extracted[1:]
			}

			// Pass the input to the function parser
			child, parsingError := function.ProcessFunction(extracted, isPublic)

			if parsingError != nil {
				return nil, parsingError
			}

			// Append the function to the block's children
			*blockAST.Children = append(*blockAST.Children, child)
		default:
			// Invalid token
			return nil, &error.Error{
				Line:     unit.Line,
				Column:   unit.Column,
				File:     &unit.File,
				Expected: []ast.Rule{ast.IncompleteDeclaration, ast.Function, ast.Public},
			}
		}
	}

	return &result, nil
}
