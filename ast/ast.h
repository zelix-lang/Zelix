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

#ifndef FLUENT_AST_H
#define FLUENT_AST_H

// ============= FLUENT LIB C =============
#include <fluent/vector/vector.h> // fluent_libc

// ============= INCLUDES =============
#include "rule.h"

/**
 * @brief Abstract Syntax Tree (AST) node structure for the Fluent language.
 *
 * Represents a node in the AST, containing rule information, child nodes,
 * optional value, and source location metadata.
 */
typedef struct ast_t
{
    ast_rule_t rule;      /**< The grammar rule associated with this node. */
    vector_generic_t *children;   /**< Vector of child AST nodes. */
    char *value;          /**< Optional value for the node, e.g., identifier name, string literal, etc. */
    size_t line;          /**< Line number in the source file. */
    size_t column;        /**< Column number in the source file. */
    size_t col_start;     /**< Column where the node starts. */
    char *file;           /**< File where the node was defined. */
} ast_t;

#endif //FLUENT_AST_H
