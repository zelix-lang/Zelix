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
	"bufio"
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
//   - idx: The index at which the line to be printed lives.
//   - line: The line number to be formatted (1-based).
//
// Returns:
//
//	A formatted string with the line number and the line content, or an empty string if the line number is out of range.
func buildLineString(lines []string, idx int, line int) string {
	// Use a builder
	builder := strings.Builder{}

	// Append the line number
	builder.WriteString(strconv.Itoa(line))
	builder.WriteString(" | ")

	// Append the line
	builder.WriteString(lines[idx])

	return builder.String()
}

// getLines extracts specific lines from the content based on the provided line numbers.
// It scans the content line by line and collects the lines that match the target line numbers.
//
// Parameters:
//   - content: A pointer to the string containing the full content.
//   - minimumValue: The minimum line number to start scanning from (1-based).
//   - targetLines: A map of line numbers to be extracted (1-based).
//
// Returns:
//   - A slice of strings containing the extracted lines.
//   - The number of lines that were successfully extracted.
func getLines(content *string, minimumValue int, targetLines map[int]bool) ([]string, int) {
	scanner := bufio.NewScanner(strings.NewReader(*content))
	lineNumber := 1

	elementsInserted := 0
	result := make([]string, 3)

	for scanner.Scan() {
		if lineNumber < minimumValue {
			lineNumber++
			continue
		}

		if elementsInserted == 3 {
			break
		}

		// Push the line
		if targetLines[lineNumber] {
			result[elementsInserted] = scanner.Text()
			elementsInserted++
		}

		lineNumber++
	}

	return result, elementsInserted
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
	lines, insertedLines := getLines(contents, line-1, map[int]bool{
		line - 1: true,
		line:     true,
		line + 1: true,
	})

	// Used to determine where to start
	startAt := 0

	// Print some context
	if insertedLines == 3 {
		startAt++
		fmt.Print("         ")
		// See if line - 1 ends in 9
		if (line-1)%10 == 9 {
			fmt.Print(" ")
		}

		// Build the line string
		fmt.Println(
			fmt.Sprintf(
				"%s%s%s",
				ansi.BrightBlack,
				buildLineString(lines, 0, line-1),
				ansi.Reset,
			),
		)
	}

	// Print an arrow pointing to the column
	fmt.Println(
		fmt.Sprintf(
			"%s       > %s%s",
			highlightColor,
			buildLineString(lines, 1, line),
			ansi.Reset,
		),
	)

	// Used to know if the pinpoint has met at least one character
	pinpointMet := false

	// Print a pinpoint to the column
	targetLine := lines[startAt]
	startAt++

	pinpointStr := strings.Builder{}
	for i := 0; i < len(targetLine); i++ {
		character := targetLine[i]

		if character == ' ' && !pinpointMet {
			pinpointStr.WriteRune(' ')
			continue
		}

		if character != ' ' {
			pinpointMet = true
		}

		if i == column-1 {
			pinpointStr.WriteRune('^')
		} else {
			pinpointStr.WriteRune('-')
		}
	}

	// Print the pinpoint
	fmt.Println(
		fmt.Sprintf(
			"%s       > %s | %s%s",
			highlightColor,
			strconv.Itoa(line),
			pinpointStr.String(),
			ansi.Reset,
		),
	)

	// Print some context
	if insertedLines == 3 {
		// See if line ends in 9
		if line%10 == 9 {
			fmt.Print("        ")
		} else {
			fmt.Print("         ")
		}

		// Build the line string
		fmt.Println(
			ansi.Colorize(ansi.BrightBlack, buildLineString(lines, 2, line+1)),
		)
	}

	fmt.Println(
		fmt.Sprintf(
			"%s         => %s:%s:%s%s",
			ansi.BrightPurple,
			DiscardCwd(*filepath),
			strconv.Itoa(line),
			strconv.Itoa(column),
			ansi.Reset,
		),
	)
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
	logger.Info("Full details:")

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
		message.WriteString(ansi.Colorize(ansi.BrightRed, expected.String()))

		if i < expectedLen-1 {
			message.WriteString(" or ")
		} else if expectedLen > 1 {
			message.WriteString(", ")
		}
	}

	// Return the error message
	return message.String()
}
