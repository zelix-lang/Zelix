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
// Created by rodrigo on 6/9/25.
//

#ifndef FLUENT_PARSER_RULE_FUNCTION_H
#define FLUENT_PARSER_RULE_FUNCTION_H

// ============= FLUENT LIB C =============
#include <fluent/std_bool/std_bool.h> // fluent_libc
#include <fluent/arena/arena.h> // fluent_libc

// ============= INCLUDES =============
#include "../error.h"
#include "../../ast/ast.h"
#include "../../lexer/stream.h"
#include "../extractor.h"
#include "type.h"
#include "block.h"

// ============ GLOBAL =============
static ast_t nothing_type = {
    .rule = AST_NOTHING,
    .children = NULL,
    .value = NULL,
    .line = 0,
    .column = 0,
    .col_start = 0
};

typedef enum
{
    FN_P_WAKEUP = 0, // Ready to expect the first token
    FN_P_NAME, // Expecting the function name
    FN_P_ARGS_START, // Expecting the opening parenthesis for arguments
    FN_P_ARG_NAME, // Expecting an argument name
    FN_P_COLON, // Expecting a colon after the argument name
    FN_P_ARG_TYPE, // Expecting the type of the argument
    FN_P_TOMBSTONE // Parsing of arguments is complete
} fn_parser_state_t;

static inline bool parse_function(
    const ast_t *const root,
    token_stream_t *const stream,
    arena_allocator_t *const arena,
    arena_allocator_t *const vec_arena,
    const token_t *trace
)
{
    // Define the initial state of the parser
    fn_parser_state_t state = FN_P_WAKEUP;

    // Extract the tokens from the stream
    const pair_extract_t extract = extract_tokens(
        stream,
        TOKEN_OPEN_CURLY,
        TOKEN_CLOSE_CURLY,
        stream->current,
        TRUE // Allow nested delimiters
    );

    // Get the extracted tokens and their count
    token_t **tokens = extract.first;
    const size_t count = extract.second;

    // Check if the extracted chunk is NULL or empty
    if (!tokens || count == 0)
    {
        // Create an error for unexpected end of stream
        create_error(
            trace->line,
            trace->column,
            trace->col_start,
            (ast_rule_t[]){AST_FUNCTION},
            1
        );
        return FALSE;
    }

    // Whether to allow a comma after the last argument
    bool allow_comma = FALSE;
    char *name = NULL; // Function name
    // Allocate a new AST node for the function
    ast_t *function_node = ast_new(arena, vec_arena, TRUE);
    // Have a pointer ready for the params in case we need to allocate any
    ast_t *params_node = NULL;
    // Store the parameter that we are currently parsing
    ast_t *current_param = NULL;
    size_t args_end = 0; // To track the end of the arguments
    ast_t *return_type_node = &nothing_type; // To store the return type of the function
    ast_t *block_node = ast_new(arena, vec_arena, TRUE); // To store the function block

    // Iterate over the extracted tokens
    for (size_t i = 0; i <= count; i++)
    {
        // Get the current token
        const token_t *token = tokens[i];

        // Check the current state of the parser
        switch (state)
        {
            case FN_P_WAKEUP:
            {
                // Expecting the function keyword
                if (token->type != TOKEN_FUNCTION)
                {
                    // Create an error for unexpected token
                    create_error(
                        token->line,
                        token->column,
                        token->col_start,
                        (ast_rule_t[]){AST_FUNCTION},
                        1
                    );
                    return FALSE; // Invalid function declaration
                }

                state = FN_P_NAME; // Move to the next state
                break;
            }

            case FN_P_NAME:
            {
                // Make sure we have a valid identifier
                if (token->type != TOKEN_IDENTIFIER)
                {
                    // Create an error for unexpected token
                    create_error(
                        token->line,
                        token->column,
                        token->col_start,
                        (ast_rule_t[]){AST_FUNCTION},
                        1
                    );
                    return FALSE; // Invalid function name
                }

                // Set the function name
                name = token->value->ptr;

                // Move to the next state
                state = FN_P_ARGS_START;
                break;
            }

            case FN_P_ARGS_START:
            {
                // Make sure we have an opening parenthesis
                if (token->type != TOKEN_OPEN_PAREN)
                {
                    // Create an error for unexpected token
                    create_error(
                        token->line,
                        token->column,
                        token->col_start,
                        (ast_rule_t[]){AST_FUNCTION},
                        1
                    );
                    return FALSE; // Invalid function arguments start
                }

                // Move to the next state
                state = FN_P_ARG_NAME;
                break;
            }

            case FN_P_ARG_NAME:
            {
                // Check for commas
                if (allow_comma && token->type == TOKEN_COMMA)
                {
                    allow_comma = FALSE; // Reset the comma allowance
                    // If we are allowing a comma, skip it and continue
                    continue;
                }

                // Check if we have a closing parenthesis right away
                if (token->type == TOKEN_CLOSE_PAREN)
                {
                    args_end = i; // Mark the end of arguments
                    // No arguments, just return to the next state
                    state = FN_P_TOMBSTONE;
                    break;
                }

                // Make sure we have a valid identifier for the argument name
                if (token->type != TOKEN_IDENTIFIER)
                {
                    // Create an error for unexpected token
                    create_error(
                        token->line,
                        token->column,
                        token->col_start,
                        (ast_rule_t[]){AST_FUNCTION},
                        1
                    );
                    return FALSE; // Invalid argument name
                }

                // Check if we have a params node
                if (params_node == NULL)
                {
                    // Allocate a new AST node for the parameters
                    params_node = ast_new(arena, vec_arena, TRUE);
                    if (!params_node)
                    {
                        // Failed to allocate memory for the parameters node
                        return FALSE;
                    }

                    // Set the rule for the parameters node
                    params_node->rule = AST_PARAMETERS;
                }

                // Allocate a new AST node for the current parameter
                current_param = ast_new(arena, vec_arena, TRUE);
                if (!current_param)
                {
                    // Failed to allocate memory for the current parameter node
                    return FALSE;
                }

                current_param->rule = AST_PARAMETER;

                // Allocate a new AST node for the argument name
                ast_t *arg_name_node = ast_new(arena, vec_arena, FALSE);
                if (!arg_name_node)
                {
                    // Failed to allocate memory for the argument name node
                    return FALSE;
                }

                // Set the rule for the argument name node
                arg_name_node->rule = AST_IDENTIFIER;
                // Set the value for the argument name node
                arg_name_node->value = token->value;
                // Set other metadata
                arg_name_node->line = token->line;
                arg_name_node->column = token->column;
                arg_name_node->col_start = token->col_start;

                // Append the argument name node to the current parameter
                vec_ast_push(current_param->children, arg_name_node);

                // Move to the next state
                state = FN_P_COLON;
                break;
            }

            case FN_P_COLON:
            {
                // Make sure we have a colon after the argument name
                if (token->type != TOKEN_COLON)
                {
                    // Create an error for unexpected token
                    create_error(
                        token->line,
                        token->column,
                        token->col_start,
                        (ast_rule_t[]){AST_FUNCTION},
                        1
                    );
                    return FALSE; // Invalid argument type declaration
                }

                // Move to the next state
                state = FN_P_ARG_TYPE;
                break;
            }

            case FN_P_ARG_TYPE:
            {
                // Parse the type
                const pair_type_parser_t type_parser = parse_type(
                    tokens,
                    i, // Start after the colon
                    count, // Remaining tokens
                    arena,
                    vec_arena
                );

                // Get the extracted range and tokens
                const ast_t *type_node = type_parser.first;
                const size_t skipped_range = type_parser.second;

                // Check if the type node is NULL
                if (!type_node)
                {
                    // Create an error for invalid type
                    create_error(
                        token->line,
                        token->column,
                        token->col_start,
                        (ast_rule_t[]){AST_FUNCTION},
                        1
                    );
                    return FALSE; // Invalid argument type
                }

                // Append the current parameter to the parameters node
                vec_ast_push(params_node->children, current_param);

                // Reset the current parameter
                current_param = NULL;

                // Skip the tokens that were parsed as type
                i = skipped_range;

                // Change the state to the next one
                state = FN_P_ARG_NAME;
                allow_comma = TRUE; // Allow a comma after the argument type
            }

            default:
            {
                // Impossible case
            }
        }
    }

    // Skip the argument count in the token stream
    stream->current += args_end;

    // Get the next token
    token_t *token = token_stream_next(stream);

    // Make sure we have a token
    if (!token)
    {
        return FALSE; // Failed to parse the function
    }

    // Check if we have a return type
    if (token->type == TOKEN_ARROW)
    {
        // Position the type parser at the token
        // after the arrow token
        token_stream_next(stream);

        // Parse the return type
        const pair_type_parser_t return_type_parser = parse_type(
            tokens,
            stream->current, // Start after the arrow
            count, // Remaining tokens
            arena,
            vec_arena
        );

        // Get the extracted range and tokens
        return_type_node = return_type_parser.first;
        const size_t skipped_range = return_type_parser.second;

        // Handle parsing failure
        if (return_type_node == NULL)
        {
            // Create an error for invalid return type
            create_error(
                token->line,
                token->column,
                token->col_start,
                (ast_rule_t[]){AST_FUNCTION},
                1
            );
            return FALSE; // Invalid return type
        }

        // Skip the parsed range
        args_end = skipped_range - 1;
        stream->current = skipped_range;

        // Peek the next token
        token = token_stream_peek(stream);
    }
    else
    {
        // Skip 2 tokens for the close parenthesis and open curly
        args_end += 2;
    }

    // Make sure that we have an open curly brace
    if (!token || token->type != TOKEN_OPEN_CURLY)
    {
        // Get the current token
        token = token_stream_nth(stream, stream->current);

        // Create an error for unexpected token
        create_error(
            token->line,
            token->column,
            token->col_start,
            (ast_rule_t[]){AST_FUNCTION},
            1
        );
        return FALSE; // Invalid function body start
    }

    // Consume the opening curly brace
    token_stream_next(stream);

    // Extract the function's body
    token_t **body = tokens + args_end;
    const size_t body_len = count - args_end - 1;

    // Parse the block
    if (
        !parse_block(
            block_node,
            body,
            body_len,
            arena,
            vec_arena
        )
    )
    {
        return FALSE; // Failed to parse the function block
    }

    // Skip the token count
    // Avoid positioning exactly the token next to the closing curly brace
    // since the main loop will increment the current position
    stream->current += count;

    // Set the function node metadata
    function_node->rule = AST_FUNCTION; // Set the rule for the function node

    // Create a new AST node for the function name
    ast_t *name_node = ast_new(arena, vec_arena, FALSE);
    if (!name_node)
    {
        // Failed to allocate memory for the function name node
        return FALSE;
    }

    // Set the rule for the function name node
    name_node->rule = AST_IDENTIFIER;
    // Set the value for the function name node
    name_node->value = heap_str_alloc(FALSE, FALSE, NULL, name);

    // Append the parsed nodes to the function node's children
    vec_ast_push(function_node->children, name_node);
    vec_ast_push(function_node->children, params_node); // Add the parameters node
    vec_ast_push(function_node->children, return_type_node); // Add the return type node
    vec_ast_push(function_node->children, block_node); // Add the block node

    // Add the function node to the root's children
    vec_ast_push(root->children, function_node);

    // Return success
    return TRUE;
}

#endif //FLUENT_PARSER_RULE_FUNCTION_H
