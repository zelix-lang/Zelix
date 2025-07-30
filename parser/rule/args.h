//
// Created by rodrigo on 6/21/25.
//

#ifndef FLUENT_PARSER_RULE_ARGS_H
#define FLUENT_PARSER_RULE_ARGS_H

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
#include "parser/shared/parse_pair.h"

static inline pair_parse_t parse_args(
    token_t **tokens,
    const size_t len,
    arena_allocator_t *const arena,
    arena_allocator_t *const vec_arena,
    alinked_queue_expr_t *queue
)
{
    // Create a new AST node for the arguments
    ast_t *args_node = ast_new(arena, vec_arena, TRUE);
    if (!args_node)
    {
        // Failed to allocate memory for the arguments node
        return pair_parse_new(NULL, 0);
    }

    // Set the rule type for the arguments node
    args_node->rule = AST_PARAMETERS;

    // Get the current token
    const token_t *token = tokens[0];

    // Check if the token is an open parenthesis
    if (token->type != TOKEN_OPEN_PAREN)
    {
        create_error_ranged(
            token->line,
            token->column,
            token->col_start,
            (ast_rule_t[]){AST_PARAMETERS},
            1
        );
        return pair_parse_new(NULL, 0);
    }

    // Extract all tokens before the closing parenthesis
    const pair_extract_t extract = extract_tokens(
        tokens,
        len,
        TOKEN_OPEN_PAREN,
        TOKEN_CLOSE_PAREN,
        0, // Start at the open paren to avoid nesting issues
        TRUE // Allow nested parentheses
    );

    // Get the extracted tokens and their count
    token_t **extracted_tokens = extract.first;
    const size_t extracted_count = num_max(extract.second, 1) - 1;

    // Handle failure
    if (!extracted_tokens || extracted_count == 0)
    {
        create_error_ranged(
            token->line,
            token->column,
            token->col_start,
            (ast_rule_t[]){AST_PARAMETERS},
            1
        );
        return pair_parse_new(NULL, 0); // Failed to extract tokens
    }

    // Count to skip arguments and commas
    size_t arg_skip_count = 0;
    // Extract all tokens by commas
    while (TRUE)
    {
        // Get the first arg
        const pair_extract_t arg_extract = extract_tokens(
            extracted_tokens + arg_skip_count,
            extracted_count,
            TOKEN_COMMA,
            TOKEN_COMMA,
            0, // Start at the beginning
            FALSE // Do not allow nested commas
        );

        // Get the extracted argument and its count
        token_t **arg_tokens = arg_extract.first;
        size_t arg_count = arg_extract.second;
        bool has_to_break = FALSE;

        // Check if we have reached the end of the arguments
        if (!arg_tokens)
        {
            has_to_break = TRUE; // No more arguments to process
            arg_tokens = extracted_tokens + arg_skip_count; // Use the remaining tokens
            arg_count = extracted_count - arg_skip_count; // Count the remaining tokens
        }
        else
        {
            arg_skip_count += arg_count; // Increment the skip count
        }

        // Create a new node for the arg
        ast_t *arg_node = ast_new(arena, vec_arena, TRUE);
        if (!arg_node)
        {
            create_error_ranged(
                token->line,
                token->column,
                token->col_start,
                (ast_rule_t[]){AST_PARAMETERS},
                1
            );
            return pair_parse_new(NULL, 0); // Failed to allocate memory for the argument node
        }

        // Add the argument node to the arguments node
        vec_ast_push(args_node->children, arg_node);

        // Add the expression to the queue
        alinked_queue_expr_prepend(queue, (queue_expression_t){
            .body = arg_tokens,
            .start = 0, // Start at the beginning of the argument tokens
            .len = arg_count, // Use the count of the argument tokens
            .parent = arg_node // Set the parent to the argument node
        });

        // Check if we have to stop
        if (has_to_break)
        {
            break;
        }
    }

    return pair_parse_new(args_node, extracted_count); // Return the populated arguments node
}

#endif //FLUENT_PARSER_RULE_ARGS_H
