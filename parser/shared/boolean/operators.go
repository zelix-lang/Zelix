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

package boolean

import "fluent/token"

// Operators is a list of boolean and comparison operators
var Operators = map[token.Type]struct{}{
	token.And:                {},
	token.Or:                 {},
	token.Equal:              {},
	token.NotEqual:           {},
	token.GreaterThan:        {},
	token.LessThan:           {},
	token.GreaterThanOrEqual: {},
	token.LessThanOrEqual:    {},
}

// OperatorsSlice is a list of boolean and comparison operators as a slice
var OperatorsSlice = []token.Type{
	token.And,
	token.Equal,
	token.NotEqual,
	token.GreaterThan,
	token.LessThan,
	token.GreaterThanOrEqual,
	token.LessThanOrEqual,
}
