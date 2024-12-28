package token

// Type represents the type of any token in the source code.
type Type int

const (
	// Function represents a function token.
	Function Type = iota

	// Let represents a let token.
	Let
	// Const represents a const token.
	Const
	// If represents an "if" token.
	If
	// Else represents an else token.
	Else
	// ElseIf represents an else if token.
	ElseIf
	// Mod represents a mod token.
	Mod
	// Return represents a return token.
	Return
	// Assign represents an assignment token.
	Assign
	// AssignAdd represents an addition assignment token.
	AssignAdd
	// AssignSub represents a subtraction assignment token.
	AssignSub
	// AssignSlash represents a division assignment token.
	AssignSlash
	// AssignAsterisk represents a multiplication assignment token.
	AssignAsterisk

	// Plus represents a plus token.
	Plus
	// Minus represents a minus token.
	Minus
	// Increment represents an increment token.
	Increment
	// Decrement represents a decrement token.
	Decrement

	// Asterisk represents an asterisk token.
	Asterisk
	// Slash represents a slash token.
	Slash
	// LessThan represents a less than token.
	LessThan
	// GreaterThan represents a greater than token.
	GreaterThan
	// Equal represents an equal token.
	Equal
	// NotEqual represents a not equal token.
	NotEqual
	// GreaterThanOrEqual represents a greater than or equal token.
	GreaterThanOrEqual
	// LessThanOrEqual represents a less than or equal token.
	LessThanOrEqual
	// Percent represents a percent token.
	Percent

	// Arrow represents an arrow token.
	Arrow
	// Comma represents a comma token.
	Comma

	// Semicolon represents a semicolon token.
	Semicolon

	// OpenParen represents an open parenthesis token.
	OpenParen
	// CloseParen represents a close parenthesis token.
	CloseParen

	// OpenCurly represents an open curly brace token.
	OpenCurly
	// CloseCurly represents a close curly brace token.
	CloseCurly

	// Colon represents a colon token.
	Colon
	// Xor represents an XOR token.
	Xor
	// Not represents a not token.
	Not

	// OpenBracket represents an open bracket token.
	OpenBracket
	// CloseBracket represents a close bracket token.
	CloseBracket

	// Dot represents a dot token.
	Dot

	// String represents a string token.
	String
	// Num represents a number token.
	Num
	// Dec represents a decimal token.
	Dec
	// Nothing represents a nothing token.
	Nothing
	// Bool represents a boolean token.
	Bool

	// StringLiteral represents a string literal token.
	StringLiteral
	// NumLiteral represents a number literal token.
	NumLiteral
	// DecimalLiteral represents a decimal literal token.
	DecimalLiteral
	// BoolLiteral represents a boolean literal token.
	BoolLiteral

	// While represents a while token.
	While
	// For represents a for token.
	For

	// New represents a new token.
	New

	// In represents an in token.
	In

	// Break represents a break token.
	Break
	// Continue represents a continue token.
	Continue
	// Pub represents a public token.
	Pub

	// Ampersand represents an ampersand token.
	Ampersand
	// Bar represents a bar token.
	Bar

	// Import represents an import token.
	Import

	// Identifier represents an identifier token.
	Identifier
	// Unknown represents an unknown token.
	Unknown
)
