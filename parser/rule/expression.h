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
    while (queue.len > 0) {
        // Get the first element of the queue
        const queue_expression_t current = alinked_queue_expr_shift(&queue);

        // Get the body and its length
        token_t **input = current.body + current.start;
        size_t input_len = current.len - current.start;
        const ast_t *parent = current.parent;

        // Handle empty input
        if (input_len == 1) // Consider len is inclusive (includes the semicolon)
        {
            create_error_ranged(
                input[0]->line,
                input[0]->column,
                input[0]->col_start,
                (ast_rule_t[]){AST_EXPRESSION},
                1
            );

            return FALSE;
        }

        // Parse pointers and dereferences
        const token_t *token = input[0];
        while (
            token->type == TOKEN_ASTERISK
            || token->type == TOKEN_AMPERSAND // Single pointer (&)
            || token->type == TOKEN_AND // Double pointer (&&)
        )
        {
            // Create a new AST node for the pointer or dereference
            ast_t *ptr_node = ast_new(arena, vec_arena, FALSE);
            if (!ptr_node)
            {
                // Failed to allocate memory for the pointer node
                return FALSE;
            }

            // Set the rule based on the token type
            switch (token->type)
            {
                case TOKEN_ASTERISK:
                    ptr_node->rule = AST_DEREFERENCE;
                    break;
                case TOKEN_AMPERSAND:
                    ptr_node->rule = AST_POINTER;
                    break;
                case TOKEN_AND:
                    // Case managed later
                    ptr_node->rule = AST_POINTER;
                    break;
                default:
                    // Invalid token type for pointer or dereference
                    return FALSE;
            }

            // Set the value of the pointer node
            ptr_node->value = token->value;

            // Add the pointer node to the expression's children
            vec_ast_push(parent->children, ptr_node);

            // Handle double pointers
            if (token->type == TOKEN_AND)
            {
                // Push the same node without allocating more memory
                vec_ast_push(parent->children, ptr_node);
            }

            // Move to the next token
            input++;

            // Decrement the input length
            input_len--;

            // Update the token
            token = input[0];
        }

        // Handle invalid expressions
        if (input_len == 1)
        {
            return FALSE; // Invalid expression with only one token
        }

        // Define the candidate
        size_t start_at = 0;
        ast_t *candidate = NULL;

        // Parse the candidate
        // Check for nested expressions
        if (token->type == TOKEN_OPEN_PAREN)
        {
            // Extract all tokens until the closing parenthesis
            const pair_extract_t extract = extract_tokens(
                input,
                input_len,
                TOKEN_OPEN_PAREN,
                TOKEN_CLOSE_PAREN,
                0,
                TRUE
            );

            // Get the extracted range
            token_t **range = extract.first;
            const size_t len = extract.second;

            // Handle failure
            if (!range)
            {
                create_error_ranged(
                    token->line,
                    token->column,
                    token->col_start,
                    (ast_rule_t[]){AST_EXPRESSION},
                    1
                );
                return FALSE; // Failed to extract tokens
            }

            // Create a new AST node for the nested expression
            ast_t *nested_expression = ast_new(arena, vec_arena, TRUE);
            if (!nested_expression)
            {
                // Failed to allocate memory for the nested expression
                return FALSE;
            }

            // Set the rule for the nested expression
            nested_expression->rule = AST_EXPRESSION;

            // Update the candidate to the nested expression
            candidate = nested_expression;

            // Enqueue the nested expression
            alinked_queue_expr_append(&queue, (queue_expression_t){
                .body = range + 1, // Skip the opening parenthesis
                .start = 0,
                .len = len - 1, // Skip the closing parenthesis
                .parent = nested_expression
            });

            start_at += len; // Move the start index forward
        }

        // Check if we have more tokens to process
        if (start_at == input_len)
        {
            // Make sure we have a valid candidate
            if (!candidate)
            {
                // Create a new candidate from the single token
                candidate = parse_single_token(input[0], arena, vec_arena);
                if (!candidate)
                {
                    create_error_ranged(
                        input[0]->line,
                        input[0]->column,
                        input[0]->col_start,
                        (ast_rule_t[]){AST_EXPRESSION},
                        1
                    );
                    return FALSE; // Failed to parse the single token
                }
            }

            // Append the candidate to the parent
            vec_ast_push(parent->children, candidate);
            continue; // No more tokens to process
        }

        // Move the buffer using pointer arithmetic
        input += start_at;
        input_len -= start_at; // Decrement the input length

        // Check if we have a valid candidate
        if (!candidate)
        {
            // Make sure we have enough tokens to parse
            if (input_len < 1)
            {
                create_error_ranged(
                    token->line,
                    token->column,
                    token->col_start,
                    (ast_rule_t[]){AST_EXPRESSION},
                    1
                );
                return FALSE; // Not enough tokens to parse
            }

            // Parse the single token as a candidate
            candidate = parse_single_token(input[0], arena, vec_arena);
            if (!candidate)
            {
                return FALSE; // Irrecoverable state
            }
        }

        // Move the pointer forward and decrement the input length
        input++;
        input_len--;

        // Check if we have any tokens left
        if (input_len == 0)
        {
            // Append the candidate to the parent
            vec_ast_push(parent->children, candidate);
            continue; // No more tokens to process
        }

        // Update the first token
        token = input[0];

        // Check for function calls
        if (token->type == TOKEN_OPEN_PAREN)
        {
            // TODO: Args parser
        }

        // Check for arithmetic operations
        if (
            token->type == TOKEN_PLUS ||
            token->type == TOKEN_MINUS ||
            token->type == TOKEN_ASTERISK ||
            token->type == TOKEN_SLASH
        )
        {
            // TODO: Arithmetic parser
        }

        // Check for prop access
        if (token->type == TOKEN_DOT)
        {
            // TODO: Prop access parser
        }

        // Check for boolean operations
        if (
            token->type == TOKEN_EQUAL ||
            token->type == TOKEN_NOT_EQUAL ||
            token->type == TOKEN_GREATER_THAN ||
            token->type == TOKEN_LESS_THAN ||
            token->type == TOKEN_LESS_THAN_OR_EQUAL ||
            token->type == TOKEN_GREATER_THAN_OR_EQUAL
        )
        {
            // TODO: Boolean parser
        }

        // Check if we have tokens left
        if (start_at != input_len)
        {
            token = input[start_at];
            // Invalid expression
            create_error_ranged(
                token->line,
                token->column,
                token->col_start,
                (ast_rule_t[]){AST_EXPRESSION},
                1
            );
            return FALSE;
        }

        // Add the final candidate to the parent
        vec_ast_push(parent->children, candidate);
    }

    // Add the child to the root
    vec_ast_push(root->children, expression);
    return TRUE;
}

#endif //FLUENT_PARSER_RULE_EXPRESSION_H
