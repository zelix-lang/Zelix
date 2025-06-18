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

#ifndef FLUENT_ERROR_H
#define FLUENT_ERROR_H

// ============= INCLUDES =============
#include <ast/ast.h>

/**
 * @brief Represents a syntax error encountered during parsing.
 *
 * This structure holds information about the location of the error
 * (line and column) and an array of up to 5 expected AST rules
 * that could have been valid at the error location.
 */
typedef struct
{
    size_t line;               /**< The line number where the error occurred. */
    size_t column;             /**< The column number where the error occurred. */
    size_t col_start;          /**< The starting column of the error. */
    ast_rule_t expected[5];    /**< An array of expected rules at the error location. */
} ast_error_t;

// ============= GLOBAL VARIABLE =============
static ast_error_t global_parser_error;

/**
 * @brief Creates and initializes a syntax error record.
 *
 * This function sets the global parser error with the provided line, column,
 * and starting column information, as well as the expected AST rules at the
 * error location. It copies up to `expected_len` rules from the `expected`
 * array into the global error structure.
 *
 * @param line         The line number where the error occurred.
 * @param column       The column number where the error occurred.
 * @param col_start    The starting column of the error.
 * @param expected     Pointer to an array of expected AST rules.
 * @param expected_len The number of expected rules to copy (up to 5).
 * @return Pointer to the global parser error structure.
 */
static ast_error_t *create_error_ranged(
    const size_t line,
    const size_t column,
    const size_t col_start,
    const ast_rule_t *const expected,
    const size_t expected_len
)
{
    // Set the global parser error
    global_parser_error.line = line;
    global_parser_error.column = column;
    global_parser_error.col_start = col_start;

    // Copy the expected rules into the global parser error
    memcpy(global_parser_error.expected, expected, sizeof(ast_rule_t) * expected_len);

    return &global_parser_error;
}

static ast_error_t *create_error(
    const token_stream_t *stream,
    const ast_rule_t *const expected,
    const size_t expected_len
)
{
    // Get the current token from the stream
    const token_t *current_token = token_stream_nth(stream, stream->current);

    // Set the global parser error
    global_parser_error.line = current_token->line;
    global_parser_error.column = current_token->column;
    global_parser_error.col_start = current_token->col_start;

    // Copy the expected rules into the global parser error
    memcpy(global_parser_error.expected, expected, sizeof(ast_rule_t) * expected_len);

    return &global_parser_error;
}

#endif //FLUENT_ERROR_H
