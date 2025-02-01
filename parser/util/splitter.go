/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package util

import "fluent/token"

// SplitTokens splits the input tokens into sub-lists based on the separator token.
// It handles nested delimiters by keeping track of the nesting level.
//
// Parameters:
// - input: A slice of pointers to token.Token representing the input tokens.
// - separator: A pointer to token.Token used as the separator for splitting.
// - nestedDelimiters: A token.Type slice representing the nested delimiters.
// - nestedEndDelimiters: A token.Type slice representing the nested end delimiters.
//
// Returns:
// A slice of slices of pointers to token.Token, where each inner slice represents
// a sublist of tokens split by the separator.
func SplitTokens(
	input []token.Token,
	separator token.Type,
	nestedDelimiters []token.Type,
	nestedEndDelimiters []token.Type,
) [][]token.Token {
	// Convert the delimiters to a map for faster lookup
	nestedDelimiterMap := make(map[token.Type]bool)
	nestedEndDelimiterMap := make(map[token.Type]bool)

	for _, d := range nestedDelimiters {
		nestedDelimiterMap[d] = true
	}

	for _, d := range nestedEndDelimiters {
		nestedEndDelimiterMap[d] = true
	}

	result := make([][]token.Token, 0)

	// Used to keep track of the current nesting level
	nestingLevel := 0

	// Used to keep track of the current token list
	currentList := make([]token.Token, 0)

	for _, unit := range input {
		if unit.TokenType == separator && nestingLevel == 0 {
			if len(currentList) == 0 {
				return nil
			}

			result = append(result, currentList)
			currentList = make([]token.Token, 0)
		} else {
			// Handle nested delimiters
			if nestedDelimiterMap[unit.TokenType] {
				nestingLevel++
			} else if nestedEndDelimiterMap[unit.TokenType] {
				nestingLevel--

				if nestingLevel < 0 {
					return nil
				}
			}

			currentList = append(currentList, unit)
		}
	}

	// Append the last list
	result = append(result, currentList)

	return result
}
