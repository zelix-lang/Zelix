package util

import (
	"strings"
)

// BuildTrace builds the trace of a character in the source code
// in O(n) time
func BuildTrace(
	currentToken strings.Builder,
	currentIndex int,
	input string,
) string {
	// Use a builder for efficiency
	var trace strings.Builder // The trace itself

	startAt := currentIndex - currentToken.Len()

	// Add 25 tokens before the current token for context
	for i := startAt; i > startAt-25; i-- {
		if i < 0 {
			break
		}

		char := input[i]

		// Skip newlines
		if char == '\n' {
			break
		}

		trace.WriteByte(char)
	}

	// Invert the trace
	traceStr := trace.String()
	trace.Reset()

	for i := len(traceStr) - 1; i >= 0; i-- {
		trace.WriteByte(traceStr[i])
	}

	// Add the current token
	trace.WriteString(currentToken.String())

	// Add 25 tokens after the current token for context
	for i := startAt + currentToken.Len() + 1; i < startAt+currentToken.Len()+25; i++ {
		if i >= len(input) {
			break
		}

		char := input[i]

		// Skip newlines
		if char == '\n' {
			break
		}

		trace.WriteByte(char)
	}

	return strings.TrimSpace(trace.String())
}
