package tokenUtil

import "surf/code"

// IsValidType checks if a token is a valid data type
func IsValidType(tokenType code.TokenType) bool {
	switch tokenType {
	case code.Bool:
		return true
	case code.String:
		return true
	case code.Num:
		return true
	case code.Dec:
		return true
	case code.Nothing:
		return true
	default:
		return false
	}
}
