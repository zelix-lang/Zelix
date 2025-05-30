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
// Created by rodrigo on 5/30/25.
//

#ifndef FLUENT_LEXER_ERROR_H
#define FLUENT_LEXER_ERROR_H

// ============= FLUENT LIB C =============
#include <fluent/types/types.h>

/**
 * @enum lexer_error_code_t
 * @brief Enumerates possible lexer error codes.
 *
 * This enum defines the types of errors that can be encountered during
 * lexical analysis in the Fluent programming language.
 */
typedef enum
{
    LEXER_ERROR_UNKNOWN_TOKEN = 0,        ///< An unknown token was encountered
    LEXER_ERROR_UNKNOWN_ESCAPE,           ///< An unknown escape sequence was encountered
    LEXER_ERROR_UNTERMINATED_STRING,      ///< A string literal was not properly terminated
    LEXER_ERROR_UNTERMINATED_COMMENT,     ///< A comment was not properly terminated
} lexer_error_code_t;

/**
 * @struct lexer_error_t
 * @brief Represents a lexer error with location and type.
 *
 * This struct contains information about a lexer error, including
 * the line and column where the error occurred and the error code.
 */
typedef struct
{
    size_t line;             ///< The line number where the error occurred
    size_t column;           ///< The column number where the error occurred
    lexer_error_code_t code; ///< The error code representing the type of error
} lexer_error_t;

#endif //FLUENT_LEXER_ERROR_H
