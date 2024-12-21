package args

import (
	"surf/token"
	"surf/tokenUtil"
)

// SplitArgs splits the arguments into a slice of slices of tokens
// and an int representing the number of tokens in the
//
//	function invocation
func SplitArgs(
	statement []token.Token,
) ([][]token.Token, int) {
	// Get the parameters
	// Skip the first 2 tokens (the function name and the opening parenthesis)
	// and the last token (the closing parenthesis)
	parametersRaw := statement[2:]

	// Split by commas
	return tokenUtil.SplitTokens(
		parametersRaw,
		token.Comma,
		token.OpenParen,
		token.CloseParen,
	), len(statement) + 1
}
