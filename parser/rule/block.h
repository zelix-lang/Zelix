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

#ifndef FLUENT_PARSER_RULE_BLOCK_H
#define FLUENT_PARSER_RULE_BLOCK_H

// ============= FLUENT LIB C =============
#include <fluent/arena/arena.h> // fluent_libc
#include <fluent/std_bool/std_bool.h> // fluent_libc

// ============= INCLUDES =============
#include "../../ast/ast.h"

static inline bool parse_block(
    ast_t *const root,
    token_t **body,
    const size_t body_len,
    arena_allocator_t *const arena,
    arena_allocator_t *const vec_arena
) {
    // Track the nesting level
    size_t nest_level = 0;

    // Create a new queue for nested blocks
    alinked_queue_ast_t queue;
    alinked_queue_ast_init(&queue, 10);
    alinked_queue_ast_append(&queue, root);

    // Iterate over the block
    for (size_t i = 0; i < body_len; ++i)
    {
        // Get the current block
        ast_t *block = queue.head->data;

        // Get the current token
        const token_t *token = body[i];

        // Determine what to do based on the token type
        switch (token->type)
        {
            case TOKEN_OPEN_CURLY:
            {
                // Increment the nesting level
                nest_level++;
                break;
            }

            case TOKEN_CLOSE_CURLY:
            {
                // Check if we have a matching opening curly brace
                if (nest_level == 0)
                {
                    // Create an error for unmatched closing brace
                    create_error(
                        token->line,
                        token->column,
                        token->col_start,
                        (ast_rule_t[]){AST_BLOCK},
                        1
                    );
                    return FALSE; // Invalid block
                }

                // Decrement the nesting level
                nest_level--;
                break;
            }

            case TOKEN_LET:
            case TOKEN_CONST:
            {
                // TODO: Variable declaration parser
                break;
            }

            default:
            {
                // TODO: Expr parser
                break;
            }
        }
    }

    // Make sure we don't end up with a nested level
    if (nest_level != 0)
    {
        // Create an error for unmatched block
        create_error(
            body[0]->line,
            body[0]->column,
            body[0]->col_start,
            (ast_rule_t[]){AST_BLOCK},
            1
        );
        return FALSE; // Invalid block
    }

    return TRUE;
}

#endif //FLUENT_PARSER_RULE_BLOCK_H
