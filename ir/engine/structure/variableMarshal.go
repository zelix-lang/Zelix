package structure

import (
	"fluent/ast"
	"fluent/code/mod"
	"fluent/code/types"
	wrapper2 "fluent/code/wrapper"
	"fluent/ir/wrapper"
	"fluent/token"
	"fluent/tokenUtil/splitter"
	"strconv"
	"strings"
)

// MarshalGenericsInVar marshals the generics in the given variable into the given strings.Builder
// and returns the TypeWrapper that represents the variable's type
func MarshalGenericsInVar(
	variable []token.Token,
	ir *wrapper.IrWrapper,
	builder *strings.Builder,
	counter *int,
	fileCode *ast.FileCode,
) wrapper2.TypeWrapper {
	// Extract the type
	propType, _ := splitter.ExtractTokensBefore(
		variable[3:],
		token.Assign,
		false,
		token.Unknown,
		token.Unknown,
		true,
	)

	// See if the first token is a mod
	firstToken := propType[0]

	// We don't have to see if it's a defined mod here
	// The static analyzer has already done that

	// Wrap into a TypeWrapper
	typeWrapper := wrapper2.NewTypeWrapper(
		propType,
		firstToken,
	)

	if firstToken.GetType() != token.Identifier {
		return typeWrapper
	}

	// See if it has templates
	if len(typeWrapper.GetParameters()) < 1 {
		return typeWrapper
	}

	// See if it has already been built
	if _, built := ir.GetGenericMod(typeWrapper.Marshal()); built {
		return typeWrapper
	}

	// Iterate over the parameters
	for _, param := range typeWrapper.GetParameters() {
		// See if it's a mod
		if param.GetType() != types.ModType {
			continue
		}

		// See if it has templates
		if len(param.GetParameters()) < 1 {
			continue
		}

		// Get the mod
		nestedMod, _, _ := mod.FindMod(
			fileCode.GetModules(),
			param.GetBaseType(),
			firstToken.GetFile(),
		)

		// Build without generics
		builtMod, _ := nestedMod.BuildWithoutGenerics(typeWrapper)

		// Mark this mod as built
		ir.AddGenericMod(typeWrapper.Marshal(), &builtMod)

		// Recursively marshal the mod
		MarshalSingleMod(&builtMod, ir, builder, counter, fileCode)
	}

	// Marshal the mod itself
	module, _, _ := mod.FindMod(
		fileCode.GetModules(),
		typeWrapper.GetBaseType(),
		firstToken.GetFile(),
	)

	builtMod, _ := module.BuildWithoutGenerics(typeWrapper)
	// Mark this mod as built
	ir.AddGenericMod(typeWrapper.Marshal(), &builtMod)

	MarshalSingleMod(&builtMod, ir, builder, counter, fileCode)

	return typeWrapper
}

// MarshalVariable marshals the given variable into the given strings.Builder
func MarshalVariable(
	variable []token.Token,
	ir *wrapper.IrWrapper,
	builder *strings.Builder,
	counter *int,
	fileCode *ast.FileCode,
) {
	*counter++

	propType, _ := splitter.ExtractTokensBefore(
		variable[3:],
		token.Assign,
		false,
		token.Unknown,
		token.Unknown,
		true,
	)

	// See if the first token is a mod
	firstToken := propType[0]

	// We don't have to see if it's a defined mod here
	// The static analyzer has already done that

	// Wrap into a TypeWrapper
	typeWrapper := wrapper2.NewTypeWrapper(
		propType,
		firstToken,
	)

	// See if the type is a mod
	isMod := typeWrapper.GetType() == types.ModType

	if isMod {
		// Write the mod_mov instruction to construct a mod inside a variable
		builder.WriteString("mod_mov x")
	} else {
		// Write the mov instruction
		builder.WriteString("mov x")
	}

	builder.WriteString(strconv.Itoa(*counter))
	builder.WriteString(" ")

	if isMod {
		// See if it's a generic mod
		var retrievedMod *mod.FluentMod

		if genericMod, built := ir.GetGenericMod(typeWrapper.Marshal()); built {
			retrievedMod = genericMod
		} else {
			// Retrieve the mod from the file code
			retrievedMod, _, _ = mod.FindMod(
				fileCode.GetModules(),
				typeWrapper.GetBaseType(),
				variable[0].GetFile(),
			)
		}

		// Get the correspondent counter for this mod
		modCounter, _ := ir.GetMod(retrievedMod)
		builder.WriteString("x")
		builder.WriteString(strconv.Itoa(modCounter))
	} else {

	}

	builder.WriteString("\n")
}
