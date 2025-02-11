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

package arithmetic

import "fluent/token"

// Operators is a list of arithmetic operators
var Operators = map[token.Type]struct{}{
	token.Plus:     {},
	token.Minus:    {},
	token.Slash:    {},
	token.Asterisk: {},
}

// OperatorsSlice is a list of arithmetic operators as a slice
var OperatorsSlice = []token.Type{
	token.Plus,
	token.Minus,
	token.Slash,
	token.Asterisk,
}
