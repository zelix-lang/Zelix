package ast

import "surf/code"

// Function represents a function in the abstract syntax tree (AST).
type Function struct {
	// returnType holds the tokens representing the return type of the function.
	returnType []code.Token
	// parameters holds the tokens representing the parameters of the function.
	// Each parameter is represented as a slice of tokens.
	parameters map[string][]code.Token
	// body holds the tokens representing the body of the function.
	body []code.Token
	// public holds whether the function is public or not.
	public bool
	// std holds whether the function is a standard library function or not.
	std bool
	// trace holds the token that triggered the creation of the function.
	trace code.Token
}

// GetReturnType returns the return type of the function.
func (f Function) GetReturnType() []code.Token {
	return f.returnType
}

// GetParameters returns the parameters of the function.
func (f Function) GetParameters() map[string][]code.Token {
	return f.parameters
}

// GetBody returns the body of the function.
func (f Function) GetBody() []code.Token {
	return f.body
}

// IsPublic returns whether the function is public or not.
func (f Function) IsPublic() bool {
	return f.public
}

// IsStd returns whether the function is a standard library function or not.
func (f Function) IsStd() bool {
	return f.std
}

// GetTrace returns the token that triggered the creation of the function.
func (f Function) GetTrace() code.Token {
	return f.trace
}
