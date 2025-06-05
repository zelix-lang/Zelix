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

//
// Created by rodrigo on 6/5/25.
//

#ifndef FLUENT_AST_RULE_H
#define FLUENT_AST_RULE_H

typedef enum
{
    AST_PROGRAM_RULE,
    AST_STATEMENT,
    AST_EXPRESSION,
    AST_IDENTIFIER,
    AST_STRING,
    AST_NUMBER,
    AST_DECIMAL,
    AST_BOOL,
    AST_NOTHING,
    AST_STRING_LITERAL,
    AST_NUMBER_LITERAL,
    AST_BOOLEAN_LITERAL,
    AST_DECIMAL_LITERAL,
    AST_FUNCTION_CALL,
    AST_PUBLIC,
    AST_FUNCTION,
    AST_MODULE,
    AST_OBJECT_CREATION,
    AST_PROPERTY_ACCESS,
    AST_ASSIGNMENT,
    AST_IF,
    AST_ELSE,
    AST_ELSE_IF,
    AST_WHILE,
    AST_FOR,
    AST_RETURN,
    AST_BREAK,
    AST_CONTINUE,
    AST_IMPORT,
    AST_ARITHMETIC_EXPRESSION,
    AST_ARITHMETIC_SIGN,
    AST_POINTER,
    AST_DEREFERENCE,
    AST_ARRAY_TYPE,
    AST_ARRAY,
    AST_BOOLEAN_OPERATOR,
    AST_BOOLEAN_EXPRESSION,
    AST_PARAMETERS,
    AST_PARAMETER,
    AST_TYPE,
    AST_TEMPLATES,
    AST_GENERICS,
    AST_INFERRED_TYPE,
    AST_DECLARATION,
    AST_CONST,
    AST_LET,
    AST_INCOMPLETE_DECLARATION,
    AST_BLOCK
} ast_rule_t;

#endif //FLUENT_AST_RULE_H
