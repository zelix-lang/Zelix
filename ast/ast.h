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
#include <ast/rule.h>

typedef struct ast_t ast_t;
DEFINE_VECTOR(ast_t *, ast); // Define a vector for ast_t

/**
 * @brief Abstract Syntax Tree (AST) node structure for the Fluent language.
 *
 * Represents a node in the AST, containing rule information, child nodes,
 * optional value, and source location metadata.
 */
struct ast_t
{
    ast_rule_t rule;      /**< The grammar rule associated with this node. */
    vector_ast_t *children;   /**< Vector of child AST nodes. */
    heap_guard_str_t *value;  /**< Optional value for the node, e.g., identifier name, string literal, etc. */
    size_t line;          /**< Line number in the source file. */
    size_t column;        /**< Column number in the source file. */
    size_t col_start;     /**< Column where the node starts. */
};

/**
 * @brief Allocates and initializes a new AST node.
 *
 * This function creates a new AST node using the provided memory allocator.
 * If `children_required` is true, it also allocates and initializes a vector
 * for the node's children using the specified vector allocator.
 *
 * @param allocator         Arena allocator for the AST node itself.
 * @param vec_allocator     Arena allocator for the children vector (if required).
 * @param children_required Whether to allocate and initialize the children vector.
 * @return Pointer to the newly created ast_t node, or NULL on allocation failure.
 */
static inline ast_t *ast_new(
    arena_allocator_t *const allocator,
    arena_allocator_t *const vec_allocator,
    const bool children_required
)
{
    // Allocate memory for a new AST node
    ast_t *node = (ast_t *)arena_malloc(allocator);
    if (!node)
    {
        return NULL; // Allocation failed
    }

    // Initialize the children vector
    if (children_required)
    {
        node->children = arena_malloc(vec_allocator);
        if (!node->children)
        {
            return NULL;
        }

        // Initialize the vector
        vec_ast_init(node->children, 15, 1.5);
    }

    // Initialize other fields
    node->value = NULL;
    node->line = 0;
    node->column = 0;
    node->col_start = 0;

    return node; // Return the newly created AST node
}

#endif //FLUENT_AST_H
