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

package token

// Token represents a lexical token with its type, value, and position in the source file.
type Token struct {
	TokenType Type   // The type of the token.
	Value     string // The literal value of the token.
	File      string // The source file where the token was found.
	Line      int    // The line number in the source file.
	Column    int    // The column number in the source file.
}

// NewToken creates a new Token with the given type, value, file, line, and column.
func NewToken(
	tokenType Type,
	value string,
	file string,
	line int,
	column int,
) Token {
	return Token{
		tokenType,
		value,
		file,
		line,
		column,
	}
}
