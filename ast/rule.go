/*
   The Fluent Programming Language
   -----------------------------------------------------
   This code is released under the GNU GPL v3 license.
   For more information, please visit:
   https://www.gnu.org/licenses/gpl-3.0.html
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package ast

type Rule int

const (
	Program Rule = iota
	Statement
	Expression
	Identifier
	String
	Number
	Decimal
	Bool
	Nothing
	StringLiteral
	NumberLiteral
	BooleanLiteral
	DecimalLiteral
	FunctionCall
	Public
	Function
	Module
	ObjectCreation
	PropertyAccess
	Assignment
	If
	Else
	ElseIf
	While
	For
	Return
	Break
	Continue
	Import
	ArithmeticExpression
	ArithmeticSign
	Pointer
	Dereference
	ArrayType
	Array
	BooleanOperator
	BooleanExpression
	Parameters
	Parameter
	Type
	Templates
	Generics
	InferredType
	Declaration
	DeclarationType
	IncompleteDeclaration
	Block
)

// String returns the string representation of the rule.
func (r Rule) String() string {
	switch r {
	case Program:
		return "Program"
	case Statement:
		return "Statement"
	case Expression:
		return "Expression"
	case Identifier:
		return "Identifier"
	case String:
		return "String"
	case Nothing:
		return "Nothing"
	case Bool:
		return "Bool"
	case Decimal:
		return "Decimal"
	case StringLiteral:
		return "StringLiteral"
	case NumberLiteral:
		return "NumberLiteral"
	case BooleanLiteral:
		return "BooleanLiteral"
	case DecimalLiteral:
		return "DecimalLiteral"
	case FunctionCall:
		return "FunctionCall"
	case Public:
		return "Public"
	case Function:
		return "Function"
	case Module:
		return "Module"
	case ObjectCreation:
		return "ObjectCreation"
	case PropertyAccess:
		return "PropertyAccess"
	case Assignment:
		return "Assignment"
	case If:
		return "If"
	case Else:
		return "Else"
	case ElseIf:
		return "ElseIf"
	case While:
		return "While"
	case For:
		return "For"
	case Return:
		return "Return"
	case Break:
		return "Break"
	case Continue:
		return "Continue"
	case Import:
		return "Import"
	case ArithmeticExpression:
		return "ArithmeticExpression"
	case ArithmeticSign:
		return "ArithmeticSign"
	case Pointer:
		return "Pointer"
	case Dereference:
		return "Dereference"
	case ArrayType:
		return "ArrayType"
	case Array:
		return "Array"
	case BooleanOperator:
		return "BooleanOperator"
	case BooleanExpression:
		return "BooleanExpression"
	case Parameters:
		return "Parameters"
	case Parameter:
		return "Parameter"
	case Type:
		return "Type"
	case Templates:
		return "Templates"
	case Generics:
		return "Generics"
	case InferredType:
		return "InferredType"
	case Declaration:
		return "Declaration"
	case DeclarationType:
		return "DeclarationType"
	case IncompleteDeclaration:
		return "IncompleteDeclaration"
	case Block:
		return "Block"
	}

	return ""
}
