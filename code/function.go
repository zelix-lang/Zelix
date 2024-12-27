package code

import (
	"fluent/code/wrapper"
	"fluent/token"
	"fluent/tokenUtil/generic"
	"fluent/tokenUtil/splitter"
	"log"
	"time"
)

// FunctionParam represents a parameter in the abstract syntax tree (AST).
type FunctionParam struct {
	// name holds the name of the parameter.
	name string
	// typ holds the type of the parameter.
	typ wrapper.TypeWrapper
	// tokens holds the tokens representing the parameter.
	tokens []token.Token
}

// Function represents a function in the abstract syntax tree (AST).
type Function struct {
	// returnType holds the tokens representing the return type of the function.
	returnType wrapper.TypeWrapper
	// parameters holds the tokens representing the parameters of the function.
	// Each parameter is represented as a slice of tokens.
	parameters []FunctionParam
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

// NewFunctionParam creates a new FunctionParam.
func NewFunctionParam(name string, typ wrapper.TypeWrapper, tokens []token.Token) FunctionParam {
	return FunctionParam{
		name:   name,
		typ:    typ,
		tokens: tokens,
	}
}

// NewFunction creates a new Function
func NewFunction(
	returnType []token.Token,
	parameters []FunctionParam,
	body []token.Token,
	public bool,
	std bool,
	trace token.Token,
) Function {
	wrappers := make([]FunctionParam, len(parameters))

	for i, val := range parameters {
		wrappers[i] = FunctionParam{
			name:   val.name,
			typ:    wrapper.NewTypeWrapper(val.tokens, val.tokens[0]),
			tokens: val.tokens,
		}
	}

	return Function{
		returnType: wrapper.NewTypeWrapper(returnType, trace),
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
func (f *Function) GetParameters() []FunctionParam {
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

// BuildWithoutGenerics builds a new function, replacing
// generics in constants and vars declarations with the given types
func (f *Function) BuildWithoutGenerics(types map[string]wrapper.TypeWrapper) Function {
	params := f.GetParameters()
	newParams := make([]FunctionParam, len(params))

	for key, value := range params {
		newValue := generic.ConvertGeneric(value.typ, types)
		log.Printf("Converting parameter: %v from %v to %v", key, value, newValue)

		newParams[key] = FunctionParam{
			name:   value.name,
			typ:    newValue,
			tokens: value.tokens,
		}
	}

	body := f.body
	newBody := make([]token.Token, 0)

	for i, t := range f.GetBody() {
		tokenType := t.GetType()

		if tokenType == token.Let || tokenType == token.Const {
			// Declarations go like:
			// let my_var : str = "hello";
			// we have to skip 3 tokens to get to the type
			declaration, _ := splitter.ExtractTokensBefore(
				body[i:],
				token.Semicolon,
				false,
				token.Unknown,
				token.Unknown,
				true,
			)

			newBody = append(newBody, generic.ConvertVariableGenerics(declaration, types)...)
			continue
		}

		newBody = append(newBody, t)
	}

	return Function{
		returnType: generic.ConvertGeneric(f.returnType, types),
		parameters: newParams,
		body:       newBody,
		public:     f.IsPublic(),
		std:        f.IsStd(),
		trace:      f.GetTrace(),
	}
}

// GetName returns the name of the parameter.
func (p FunctionParam) GetName() string {
	return p.name
}

// GetType returns the type of the parameter.
func (p FunctionParam) GetType() wrapper.TypeWrapper {
	return p.typ
}

// GetTokens returns the tokens of the parameter.
func (p FunctionParam) GetTokens() []token.Token {
	return p.tokens
}
