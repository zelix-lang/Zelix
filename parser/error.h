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
#include "../ast/rule.h"

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

#endif //FLUENT_ERROR_H
