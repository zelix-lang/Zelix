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

package util

import "fluent/token"

// TokenSliceContains checks if any token in the slice is present in the elements map.
//
// Parameters:
//   - slice: a slice of token.Token to search through.
//   - elements: a map of token.Type to an empty struct, representing the elements to check for.
//
// Returns:
//   - bool: true if any token in the slice is found in the elements map, false otherwise.
func TokenSliceContains(slice []token.Token, elements map[token.Type]struct{}) bool {
	for _, e := range slice {
		if _, ok := elements[e.TokenType]; ok {
			return true
		}
	}
	return false
}
