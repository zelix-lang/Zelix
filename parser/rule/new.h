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
// Created by rodrigo on 6/18/25.
//

#ifndef FLUENT_PARSER_RULE_NEW_H
#define FLUENT_PARSER_RULE_NEW_H

// ============= FLUENT LIB C =============
#include <fluent/pair/pair.h> // fluent_libc

#include "parser/error.h"
DEFINE_PAIR_T(ast_t *, size_t, obj_creation);

static inline pair_obj_creation_t parse_args(
    alinked_queue_expr_t *queue,
    token_t **body,
    size_t start,
    const size_t len,
    arena_allocator_t *const arena,
    arena_allocator_t *const vec_arena
)
{
    const token_t *open_paren = body[start];

    // Validate the token
    if (open_paren->type != TOKEN_OPEN_PAREN)
    {
        // Create an error for unexpected tokens
        create_error_ranged(
            open_paren->line,
            open_paren->column,
            open_paren->col_start,
            (ast_rule_t[]){AST_IDENTIFIER, AST_FUNCTION_CALL},
            2
        );
        return pair_obj_creation_new(NULL, 0);
    }

    // Create a new parameters node
    ast_t *params_node = ast_new(arena, vec_arena, TRUE);
    if (!params_node)
    {
        // Failed to allocate memory for the parameters node
        return pair_obj_creation_new(NULL, 0);
    }

    // Get the 2rd token to determine if we have to parse arguments
    const token_t *close_paren = body[start + 1];
    if (close_paren->type == TOKEN_CLOSE_PAREN)
    {
        // No arguments have to be parsed, just return the object creation node
        return pair_obj_creation_new(NULL, 3);
    }

    // Skip 1 token for the open parenthesis
    start += 2; // + 1 to start at the first argument

    // Parse all arguments by extracting tokens
    // until we reach the closing parenthesis
    while (TRUE)
    {
        bool has_to_break = FALSE;
        // Extract all tokens before the next comma
        const pair_extract_t extracted = extract_tokens(
            body,
            len,
            TOKEN_COMMA,
            TOKEN_COMMA,
            start, // Start after the open parenthesis
            FALSE
        );

        // Get the extracted tokens and their count
        token_t **arg = extracted.first;
        size_t arg_len = extracted.second;

        // Check if the extraction failed
        if (!arg)
        {
            // This means we are parsing the last param
            arg = body + start; // Set the argument to the remaining tokens
            // Set the length to the remaining tokens
            arg_len = len - start - 1; // -1 to exclude the closing parenthesis
            has_to_break = TRUE; // We have to break the loop after this
        }

        // Create a new parameter AST node
        ast_t *param_node = ast_new(arena, vec_arena, TRUE);
        if (!param_node)
        {
            // Failed to allocate memory for the parameter node
            return pair_obj_creation_new(NULL, 0);
        }

        // Set the rule for the parameter node
        param_node->rule = AST_PARAMETER;

        // Create a new AST node for the expression
        ast_t *expression_node = ast_new(arena, vec_arena, TRUE);
        if (!expression_node)
        {
            // Failed to allocate memory for the expression node
            return pair_obj_creation_new(NULL, 0);
        }

        // Set the rule for the expression node
        expression_node->rule = AST_EXPRESSION;

        // Add the expression node to the parameter node
        vec_ast_push(param_node->children, expression_node);

        // Add the parameter node to the params node
        vec_ast_push(params_node->children, param_node);

        // Add the expression to the queue
        alinked_queue_expr_append(queue, (queue_expression_t){
            .body = arg,
            .start = 0, // Start at the beginning of the argument
            .len = arg_len, // Length of the argument
            .parent = expression_node // Parent is the expression node
        });

        // Break the loop if we have reached the closing parenthesis
        if (has_to_break)
        {
            // Check if the next token is a closing parenthesis
            if (body[start + arg_len]->type != TOKEN_CLOSE_PAREN)
            {
                // Create an error for unexpected end of stream
                create_error_ranged(
                    body[start + arg_len]->line,
                    body[start + arg_len]->column,
                    body[start + arg_len]->col_start,
                    (ast_rule_t[]){AST_FUNCTION_CALL},
                    1
                );
                return pair_obj_creation_new(NULL, 0);
            }
            break; // Break the loop since we have reached the closing parenthesis
        }
    }

    return pair_obj_creation_new(params_node, start + 1); // +1 to include the closing parenthesis
}

static inline pair_obj_creation_t parse_new(
    alinked_queue_expr_t *queue,
    token_t **body,
    const size_t start,
    const size_t len,
    arena_allocator_t *const arena,
    arena_allocator_t *const vec_arena
)
{
    // Perform bounds checking
    if (start + 2 > len)
    {
        // Create an error for unexpected end of stream
        return pair_obj_creation_new(NULL, 0);
    }

    // Allocate a new AST node for the object creation
    ast_t *object_creation_node = ast_new(arena, vec_arena, TRUE);

    // Handle allocation failure
    if (!object_creation_node)
    {
        // Failed to allocate memory for the node
        return pair_obj_creation_new(NULL, 0);
    }

    // Get the identifier
    const token_t *identifier_token = body[start];

    // Validate the tokens
    if (identifier_token->type != TOKEN_IDENTIFIER)
    {
        // Create an error for unexpected tokens
        create_error_ranged(
            identifier_token->line,
            identifier_token->column,
            identifier_token->col_start,
            (ast_rule_t[]){AST_IDENTIFIER, AST_FUNCTION_CALL},
            2
        );
        return pair_obj_creation_new(NULL, 0);
    }

    // Create a new node for the identifier
    ast_t *identifier_node = ast_new(arena, vec_arena, FALSE);
    if (!identifier_node)
    {
        // Failed to allocate memory for the identifier node
        return pair_obj_creation_new(NULL, 0);
    }

    identifier_node->rule = AST_IDENTIFIER; // Set the rule to identifier
    identifier_node->value = identifier_token->value; // Set the value to the identifier
    vec_ast_push(object_creation_node->children, identifier_node);

    // Parse the parameters
    const pair_obj_creation_t params_result = parse_args(
        queue,
        body,
        start + 1, // +1 to skip the 'new' token
        len,
        arena,
        vec_arena
    );

    // Get the parameters node and its length
    ast_t *params_node = params_result.first;
    const size_t params_len = params_result.second;

    // Handle failure
    if (!params_node && params_len == 0)
    {
        return pair_obj_creation_new(NULL, 0);
    } else if (!params_node)
    {
        return pair_obj_creation_new(object_creation_node, 3);
    }

    // Add the parameters node to the object creation node
    vec_ast_push(object_creation_node->children, params_node);

    return pair_obj_creation_new(object_creation_node, len);
}

#endif //FLUENT_PARSER_RULE_NEW_H
