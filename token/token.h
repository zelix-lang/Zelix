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

/**
 * @enum token_type_t
 * @brief Enumerates all possible token types in the Fluent programming language.
 *
 * Each value represents a distinct kind of token that can be identified by the lexer,
 * including keywords, operators, literals, punctuation, and special symbols.
 * The integer values are assigned sequentially, starting from 0.
 */
typedef enum
{
    TOKEN_UNKNOWN = 0,       ///< Unknown or invalid token - 0
    TOKEN_FUNCTION,          ///< Function keyword - 1
    TOKEN_LET,               ///< Let keyword (variable declaration) - 2
    TOKEN_CONST,             ///< Const keyword (constant declaration) - 3
    TOKEN_IF,                ///< If keyword (conditional) - 4
    TOKEN_ELSE,              ///< Else keyword - 5
    TOKEN_ELSE_IF,           ///< Else if keyword - 6
    TOKEN_MOD,               ///< Module keyword - 7
    TOKEN_RETURN,            ///< Return keyword - 8
    TOKEN_ASSIGN,            ///< Assignment operator (=) - 9
    TOKEN_PLUS,              ///< Addition operator (+) - 10
    TOKEN_MINUS,             ///< Subtraction operator (-) - 11
    TOKEN_ASTERISK,          ///< Multiplication operator (*) - 12
    TOKEN_SLASH,             ///< Division operator (/) - 13
    TOKEN_LESS_THAN,         ///< Less than operator (<) - 14
    TOKEN_GREATER_THAN,      ///< Greater than operator (>) - 15
    TOKEN_EQUAL,             ///< Equality operator (==) - 16
    TOKEN_NOT_EQUAL,         ///< Not equal operator (!=) - 17
    TOKEN_GREATER_THAN_OR_EQUAL, ///< Greater than or equal operator (>=) - 18
    TOKEN_LESS_THAN_OR_EQUAL,    ///< Less than or equal operator (<=) - 19
    TOKEN_ARROW,             ///< Arrow operator (->) - 20
    TOKEN_COMMA,             ///< Comma (,) - 21
    TOKEN_SEMICOLON,         ///< Semicolon (;) - 22
    TOKEN_OPEN_PAREN,        ///< Open parenthesis (() - 23
    TOKEN_CLOSE_PAREN,       ///< Close parenthesis ()) - 24
    TOKEN_OPEN_CURLY,        ///< Open curly brace ({) - 25
    TOKEN_CLOSE_CURLY,       ///< Close curly brace (}) - 26
    TOKEN_COLON,             ///< Colon (:) - 27
    TOKEN_NOT,               ///< Logical NOT operator (!) - 28
    TOKEN_OR,                ///< Logical OR operator (||) - 29
    TOKEN_AND,               ///< Logical AND operator (&&) - 30
    TOKEN_OPEN_BRACKET,      ///< Open square bracket ([) - 31
    TOKEN_CLOSE_BRACKET,     ///< Close square bracket (]) - 32
    TOKEN_DOT,               ///< Dot (.) - 33
    TOKEN_STRING,            ///< String type keyword - 34
    TOKEN_NUM,               ///< Number type keyword - 35
    TOKEN_DEC,               ///< Decimal type keyword - 36
    TOKEN_NOTHING,           ///< Nothing/null type keyword - 37
    TOKEN_BOOL,              ///< Boolean type keyword - 38
    TOKEN_STRING_LITERAL,    ///< String literal - 39
    TOKEN_NUM_LITERAL,       ///< Integer literal - 40
    TOKEN_DECIMAL_LITERAL,   ///< Decimal literal - 41
    TOKEN_BOOL_LITERAL,      ///< Boolean literal - 42
    TOKEN_WHILE,             ///< While keyword (loop) - 43
    TOKEN_FOR,               ///< For keyword (loop) - 44
    TOKEN_NEW,               ///< New keyword (object/instance creation) - 45
    TOKEN_IN,                ///< In keyword (membership/iteration) - 46
    TOKEN_TO,                ///< To keyword (range/iteration) - 47
    TOKEN_BREAK,             ///< Break keyword (loop control) - 48
    TOKEN_CONTINUE,          ///< Continue keyword (loop control) - 49
    TOKEN_PUB,               ///< Public visibility keyword - 50
    TOKEN_AMPERSAND,         ///< Ampersand (&) - 51
    TOKEN_BAR,               ///< Bar/pipe (|) - 52
    TOKEN_IMPORT,            ///< Import keyword (module import) - 53
    TOKEN_IDENTIFIER,        ///< Identifier (variable/function name) - 54
} token_type_t;

DEFINE_HEAP_GUARD(char, str, 512);

/**
 * @brief Represents a token in the Fluent programming language.
 *
 * This structure holds information about a single token, including its
 * string value, type, and the position (line and column) where it was found
 * in the source code.
 */
typedef struct
{
    heap_guard_str_t *value; // The value of the token
    token_type_t type;   // The type of the token (e.g., identifier, keyword, operator)
    size_t line;         // The line number where the token was found
    size_t column;       // The column number where the token was found
    size_t col_start;    // The column number where the token starts
} token_t;

#endif //FLUENT_TOKEN_H
