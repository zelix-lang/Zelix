package token

import (
	"fluent/util"
	"strconv"
	"strings"
)

// Token represents a small piece of the source code
type Token struct {
	tokenType    Type
	value        string
	file         string
	line         int
	column       int
	trace        string
	traceContext string
}

func NewToken(
	tokenType Type,
	value string,
	file string,
	line int,
	column int,
	input string,
	currentIndex int,
	currentToken strings.Builder,
) *Token {
	context := util.BuildTrace(currentToken, currentIndex, input)

	return &Token{
		tokenType,
		value,
		file,
		line,
		column,
		"At " + file + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(column),
		context,
	}
}

// String returns the string representation of the Token.
func (t Token) String() string {
	return t.value
}

// GetType returns the type of the Token.
func (t Token) GetType() Type {
	return t.tokenType
}

// GetValue returns the value of the Token.
func (t Token) GetValue() string {
	return t.value
}

// GetFile returns the file name where the Token is located.
func (t Token) GetFile() string {
	return t.file
}

// GetLine returns the line number where the Token is located.
func (t Token) GetLine() int {
	return t.line
}

// GetColumn returns the column number where the Token is located.
func (t Token) GetColumn() int {
	return t.column
}

// GetTrace returns the trace information of the Token.
func (t Token) GetTrace() string {
	return t.trace
}

// GetTraceContext returns the trace context of the Token.
func (t Token) GetTraceContext() string {
	return t.traceContext
}
