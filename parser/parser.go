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

package parser

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/rule/access"
	"fluent/parser/rule/function"
	"fluent/parser/rule/import"
	"fluent/parser/rule/mod"
	"fluent/parser/util"
	"fluent/token"
)

// Parse parses a list of tokens and returns an abstract syntax tree (AST) and an error if any.
//
// Parameters:
//   - tokens: a slice of token.Token representing the tokens to be parsed.
//   - file: a string representing the file name from which the tokens were generated.
//
// Returns:
//   - ast.AST: the resulting abstract syntax tree.
//   - error.Error: an error object if an error occurred during parsing.
func Parse(tokens []token.Token, file string) (ast.AST, error.Error) {
	emptyString := ""
	result := ast.AST{
		Rule:     ast.Program,
		Children: &[]*ast.AST{},
		Value:    &emptyString,
		Line:     0,
		Column:   0,
		File:     &file,
	}

	// Used to skip indexes that have already been processed
	skip := 0

	for i, unit := range tokens {
		if i < skip {
			continue
		}

		tokenType := unit.TokenType

		switch tokenType {
		case token.Import:
			// Extract the tokens before the next semicolon
			extracted := util.ExtractTokensBefore(
				tokens[i:],
				[]token.Type{token.Semicolon},
				false,
				token.Unknown,
				token.Unknown,
				true,
			)

			if extracted == nil {
				return result, error.Error{
					Line:     unit.Line,
					Column:   unit.Column,
					File:     &unit.File,
					Expected: []ast.Rule{ast.StringLiteral},
				}
			}

			skip = i + len(extracted) + 1

			// Parse the import rule
			child, parsingError := _import.ProcessImport(extracted)

			if parsingError.IsError() {
				return result, parsingError
			}

			*result.Children = append(*result.Children, &child)
		case token.Pub, token.Function, token.Mod:
			// Extract the whole block
			extracted := util.ExtractTokensBefore(
				tokens[i:],
				[]token.Type{token.CloseCurly},
				true,
				token.OpenCurly,
				token.CloseCurly,
				true,
			)

			if extracted == nil {
				return result, error.Error{
					Line:     unit.Line,
					Column:   unit.Column,
					File:     &unit.File,
					Expected: []ast.Rule{ast.Function},
				}
			}

			skip = i + len(extracted) + 1

			// Parse the access modifier
			var child ast.AST
			var parsingError error.Error

			if tokenType == token.Pub {
				child, parsingError = access.ProcessPub(extracted)
			} else if tokenType == token.Mod {
				child, parsingError = mod.ProcessMod(extracted, false)
			} else {
				child, parsingError = function.ProcessFunction(extracted, false)
			}

			if parsingError.IsError() {
				return result, parsingError
			}

			*result.Children = append(*result.Children, &child)
		default:
			return ast.AST{}, error.Error{
				Line:     unit.Line,
				Column:   unit.Column,
				File:     &unit.File,
				Expected: []ast.Rule{ast.Import, ast.Function, ast.Module},
			}
		}
	}

	return result, error.Error{}
}
