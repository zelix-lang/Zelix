package lexer

import "surf/code"

var knownTokens = map[string]code.TokenType{
	// Keywords
	"fun":      code.Function,
	"let":      code.Let,
	"const":    code.Const,
	"while":    code.While,
	"for":      code.For,
	"break":    code.Break,
	"continue": code.Continue,
	"if":       code.If,
	"else":     code.Else,
	"elseif":   code.ElseIf,
	"return":   code.Return,

	// Operators and symbols
	"=":  code.Assign,
	"+":  code.Plus,
	"-":  code.Minus,
	"++": code.Increment,
	"--": code.Decrement,
	"*":  code.Asterisk,
	"/":  code.Slash,
	"<":  code.LessThan,
	">":  code.GreaterThan,
	"+=": code.AssignAdd,
	"-=": code.AssignSub,
	"/=": code.AssignSlash,
	"*=": code.AssignAsterisk,
	"==": code.Equal,
	"!=": code.NotEqual,
	">=": code.GreaterThanOrEqual,
	"<=": code.LessThanOrEqual,
	"&":  code.Ampersand,
	"|":  code.Bar,
	"^":  code.Xor,
	"!":  code.Not,
	",":  code.Comma,
	";":  code.Semicolon,
	"(":  code.OpenParen,
	")":  code.CloseParen,
	"{":  code.OpenCurly,
	"}":  code.CloseCurly,
	":":  code.Colon,
	"->": code.Arrow,
	"[":  code.OpenBracket,
	"]":  code.CloseBracket,
	".":  code.Dot,
	"%":  code.Percent,

	// Data types
	"str":     code.String,
	"num":     code.Num,
	"dec":	   code.Dec,
	"bool":    code.Bool,
	"nothing": code.Nothing,

	// Access modifiers
	"pub": code.Pub,

	// Boolean literals
	"true":  code.BoolLiteral,
	"false": code.BoolLiteral,

	// Imports
	"import": code.Import,
}

func GetKnownToken(entry string) (code.TokenType, bool) {
	val, ok := knownTokens[entry]

	if !ok {
		return code.Unknown, ok
	}

	return val, ok
}
