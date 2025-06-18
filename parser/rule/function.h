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
#include <parser/error.h>
#include <ast/ast.h>
#include <lexer/stream.h>
#include <parser/rule/type.h>

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
    FN_P_ARG_NAME = 0, // Expecting an argument name
    FN_P_COLON, // Expecting a colon after the argument name
    FN_P_ARG_TYPE, // Expecting the type of the argument
    FN_P_TOMBSTONE // Parsing commas or closing parenthesis
} fn_parser_state_t;

static inline ast_t *parse_function(
    const ast_t *const root,
    token_stream_t *const stream,
    arena_allocator_t *const arena,
    arena_allocator_t *const vec_arena
)
{
    // Create a new AST node for the function
    ast_t *function_node = ast_new(arena, vec_arena, TRUE);
    if (!function_node)
    {
        // Failed to allocate memory for the function node
        return NULL;
    }

    // Get the next token
    const token_t *token = token_stream_next(stream);
    if (!token)
    {
        create_error(
            stream,
            (ast_rule_t[]){AST_FUNCTION},
            1
        );

        return NULL; // Failed to get the next token
    }

    // Make sure the token is an identifier
    if (token->type != TOKEN_IDENTIFIER)
    {
        // Create an error for unexpected token
        create_error(
            stream,
            (ast_rule_t[]){AST_IDENTIFIER},
            1
        );
        return NULL; // Invalid function name
    }

    // Create a new identifier node
    ast_t *name_node = ast_new(arena, vec_arena, FALSE);
    if (!name_node)
    {
        // Failed to allocate memory for the identifier node
        return NULL;
    }

    name_node->rule = AST_IDENTIFIER; // Set the rule to identifier
    name_node->value = token->value; // Set the value to the identifier
    vec_ast_push(function_node->children, name_node); // Add the identifier to the function node

    // Get the next token
    token = token_stream_next(stream);
    if (!token)
    {
        create_error(
            stream,
            (ast_rule_t[]){AST_FUNCTION},
            1
        );

        return NULL; // Failed to get the next token
    }

    // Check if the next token is an opening parenthesis
    if (token->type != TOKEN_OPEN_PAREN)
    {
        // Create an error for unexpected token
        create_error(
            stream,
            (ast_rule_t[]){AST_FUNCTION},
            1
        );
        return NULL; // Invalid function signature
    }

    // Move to the next token
    token = token_stream_next(stream);

    // Check if we have a closing parenthesis
    if (token->type == TOKEN_CLOSE_PAREN)
    {
        // Return the function node without parameters
        return function_node;
    }

    // Create a new parameters node
    ast_t *params_node = ast_new(arena, vec_arena, TRUE);
    if (!params_node)
    {
        // Failed to allocate memory for the identifier node
        return NULL;
    }

    // Set the metadata for the node
    params_node->rule = AST_PARAMETERS;
    vec_ast_push(function_node->children, params_node);

    // Set the parser state to expect arguments
    fn_parser_state_t state = FN_P_ARG_NAME;
    ast_t *param_name = NULL;

    // Parse all arguments
    while (token != NULL)
    {
        bool has_to_break = FALSE;

        // Switch on the state to determine what to do
        switch (state)
        {
            case FN_P_ARG_NAME:
            {
                // Make sure the token is an identifier
                if (token->type != TOKEN_IDENTIFIER)
                {
                    create_error(
                        stream,
                        (ast_rule_t[]){AST_FUNCTION},
                        1
                    );
                    return NULL;
                }

                // Create a new node for the identifier
                param_name = ast_new(arena, vec_arena, FALSE);
                if (!param_name)
                {
                    create_error(
                        stream,
                        (ast_rule_t[]){AST_FUNCTION},
                        1
                    );
                    return NULL;
                }

                // Update the param name
                param_name->value = token->value;
                // Move to the next state
                state = FN_P_COLON;
                break;
            }

            case FN_P_COLON:
            {
                // Make sure the token is a colon
                if (token->type != TOKEN_COLON)
                {
                    create_error(
                        stream,
                        (ast_rule_t[]){AST_PARAMETER},
                        1
                    );
                    return NULL;
                }

                // Move to the next state
                state = FN_P_ARG_TYPE;
                break;
            }

            case FN_P_ARG_TYPE:
            {
                // Parse the argument type
                ast_t *type = parse_type(stream, arena, vec_arena);

                // Handle failure
                if (!type)
                {
                    create_error(
                        stream,
                        (ast_rule_t[]){AST_PARAMETER},
                        1
                    );
                    return NULL;
                }

                // Create a new param node
                ast_t *param = ast_new(arena, vec_arena, TRUE);
                if (!param)
                {
                    create_error(
                        stream,
                        (ast_rule_t[]){AST_FUNCTION},
                        1
                    );
                    return NULL;
                }

                // Append the name and the type to the param
                vec_ast_push(param->children, param_name);
                vec_ast_push(param->children, type);
                vec_ast_push(params_node->children, param);

                // Move to the next state
                state = FN_P_TOMBSTONE;
                break;
            }

            case FN_P_TOMBSTONE:
            {
                // Check if we have reached the end of the parameters
                if (token->type == TOKEN_CLOSE_PAREN)
                {
                    has_to_break = TRUE;
                    break;
                }

                // Make sure we have a comma
                if (token->type != TOKEN_COMMA)
                {
                    create_error(
                        stream,
                        (ast_rule_t[]){AST_FUNCTION},
                        1
                    );
                    return NULL;
                }

                // Move to the name state again
                state = FN_P_ARG_NAME;
                break;
            }
        }

        // Check if we have to break
        if (has_to_break)
        {
            break;
        }

        // Move to the next token
        token = token_stream_next(stream);
    }

    // Make sure we end parsing a name
    if (state != FN_P_ARG_NAME)
    {
        create_error(
            stream,
            (ast_rule_t[]){AST_FUNCTION},
            1
        );
        return NULL;
    }

    // Get the next token
    token = token_stream_next(stream);

    // Check if we have an arrow
    if (!token)
    {
        create_error(
            stream,
            (ast_rule_t[]){AST_FUNCTION},
            1
        );
        return NULL;
    }

    // Check if we have to parse a return type
    if (token->type == TOKEN_ARROW)
    {
        // Consume the arrow token
        token_stream_next(stream);

        // Parse the return type
        ast_t *type = parse_type(stream, arena, vec_arena);

        // Handle failure
        if (!type)
        {
            create_error(
                stream,
                (ast_rule_t[]){AST_PARAMETER},
                1
            );
            return NULL;
        }

        // Append the return type
        vec_ast_push(function_node->children, type);

        // Move to the next token
        token = token_stream_next(stream);
    }
    else
    {
        // Append the nothing return type
        vec_ast_push(function_node->children, &nothing_type);
    }

    // Make sure the token type is an open curly
    if (token->type != TOKEN_OPEN_CURLY)
    {
        create_error(
            stream,
            (ast_rule_t[]){AST_FUNCTION},
            1
        );
        return NULL;
    }

    // Create a new block AST node
    ast_t *block = ast_new(arena, vec_arena, TRUE);
    if (!block)
    {
        return NULL;
    }

    // Append the block node to the function
    vec_ast_push(function_node->children, block);

    // Append the function node to the root
    vec_ast_push(root->children, function_node);

    // Return the node
    return block;
}

#endif //FLUENT_PARSER_RULE_FUNCTION_H
