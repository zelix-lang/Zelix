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
#include <ast/ast.h>
#include <token/token.h>
#include <parser/queue/expression.h>
#include <parser/rule/new.h>

static ast_t *parse_single_token(
    const token_t *token,
    arena_allocator_t *const arena,
    arena_allocator_t *const vec_arena
)
{
    // Create a new AST node for the single token
    ast_t *node = ast_new(arena, vec_arena, FALSE);
    if (!node)
    {
        // Failed to allocate memory for the node
        return NULL;
    }

    // Switch on the token type to set the appropriate rule
    switch (token->type)
    {
        case TOKEN_IDENTIFIER:
            node->rule = AST_IDENTIFIER;
            break;
        case TOKEN_STRING_LITERAL:
            node->rule = AST_STRING_LITERAL;
            break;
        case TOKEN_NUM_LITERAL:
            node->rule = AST_NUMBER_LITERAL;
            break;
        case TOKEN_DECIMAL_LITERAL:
            node->rule = AST_DECIMAL_LITERAL;
            break;
        case TOKEN_BOOL_LITERAL:
            node->rule = AST_BOOLEAN_LITERAL;
            break;
        default:
            // Invalid token type for a single token expression
            return NULL;
    }

    // Set the token's value
    node->value = token->value;
    return node;
}

static inline bool parse_expression(
    const ast_t *const root,
    token_t **body,
    const size_t body_start,
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
        .start = body_start,
        .len = body_len - body_start,
        .parent = expression
    });

    // Iterate over the queue until it's empty
    while (queue.len > 0)
    {
        // Get the first element
        const queue_expression_t expr = alinked_queue_expr_shift(&queue);
        token_t **input = expr.body + expr.start;
        const size_t len = expr.len - expr.start;
        const ast_t *parent = expr.parent;
        size_t start = 0;

        // Handle empty input
        if (len == 0)
        {
            return FALSE;
        }

        // Parse pointers
        for (size_t i = 0; i <= len; i++)
        {
            // Get the current token
            const token_t *token = input[i];
            ast_rule_t rule = AST_PROGRAM_RULE;

            // Check if we have a dereference operation
            if (token->type == TOKEN_ASTERISK)
            {
                rule = AST_DEREFERENCE; // Set the rule to dereference
            }
            // Check for pointers
            else if (
                token->type == TOKEN_AMPERSAND ||
                token->type == TOKEN_AND
            )
            {
                rule = AST_POINTER; // Set the rule to pointer
            }
            else
            {
                break;
            }

            // Increment the start counter
            start++;

            // Allocate a new AST node for the pointer or dereference
            ast_t *ptr_node = ast_new(arena, vec_arena, FALSE);
            if (!ptr_node)
            {
                // Failed to allocate memory for the pointer node
                return FALSE;
            }

            // Set the metadata for the pointer node
            ptr_node->rule = rule;
            ptr_node->line = token->line;
            ptr_node->column = token->column;
            ptr_node->col_start = token->col_start;

            // Append the pointer node to the expression
            vec_ast_push(parent->children, ptr_node);

            // Check for double pointers
            if (token->type == TOKEN_AND)
            {
                // Append the same node for a double pointer
                // without allocating a new one
                vec_ast_push(parent->children, ptr_node);
            }
        }

        // Handle invalid expressions
        if (start == len)
        {
            return FALSE; // Invalid expression, no valid tokens found
        }

        // Get the current token
        token_t *token = input[start];
        ast_t *candidate = NULL;
        // Whether the expression is likely to be arithmetic
        bool is_arithmetic = FALSE;
        // Whether the expression is likely to be a prop access
        bool is_prop_access = FALSE;

        // Look for nested expressions
        if (token->type == TOKEN_OPEN_PAREN)
        {
            // Extract all tokens before the next close paren
            const pair_extract_t nested_expr_result = extract_tokens(
                input,
                len,
                TOKEN_OPEN_PAREN,
                TOKEN_CLOSE_PAREN,
                start,
                TRUE
            );

            // Obtain the extracted range and len
            token_t **nested_expr = nested_expr_result.first;
            const size_t skipped = nested_expr_result.second;

            // Handle failure
            if (!nested_expr)
            {
                return FALSE;
            }

            // Create a new AST node for the nested expression
            ast_t *nested_expr_node = ast_new(arena, vec_arena, TRUE);
            if (!nested_expr_node)
            {
                // Failed to allocate memory for the nested expression
                return FALSE;
            }

            // Set the rule for the nested expression
            nested_expr_node->rule = AST_EXPRESSION;

            // Create a new queue element and add it to the queue
            alinked_queue_expr_append(&queue, (queue_expression_t){
                .body = nested_expr,
                .start = 0,
                .len = skipped,
                .parent = nested_expr_node
            });

            // Skip the nested expression
            start += skipped;

            // Update the candidate to the nested expression node
            candidate = nested_expr_node;
        }

        // Handle literals and identifiers
        if (candidate == NULL)
        {
            // Parse object creation tokens
            if (token->type == TOKEN_NEW)
            {
                // Parse the object creation expression
                const pair_obj_creation_t obj_creation_result = parse_new(
                    &queue,
                    input,
                    start + 1, // +1 to start after the 'new' token
                    len,
                    arena,
                    vec_arena
                );

                // Get the extracted object creation node and its length
                candidate = obj_creation_result.first;
                const size_t obj_creation_len = obj_creation_result.second;

                if (!candidate)
                {
                    create_error(
                        token->line,
                        token->column,
                        token->col_start,
                        (ast_rule_t[]){AST_FUNCTION_CALL},
                        1
                    );

                    // Failed to parse the object creation expression
                    return FALSE;
                }

                // Update the start counter
                start += obj_creation_len + 1; // +1 to skip the 'new' token
                is_arithmetic = FALSE;
                is_prop_access = TRUE;
            }
            else
            {
                // Update the start counter
                start++;

                // Parse a single token
                ast_t *single = parse_single_token(token, arena, vec_arena);
                if (!single)
                {
                    return FALSE;
                }

                // Update the candidate
                candidate = single;
                is_arithmetic = single->rule != AST_BOOLEAN_LITERAL
                    && single->rule != AST_STRING_LITERAL;
                is_prop_access = single->rule == AST_IDENTIFIER;
            }
        }

        // Handle end of the expression
        if (start == len)
        {
            // Append the nested expression node to the parent
            vec_ast_push(parent->children, candidate);
            continue; // No more tokens to process
        }

        // Update the token to the next one
        token = input[start];

        // Check for prop access
        if (token->type == TOKEN_DOT)
        {
            // Check if we are allowed to parse prop access
            if (!is_prop_access)
            {
                // Create an error for unexpected prop access
                create_error(
                    token->line,
                    token->column,
                    token->col_start,
                    (ast_rule_t[]){AST_EXPRESSION},
                    1
                );
                return FALSE; // Invalid prop access expression
            }

            // TODO: Prop access parser
        }

        // Check for arithmetic operations
        if (
            token->type == TOKEN_PLUS
            || token->type == TOKEN_MINUS
            || token->type == TOKEN_ASTERISK
            || token->type == TOKEN_SLASH
        )
        {
            // Check if we are allowed to parse arithmetic expressions
            if (!is_arithmetic)
            {
                // Create an error for unexpected arithmetic operation
                create_error(
                    token->line,
                    token->column,
                    token->col_start,
                    (ast_rule_t[]){AST_EXPRESSION},
                    1
                );
                return FALSE; // Invalid arithmetic expression
            }

            // TODO: Arithmetic parser
        }

        // Check for comparison operations
        if (
            token->type == TOKEN_EQUAL
            || token->type == TOKEN_NOT_EQUAL
            || token->type == TOKEN_LESS_THAN
            || token->type == TOKEN_GREATER_THAN
            || token->type == TOKEN_LESS_THAN_OR_EQUAL
            || token->type == TOKEN_GREATER_THAN_OR_EQUAL
        )
        {
            // TODO: boolean parser
        }

        // Handle end of the expression
        if (start == len)
        {
            // Append the nested expression node to the parent
            vec_ast_push(parent->children, candidate);
            continue; // No more tokens to process
        }

        // No more tokens to process, expression is invalid
        create_error(
            token->line,
            token->column,
            token->col_start,
            (ast_rule_t[]){AST_EXPRESSION},
            1
        );

        return FALSE; // Invalid expression
    }

    // Add the child to the root
    vec_ast_push(root->children, expression);
    return TRUE;
}

#endif //FLUENT_PARSER_RULE_EXPRESSION_H
