package code

import (
	"time"
	"zyro/code/wrapper"
	"zyro/token"
)

// Function represents a function in the abstract syntax tree (AST).
type Function struct {
	// returnType holds the tokens representing the return type of the function.
	returnType wrapper.TypeWrapper
	// parameters holds the tokens representing the parameters of the function.
	// Each parameter is represented as a slice of tokens.
	parameters map[string]wrapper.TypeWrapper
	// body holds the tokens representing the body of the function.
	body []token.Token
	// public holds whether the function is public or not.
	public bool
	// std holds whether the function is a standard library function or not.
	std bool
	// trace holds the token that triggered the creation of the function.
	trace token.Token
	// timesCalled holds the number of times the function has been called.
	timesCalled int
	// lastCalled holds the time the function was last called.
	lastCalled time.Time
}

// NewFunction creates a new Function
func NewFunction(
	returnType []token.Token,
	parameters map[string][]token.Token,
	body []token.Token,
	public bool,
	std bool,
	trace token.Token,
) Function {
	wrappers := make(map[string]wrapper.TypeWrapper)

	for key, value := range parameters {
		wrappers[key] = wrapper.NewTypeWrapper(value, trace, false)
	}

	return Function{
		returnType: wrapper.NewTypeWrapper(returnType, trace, false),
		parameters: wrappers,
		body:       body,
		public:     public,
		std:        std,
		trace:      trace,
	}
}

// GetReturnType returns the return type of the function.
func (f *Function) GetReturnType() wrapper.TypeWrapper {
	return f.returnType
}

// GetParameters returns the parameters of the function.
func (f *Function) GetParameters() map[string]wrapper.TypeWrapper {
	return f.parameters
}

// GetBody returns the body of the function.
func (f *Function) GetBody() []token.Token {
	return f.body
}

// IsPublic returns whether the function is public or not.
func (f *Function) IsPublic() bool {
	return f.public
}

// IsStd returns whether the function is a standard library function or not.
func (f *Function) IsStd() bool {
	return f.std
}

// GetTrace returns the token that triggered the creation of the function.
func (f *Function) GetTrace() token.Token {
	return f.trace
}

// GetTimesCalled returns the number of times the function has been called.
func (f *Function) GetTimesCalled() int {
	return f.timesCalled
}

// GetLastCalled returns the time the function was last called.
func (f *Function) GetLastCalled() time.Time {
	return f.lastCalled
}

// SetTimesCalled sets the number of times the function has been called.
func (f *Function) SetTimesCalled(timesCalled int) {
	f.timesCalled = timesCalled
}

// SetLastCalled sets the time the function was last called.
func (f *Function) SetLastCalled(lastCalled time.Time) {
	f.lastCalled = lastCalled
}
