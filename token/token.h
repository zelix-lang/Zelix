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
    TOKEN_UNKNOWN = 0,       ///< Unknown or invalid token
    TOKEN_FUNCTION,          ///< Function keyword
    TOKEN_LET,               ///< Let keyword (variable declaration)
    TOKEN_CONST,             ///< Const keyword (constant declaration)
    TOKEN_IF,                ///< If keyword (conditional)
    TOKEN_ELSE,              ///< Else keyword
    TOKEN_ELSE_IF,           ///< Else if keyword
    TOKEN_MOD,               ///< Modulo operator
    TOKEN_RETURN,            ///< Return keyword
    TOKEN_ASSIGN,            ///< Assignment operator (=)
    TOKEN_PLUS,              ///< Addition operator (+)
    TOKEN_MINUS,             ///< Subtraction operator (-)
    TOKEN_ASTERISK,          ///< Multiplication operator (*)
    TOKEN_SLASH,             ///< Division operator (/)
    TOKEN_LESS_THAN,         ///< Less than operator (<)
    TOKEN_GREATER_THAN,      ///< Greater than operator (>)
    TOKEN_EQUAL,             ///< Equality operator (==)
    TOKEN_NOT_EQUAL,         ///< Not equal operator (!=)
    TOKEN_GREATER_THAN_OR_EQUAL, ///< Greater than or equal operator (>=)
    TOKEN_LESS_THAN_OR_EQUAL,    ///< Less than or equal operator (<=)
    TOKEN_ARROW,             ///< Arrow operator (->)
    TOKEN_COMMA,             ///< Comma (,)
    TOKEN_SEMICOLON,         ///< Semicolon (;)
    TOKEN_OPEN_PAREN,        ///< Open parenthesis (()
    TOKEN_CLOSE_PAREN,       ///< Close parenthesis ())
    TOKEN_OPEN_CURLY,        ///< Open curly brace ({)
    TOKEN_CLOSE_CURLY,       ///< Close curly brace (})
    TOKEN_COLON,             ///< Colon (:)
    TOKEN_NOT,               ///< Logical NOT operator (!)
    TOKEN_OR,                ///< Logical OR operator (||)
    TOKEN_AND,               ///< Logical AND operator (&&)
    TOKEN_OPEN_BRACKET,      ///< Open square bracket ([)
    TOKEN_CLOSE_BRACKET,     ///< Close square bracket (])
    TOKEN_DOT,               ///< Dot (.)
    TOKEN_STRING,            ///< String type keyword
    TOKEN_NUM,               ///< Number type keyword
    TOKEN_DEC,               ///< Decimal type keyword
    TOKEN_NOTHING,           ///< Nothing/null type keyword
    TOKEN_BOOL,              ///< Boolean type keyword
    TOKEN_STRING_LITERAL,    ///< String literal
    TOKEN_NUM_LITERAL,       ///< Integer literal
    TOKEN_DECIMAL_LITERAL,   ///< Decimal literal
    TOKEN_BOOL_LITERAL,      ///< Boolean literal
    TOKEN_WHILE,             ///< While keyword (loop)
    TOKEN_FOR,               ///< For keyword (loop)
    TOKEN_NEW,               ///< New keyword (object/instance creation)
    TOKEN_IN,                ///< In keyword (membership/iteration)
    TOKEN_TO,                ///< To keyword (range/iteration)
    TOKEN_BREAK,             ///< Break keyword (loop control)
    TOKEN_CONTINUE,          ///< Continue keyword (loop control)
    TOKEN_PUB,               ///< Public visibility keyword
    TOKEN_AMPERSAND,         ///< Ampersand (&)
    TOKEN_BAR,               ///< Bar/pipe (|)
    TOKEN_IMPORT,            ///< Import keyword (module import)
    TOKEN_IDENTIFIER,        ///< Identifier (variable/function name)
} token_type_t;

/**
 * @brief Represents a token in the Fluent programming language.
 *
 * This structure holds information about a single token, including its
 * string value, type, and the position (line and column) where it was found
 * in the source code.
 */
typedef struct
{
    heap_guard_t *value; // The value of the token
    token_type_t type;   // The type of the token (e.g., identifier, keyword, operator)
    int line;            // The line number where the token was found
    int column;          // The column number where the token was found
} token_t;

#endif //FLUENT_TOKEN_H
