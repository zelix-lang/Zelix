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
	TOKEN_FUNCTION = 0,      // 0
	TOKEN_LET,               // 1
	TOKEN_CONST,             // 2
	TOKEN_IF,                // 3
	TOKEN_ELSE,              // 4
	TOKEN_ELSE_IF,           // 5
	TOKEN_MOD,               // 6
	TOKEN_RETURN,            // 7
	TOKEN_ASSIGN,            // 8
	TOKEN_PLUS,              // 9
	TOKEN_MINUS,             // 10
	TOKEN_ASTERISK,          // 11
	TOKEN_SLASH,             // 12
	TOKEN_LESS_THAN,         // 13
	TOKEN_GREATER_THAN,      // 14
	TOKEN_EQUAL,             // 15
	TOKEN_NOT_EQUAL,         // 16
	TOKEN_GREATER_THAN_OR_EQUAL, // 17
	TOKEN_LESS_THAN_OR_EQUAL,    // 18
	TOKEN_ARROW,             // 19
	TOKEN_COMMA,             // 20
	TOKEN_SEMICOLON,         // 21
	TOKEN_OPEN_PAREN,        // 22
	TOKEN_CLOSE_PAREN,       // 23
	TOKEN_OPEN_CURLY,        // 24
	TOKEN_CLOSE_CURLY,       // 25
	TOKEN_COLON,             // 26
	TOKEN_NOT,               // 27
	TOKEN_OR,                // 28
	TOKEN_AND,               // 29
	TOKEN_OPEN_BRACKET,      // 30
	TOKEN_CLOSE_BRACKET,     // 31
	TOKEN_DOT,               // 32
	TOKEN_STRING,            // 33
	TOKEN_NUM,               // 34
	TOKEN_DEC,               // 35
	TOKEN_NOTHING,           // 36
	TOKEN_BOOL,              // 37
	TOKEN_STRING_LITERAL,    // 38
	TOKEN_NUM_LITERAL,       // 39
	TOKEN_DECIMAL_LITERAL,   // 40
	TOKEN_BOOL_LITERAL,      // 41
	TOKEN_WHILE,             // 42
	TOKEN_FOR,               // 43
	TOKEN_NEW,               // 44
	TOKEN_IN,                // 45
	TOKEN_TO,                // 46
	TOKEN_BREAK,             // 47
	TOKEN_CONTINUE,          // 48
	TOKEN_PUB,               // 49
	TOKEN_AMPERSAND,         // 50
	TOKEN_BAR,               // 51
	TOKEN_IMPORT,            // 52
	TOKEN_IDENTIFIER,        // 53
	TOKEN_UNKNOWN            // 54
} token_type_t;

typedef struct
{
    char *value;        // The value of the token
    token_type_t type;  // The type of the token (e.g., identifier, keyword, operator)
    int line;           // The line number where the token was found
    int column;         // The column number where the token was found
} token_t;

#endif //FLUENT_TOKEN_H
