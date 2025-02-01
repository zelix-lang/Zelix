/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package lexer

// Error represents a lexical error with details about its location and message.
type Error struct {
	Line    int    // Line number where the error occurred.
	Column  int    // Column number where the error occurred.
	File    string // File name where the error occurred.
	Message string // Description of the error.
}
