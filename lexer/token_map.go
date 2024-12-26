package lexer

import (
	"fluent/token"
)

var knownTokens = map[string]token.Type{
	// Keywords
	"fun":      token.Function,
	"let":      token.Let,
	"const":    token.Const,
	"while":    token.While,
	"for":      token.For,
	"break":    token.Break,
	"continue": token.Continue,
	"if":       token.If,
	"else":     token.Else,
	"elseif":   token.ElseIf,
	"return":   token.Return,
	"mod":      token.Mod,
	"new":      token.New,

	// Operators and symbols
	"=":  token.Assign,
	"+":  token.Plus,
	"-":  token.Minus,
	"++": token.Increment,
	"--": token.Decrement,
	"*":  token.Asterisk,
	"/":  token.Slash,
	"<":  token.LessThan,
	">":  token.GreaterThan,
	"+=": token.AssignAdd,
	"-=": token.AssignSub,
	"/=": token.AssignSlash,
	"*=": token.AssignAsterisk,
	"==": token.Equal,
	"!=": token.NotEqual,
	">=": token.GreaterThanOrEqual,
	"<=": token.LessThanOrEqual,
	"&":  token.Ampersand,
	"|":  token.Bar,
	"^":  token.Xor,
	"!":  token.Not,
	",":  token.Comma,
	";":  token.Semicolon,
	"(":  token.OpenParen,
	")":  token.CloseParen,
	"{":  token.OpenCurly,
	"}":  token.CloseCurly,
	":":  token.Colon,
	"->": token.Arrow,
	"[":  token.OpenBracket,
	"]":  token.CloseBracket,
	".":  token.Dot,
	"%":  token.Percent,

	// Data types
	"str":     token.String,
	"num":     token.Num,
	"dec":     token.Dec,
	"bool":    token.Bool,
	"nothing": token.Nothing,

	// Access modifiers
	"pub": token.Pub,

	// Boolean literals
	"true":  token.BoolLiteral,
	"false": token.BoolLiteral,

	// Imports
	"import": token.Import,
}

func GetKnownToken(entry string) (token.Type, bool) {
	val, ok := knownTokens[entry]

	if !ok {
		return token.Unknown, ok
	}

	return val, ok
}
