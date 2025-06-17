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
// Created by rodrigo on 6/16/25.
//

#ifndef FLUENT_PARSER_RULE_EXPRESSION_H
#define FLUENT_PARSER_RULE_EXPRESSION_H

// ============= FLUENT LIB C =============
#include <fluent/arena/arena.h> // fluent_libc
#include <fluent/std_bool/std_bool.h> // fluent_libc

// ============= INCLUDES =============
#include "../../ast/ast.h"
#include "../../token/token.h"
#include "../queue/expression.h"

static inline bool parse_expression(
    ast_t *const root,
    token_t **body,
    const size_t start,
    const size_t body_len,
    arena_allocator_t *const arena,
    arena_allocator_t *const vec_arena
)
{
    // Create a new expression AST for the result
    ast_t *expression = ast_new(arena, vec_arena, TRUE);

    // Handle failure
    if (!expression)
    {
        // Failed to allocate memory for the expression
        return FALSE;
    }

    // Set the metadata for the expression
    expression->rule = AST_EXPRESSION;

    // Create a new queue for the expression body
    alinked_queue_expr_t queue;
    alinked_queue_expr_init(&queue, 10);
    alinked_queue_expr_append(&queue, (queue_expression_t){
        .body = body,
        .start = start,
        .len = body_len - start
    });

    // Iterate over the queue until it's empty
    while (queue.len > 0)
    {
        // Get the first element
        const queue_expression_t expr = alinked_queue_expr_shift(&queue);
        token_t **input = expr.body + expr.start;
        const size_t len = expr.len - expr.start;
        size_t start = 0;

        // Handle empty input
        if (len == 0)
        {
            return FALSE;
        }

        // Parse pointers
        for (size_t i = 0; i < len; i++)
        {
            // Get the current token
            const token_t *token = input[i];

            // Check if we have a dereference operation
            if (token->type == TOKEN_ASTERISK)
            {

            }
            // Check for pointers
            else if (
                token->type == TOKEN_AMPERSAND ||
                token->type == TOKEN_AND
            )
            {

            }
        }
    }

    // Add the child to the root
    vec_ast_push(root->children, expression);
    return TRUE;
}

#endif //FLUENT_PARSER_RULE_EXPRESSION_H
