/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package logger

import (
	"fluent/ansi"
	"fmt"
	"strings"
)

// Prefixes for the different log levels
var errorPrefix = ansi.Colorize(ansi.BoldBrightRed, "[ERROR] ")
var infoPrefix = ansi.Colorize(ansi.BoldBrightBlue, "[INFO]  ")
var warnPrefix = ansi.Colorize(ansi.BoldBrightYellow, "[WARN]  ")
var helpPrefix = ansi.Colorize(ansi.BoldBrightGreen, "[HELP]  ")

// buildMessageImpl constructs a formatted log message.
// Parameters:
//   - prefix: the prefix to be added to each message.
//   - color: the color code to be used for colorizing the message.
//   - colorize: a boolean indicating whether to colorize the message.
//   - message: variadic string arguments representing the log messages.
//
// Returns:
//
//	A single string containing the formatted log messages.
func buildMessageImpl(prefix string, color string, colorize bool, message ...string) string {
	// Use a string builder
	builder := strings.Builder{}

	for _, m := range message {
		builder.WriteString(prefix)
		if colorize {
			builder.WriteString(ansi.Colorize(color, m))
		} else {
			builder.WriteString(m)
		}

		// Write a newline
		builder.WriteString("\n")
	}

	return builder.String()
}

// Warn prints the provided messages with a warning prefix.
// Parameters:
//   - message: variadic string arguments to be logged.
func Warn(message ...string) {
	fmt.Print(BuildWarn(message...))
}

// Info prints the provided messages with an info prefix.
// Parameters:
//   - message: variadic string arguments to be logged.
func Info(message ...string) {
	fmt.Print(BuildInfo(message...))
}

// Help prints the provided messages with a help prefix.
// Parameters:
//   - message: variadic string arguments to be logged.
func Help(message ...string) {
	fmt.Print(BuildHelp(message...))
}

// Error prints the provided messages with an error prefix.
// Parameters:
//   - message: variadic string arguments to be logged.
func Error(message ...string) {
	fmt.Print(BuildError(message...))
}

// BuildInfo constructs an info log message.
// Parameters:
//   - message: variadic string arguments representing the log messages.
//
// Returns:
//
//	A single string containing the formatted log messages.
func BuildInfo(message ...string) string {
	return buildMessageImpl(infoPrefix, "", false, message...)
}

// BuildWarn constructs a warning log message.
// Parameters:
//   - message: variadic string arguments representing the log messages.
//
// Returns:
//
//	A single string containing the formatted log messages.
func BuildWarn(message ...string) string {
	return buildMessageImpl(warnPrefix, ansi.BrightYellow, true, message...)
}

// BuildHelp constructs a help log message.
// Parameters:
//   - message: variadic string arguments representing the log messages.
//
// Returns:
//
//	A single string containing the formatted log messages.
func BuildHelp(message ...string) string {
	return buildMessageImpl(helpPrefix, "", false, message...)
}

// BuildError constructs an error log message.
// Parameters:
//   - message: variadic string arguments representing the log messages.
//
// Returns:
//
//	A single string containing the formatted log messages.
func BuildError(message ...string) string {
	return buildMessageImpl(errorPrefix, ansi.BrightRed, true, message...)
}
