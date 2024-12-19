package tokenUtil

import "surf/code"

// TokenTypeToString converts a token type to a string
func TokenTypeToString(_type code.TokenType) string {
	switch _type {
	case code.Identifier:
		return "identifier"
	case code.String:
		return "string"
	case code.Num, code.NumLiteral:
		return "number"
	case code.Colon:
		return "colon"
	case code.Comma:
		return "comma"
	case code.Semicolon:
		return "semicolon"
	case code.OpenParen:
		return "left parenthesis"
	case code.CloseParen:
		return "right parenthesis"
	case code.OpenCurly:
		return "left brace"
	case code.CloseCurly:
		return "right brace"
	case code.OpenBracket:
		return "left bracket"
	case code.CloseBracket:
		return "right bracket"
	case code.Plus:
		return "plus"
	case code.Minus:
		return "minus"
	case code.Asterisk:
		return "asterisk"
	case code.Slash:
		return "slash"
	case code.Assign:
		return "assign"
	case code.Equal:
		return "equal"
	default:
		return "not yet implemented"
	}
}
