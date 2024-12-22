package tokenUtil

import (
	"zyro/token"
)

// IsValidType checks if a token is a valid data type
func IsValidType(tokenType token.Type) bool {
	switch tokenType {
	case token.Bool:
		return true
	case token.String:
		return true
	case token.Num:
		return true
	case token.Dec:
		return true
	case token.Nothing:
		return true
	default:
		return false
	}
}
