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
	"in":       token.In,
	"to":       token.To,

	// Operators and symbols
	"=":  token.Assign,
	"+":  token.Plus,
	"-":  token.Minus,
	"*":  token.Asterisk,
	"/":  token.Slash,
	"<":  token.LessThan,
	">":  token.GreaterThan,
	"==": token.Equal,
	"!=": token.NotEqual,
	">=": token.GreaterThanOrEqual,
	"<=": token.LessThanOrEqual,
	"&":  token.Ampersand,
	"&&": token.And,
	"||": token.Or,
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

func getKnownToken(entry string) (token.Type, bool) {
	val, ok := knownTokens[entry]

	if !ok {
		return token.Unknown, ok
	}

	return val, ok
}
