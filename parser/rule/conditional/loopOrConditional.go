package conditional

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/queue"
	"fluent/parser/rule/expression"
	"fluent/parser/util"
	"fluent/token"
)

// ProcessLoopOrConditional processes a loop or conditional statement based on the provided tokens.
// It determines whether the statement is an `if` or `while` based on the `isIf` parameter.
// The function returns an AST (Abstract Syntax Tree) and an error if any issues are encountered during parsing.
//
// Parameters:
// - input: A slice of tokens representing the input code.
// - exprQueue: A pointer to a slice of queue elements used for expression parsing.
// - isIf: A boolean indicating whether the statement is an `if` (true) or `while` (false).
//
// Returns:
// - ast.AST: The resulting abstract syntax tree.
// - error.Error: An error object if any parsing errors occur.
func ProcessLoopOrConditional(
	input []token.Token,
	exprQueue *[]queue.Element,
	// The logic for the while loop's syntax implementation
	// is exactly the same as the one in this function
	// therefore, it is re-utilized
	isIf bool,
	isElseIf bool,
) (*ast.AST, *error.Error) {
	var rule ast.Rule

	if isIf {
		rule = ast.If
	} else if isElseIf {
		rule = ast.ElseIf
	} else {
		rule = ast.While
	}

	result := ast.AST{
		Rule:     rule,
		Line:     input[0].Line,
		Column:   input[0].Column,
		File:     &input[0].File,
		Children: &[]*ast.AST{},
	}

	// Check the input length
	if len(input) < 3 {
		return nil, &error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{rule},
		}
	}

	// Strip the first token
	input = input[1:]

	// Extract everything before the first opening brace
	// to find the condition
	condition := util.ExtractTokensBefore(
		input,
		[]token.Type{token.OpenCurly},
		false,
		token.Unknown,
		token.Unknown,
		true,
	)

	if condition == nil {
		return nil, &error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{rule},
		}
	}

	// Create a new condition node and append it
	conditionNode, parsingError := expression.ProcessExpression(condition)

	if parsingError != nil {
		return nil, parsingError
	}

	*result.Children = append(*result.Children, conditionNode)

	// Extract the block
	var block []token.Token

	// See if there is a block
	if len(input) > 1+len(condition) {
		block = input[1+len(condition):]
	} else {
		block = make([]token.Token, 0)
	}

	// Create a new block node
	blockNode := ast.AST{
		Rule:     ast.Block,
		Line:     input[0].Line,
		Column:   input[0].Column,
		Children: &[]*ast.AST{},
	}

	// Append the block to the result
	*result.Children = append(*result.Children, &blockNode)

	// Queue the block for parsing
	*exprQueue = append(*exprQueue, queue.Element{
		Tokens: block,
		Parent: &blockNode,
	})

	return &result, nil
}
