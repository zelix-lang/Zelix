package marshall

import (
	"fluent/token"
)

// TokenTypeToString converts a token type to a string
func TokenTypeToString(_type token.Type) string {
	switch _type {
	case token.Identifier:
		return "identifier"
	case token.String:
		return "string"
	case token.Num, token.NumLiteral:
		return "number"
	case token.Colon:
		return "colon"
	case token.Comma:
		return "comma"
	case token.Semicolon:
		return "semicolon"
	case token.OpenParen:
		return "left parenthesis"
	case token.CloseParen:
		return "right parenthesis"
	case token.OpenCurly:
		return "left brace"
	case token.CloseCurly:
		return "right brace"
	case token.OpenBracket:
		return "left bracket"
	case token.CloseBracket:
		return "right bracket"
	case token.Plus:
		return "plus"
	case token.Minus:
		return "minus"
	case token.Asterisk:
		return "asterisk"
	case token.Slash:
		return "slash"
	case token.Assign:
		return "assign"
	case token.Equal:
		return "equal"
	default:
		return "not yet implemented"
	}
}
