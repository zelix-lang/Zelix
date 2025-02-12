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

package format

import "regexp"

type Case int

const (
	// SnakeCase represents the snake case format.
	SnakeCase Case = iota
	// PascalCase represents the pascal case format.
	PascalCase
)

// The snake case regex
var snakeCaseRegex = regexp.MustCompile("^[a-z0-9_]+(_[a-z0-9]+)*$")
var pascalCaseRegex = regexp.MustCompile("^[A-Z][a-z0-9]+([A-Z][a-z0-9]+)*$")

// CheckCase checks if the given string matches the specified case format.
// It returns true if the string matches the format, otherwise false.
//
// Parameters:
//   - str: The string to be checked.
//   - format: The case format to check against (SnakeCase or PascalCase).
//
// Returns:
//   - bool: true if the string matches the specified case format, otherwise false.
func CheckCase(str *string, format Case) bool {
	var regexPattern *regexp.Regexp

	// Update the pattern according to the format
	switch format {
	case SnakeCase:
		regexPattern = snakeCaseRegex
	case PascalCase:
		regexPattern = pascalCaseRegex
	default:
		return false
	}

	// Check if the string matches the regex pattern
	return regexPattern.MatchString(*str)
}
