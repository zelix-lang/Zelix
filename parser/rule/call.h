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
// Created by rodrigo on 7/4/25.
//

#ifndef FLUENT_PARSER_RULE_CALL_H
#define FLUENT_PARSER_RULE_CALL_H

// ============= FLUENT LIB C =============
#include <fluent/arena/arena.h> // fluent_libc
#include <fluent/std_bool/std_bool.h> // fluent_libc
#include <fluent/std_math/std_math.h> // fluent_libc

// ============= INCLUDES =============
#include "ast/ast.h"
#include "token/token.h"
#include "parser/error.h"
#include "parser/extractor.h"
#include "parser/queue/expression.h"
#include "parser/rule/args.h"
#include "parser/shared/parse_pair.h"

static inline pair_parse_t parse_call(
    token_t **tokens,
    const size_t len,
    arena_allocator_t *const arena,
    arena_allocator_t *const vec_arena,
    alinked_queue_expr_t *queue,
    ast_t *const candidate
)
{
    // Create a new AST node for the function call
    ast_t *call_node = ast_new(arena, vec_arena, TRUE);
    if (!call_node)
    {
        // Failed to allocate memory for the function call node
        return pair_parse_new(NULL, 0); // Return a pair with NULL and 0
    }

    // Set the rule type for the function call
    call_node->rule = AST_FUNCTION_CALL;
    // Set the line and column information from the first token
    call_node->line = tokens[0]->line;
    call_node->column = tokens[0]->column;
    call_node->col_start = tokens[0]->col_start;

    // Append the candidate as the first child of the call node
    vec_ast_push(call_node->children, candidate);

    // Check if we have an empty function call
    if (tokens[1]->type == TOKEN_CLOSE_PAREN)
    {
        return pair_parse_new(call_node, 2); // Return the call node with 2 tokens consumed
    }

    // Parse the arguments
    const pair_parse_t args = parse_args(
        tokens,
        len,
        arena,
        vec_arena,
        queue
    );

    // Handle failure in parsing arguments
    if (!args.first)
    {
        // Failed to parse the arguments
        create_error_ranged(
            tokens[0]->line,
            tokens[0]->column,
            tokens[0]->col_start,
            (ast_rule_t[]){AST_FUNCTION_CALL},
            1
        );
        return pair_parse_new(NULL, 0); // Return a pair with NULL and 0
    }

    // +2 for the open and close parentheses
    return pair_parse_new(call_node, args.second + 2);
}

#endif //FLUENT_PARSER_RULE_CALL_H
