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
// Created by rodrigo on 5/29/25.
//

#ifndef FLUENT_TOKEN_H
#define FLUENT_TOKEN_H

typedef enum
{
	TOKEN_FUNCTION,
	TOKEN_LET,
	TOKEN_CONST,
	TOKEN_IF,
	TOKEN_ELSE,
	TOKEN_ELSE_IF,
	TOKEN_MOD,
	TOKEN_RETURN,
	TOKEN_ASSIGN,
	TOKEN_PLUS,
	TOKEN_MINUS,
	TOKEN_ASTERISK,
	TOKEN_SLASH,
	TOKEN_LESS_THAN,
	TOKEN_GREATER_THAN,
	TOKEN_EQUAL,
	TOKEN_NOT_EQUAL,
	TOKEN_GREATER_THAN_OR_EQUAL,
	TOKEN_LESS_THAN_OR_EQUAL,
	TOKEN_ARROW,
	TOKEN_COMMA,
	TOKEN_SEMICOLON,
	TOKEN_OPEN_PAREN,
	TOKEN_CLOSE_PAREN,
	TOKEN_OPEN_CURLY,
	TOKEN_CLOSE_CURLY,
	TOKEN_COLON,
	TOKEN_NOT,
	TOKEN_OR,
	TOKEN_AND,
	TOKEN_OPEN_BRACKET,
	TOKEN_CLOSE_BRACKET,
	TOKEN_DOT,
	TOKEN_STRING,
	TOKEN_NUM,
	TOKEN_DEC,
	TOKEN_NOTHING,
	TOKEN_BOOL,
	TOKEN_STRING_LITERAL,
	TOKEN_NUM_LITERAL,
	TOKEN_DECIMAL_LITERAL,
	TOKEN_BOOL_LITERAL,
	TOKEN_WHILE,
	TOKEN_FOR,
	TOKEN_NEW,
	TOKEN_IN,
	TOKEN_TO,
	TOKEN_BREAK,
	TOKEN_CONTINUE,
	TOKEN_PUB,
	TOKEN_AMPERSAND,
	TOKEN_BAR,
	TOKEN_IMPORT,
	TOKEN_IDENTIFIER,
	TOKEN_UNKNOWN
} token_type_t;

typedef struct
{
    char *value;        // The value of the token
    token_type_t type;  // The type of the token (e.g., identifier, keyword, operator)
    int line;           // The line number where the token was found
    int column;         // The column number where the token was found
} token_t;

#endif //FLUENT_TOKEN_H
