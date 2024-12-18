package util

import (
	"strings"
	"surf/ansi"
)

var indicatorBase = ansi.Colorize("magenta_bright", "~")
var indicator = ansi.Colorize("magenta_bright", "^")

// BuildTrace builds the trace of a character in the source code
// in O(n) time
func BuildTrace(
	currentToken strings.Builder,
	currentIndex int,
	input string,
) (string, string) {
	// Use a builder for efficiency
	var trace strings.Builder          // The trace itself
	var traceIndicator strings.Builder // An indicator
	// Both together make something like:
	// pub fn main() { println!("Hello, World!"); }
	// ~~~~~~~~~~~~~~~~^^^^^^^^^^^^^^^^^^^^^^^^^^~~

	startAt := currentIndex - currentToken.Len()

	// Add 25 tokens before the current token for context
	for i := startAt - 25; i < startAt; i++ {
		if i < 0 {
			break
		}

		char := input[i]

		// Skip newlines
		if char == '\n' {
			break
		}

		traceIndicator.WriteString(indicatorBase)
		trace.WriteByte(char)
	}

	// Add the current token
	trace.WriteString(currentToken.String())
	traceIndicator.WriteString(strings.Repeat(indicator, currentToken.Len()))

	// Add 25 tokens after the current token for context
	for i := startAt + currentToken.Len(); i < startAt+currentToken.Len()+25; i++ {
		if i >= len(input) {
			break
		}

		char := input[i]

		// Skip newlines
		if char == '\n' {
			break
		}

		traceIndicator.WriteString(indicatorBase)
		trace.WriteByte(char)
	}

	return strings.Trim(trace.String(), " "), traceIndicator.String()
}
