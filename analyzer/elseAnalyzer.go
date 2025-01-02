package analyzer

import (
	"fluent/logger"
	"fluent/token"
)

// AnalyzeElse analyzes an else statement
func AnalyzeElse(statement []token.Token) {
	// The statement should have exactly 0 tokens
	// As we receive everything before the curly brace
	// without counting nor the curly brace nor the else keyword
	if len(statement) != 0 {
		logger.TokenError(
			statement[0],
			"Unexpected token",
			"Expected no tokens",
			"Remove the tokens after the 'else' keyword",
		)
	}
}
