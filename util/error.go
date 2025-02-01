/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package util

import (
	"fluent/ansi"
	"fluent/logger"
	"fluent/parser/error"
	"fmt"
	"strconv"
	"strings"
)

// buildLineString constructs a string representation of a specific line from the source code.
// It includes the line number and the line content.
//
// Parameters:
//   - lines: A slice of strings representing the lines of the source code.
//   - line: The line number to be formatted (1-based).
//
// Returns:
//
//	A formatted string with the line number and the line content, or an empty string if the line number is out of range.
func buildLineString(lines []string, line int) string {
	if line-1 < 0 || line-1 >= len(lines) {
		return ""
	}

	result := ""

	// Append the line number
	result += strconv.Itoa(line)
	result += " | "

	// Append the line
	result += lines[line-1]

	return result
}

// BuildAndPrintDetails prints detailed information about a specific line and column in the source code.
// It provides context around the error location and highlights the exact position.
//
// Parameters:
//   - contents: The full contents of the source code as a string.
//   - filepath: The path to the file where the error occurred.
//   - line: The line number where the error occurred (1-based).
//   - column: The column number where the error occurred (1-based).
func BuildAndPrintDetails(contents, filepath *string, line int, column int, isError bool) {
	var highlightColor string

	if isError {
		highlightColor = ansi.BrightRed
	} else {
		highlightColor = ansi.BrightYellow
	}

	// Split the contents by lines
	lines := strings.Split(*contents, "\n")

	// Print some context
	if line > 1 {
		print("         ")
		// See if line - 1 ends in 9
		if (line-1)%10 == 9 {
			print(" ")
		}

		// Build the line string
		fmt.Print(ansi.Colorize(ansi.BrightBlack, buildLineString(lines, line-1)))
		fmt.Println()
	}

	// Print the line with the error
	print("       ")
	// Print an arrow pointing to the column
	fmt.Print(ansi.Colorize(highlightColor, "> ").Bold())
	fmt.Print(ansi.Colorize(highlightColor, buildLineString(lines, line)).Bold())
	fmt.Println()

	print("       ")
	// Print an arrow pointing to the column
	fmt.Print(ansi.Colorize(highlightColor, "> ").Bold())
	fmt.Print(ansi.Colorize(highlightColor, strconv.Itoa(line)).Bold())
	fmt.Print(ansi.Colorize(highlightColor, " | ").Bold())

	// Used to know if the pinpoint has met at least one character
	pinpointMet := false

	// Print a pinpoint to the column
	for i := 0; i < len(lines[line-1]); i++ {
		character := lines[line-1][i]

		if character == ' ' && !pinpointMet {
			print(" ")
			continue
		}

		if character != ' ' {
			pinpointMet = true
		}

		if i == column-1 {
			fmt.Print(ansi.Colorize(highlightColor, "^").Bold())
		} else {
			fmt.Print(ansi.Colorize(ansi.BrightBlack, "-").Bold())
		}
	}

	println()

	// Print some context
	if line < len(lines) {
		// See if line ends in 9
		if line%10 == 9 {
			print("        ")
		} else {
			print("         ")
		}

		// Build the line string
		fmt.Print(ansi.Colorize(ansi.BrightBlack, buildLineString(lines, line+1)))
		println()
	}

	print("         ")
	fmt.Print(ansi.Colorize(ansi.BrightPurple, "=> "))
	fmt.Print(ansi.Colorize(ansi.BrightBlue, DiscardCwd(*filepath)))
	fmt.Print(ansi.Colorize(ansi.BrightPurple, ":"+strconv.Itoa(line)+":"+strconv.Itoa(column)))
	println()
}

// PrintError prints a formatted error message with context from the source code.
// It highlights the line and column where the error occurred and provides additional context.
//
// Parameters:
//   - contents: The full contents of the source code as a string.
//   - filepath: The path to the file where the error occurred.
//   - message: The error message to be displayed.
//   - line: The line number where the error occurred (1-based).
//   - column: The column number where the error occurred (1-based).
func PrintError(contents, filepath, message *string, line int, column int) {
	logger.Error(*message)
	logger.Info("Full details")

	BuildAndPrintDetails(contents, filepath, line, column, true)
}

// BuildMessageFromParsingError constructs a detailed error message from a parsing error.
// It formats the expected tokens in a human-readable way.
//
// Parameters:
//   - parsingError: The parsing error containing the expected tokens.
//
// Returns:
//
//	A formatted string describing the expected tokens.
func BuildMessageFromParsingError(parsingError error.Error) string {
	// Use a strings.Builder for efficiency
	var message strings.Builder

	message.WriteString("Expected ")

	expectedLen := len(parsingError.Expected)
	for i, expected := range parsingError.Expected {
		message.WriteString(ansi.Colorize(ansi.BrightRed, expected.String()).Bold().String())

		if i < expectedLen-1 {
			message.WriteString(" or ")
		} else if expectedLen > 1 {
			message.WriteString(", ")
		}
	}

	// Return the error message
	return message.String()
}
