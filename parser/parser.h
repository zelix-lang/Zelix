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

#ifndef FLUENT_PARSER_H
#define FLUENT_PARSER_H

// ============= FLUENT LIB C =============
#include <fluent/pair/pair.h> // fluent_libc
#include <fluent/std_bool/std_bool.h> // fluent_libc

// ============= INCLUDES =============
#include "error.h"
#include "../lexer/stream.h"
#include "stream.h"

// ============= GLOBAL VARIABLE =============
static ast_error_t global_parser_error;

static ast_error_t *create_error(
    const size_t line,
    const size_t column,
    const size_t col_start,
    const ast_rule_t *const expected,
    const size_t expected_len
)
{
    // Set the global parser error
    global_parser_error.line = line;
    global_parser_error.column = column;
    global_parser_error.col_start = col_start;

    // Copy the expected rules into the global parser error
    memcpy(global_parser_error.expected, expected, expected_len);

    return &global_parser_error;
}

DEFINE_PAIR_T(ast_stream_t, ast_error_t *, parser_result);

static inline pair_parser_result_t parser_parse(
    token_stream_t *const stream,
    const char *const file_name
)
{
    // Create a new AST stream
    ast_stream_t ast_stream;

    // Create a new arena allocator for the AST
    ast_stream.allocator = arena_new(50, sizeof(ast_t));
    // Handle allocation failure
    if (!ast_stream.allocator)
    {
        // Failed to create the arena allocator
        return (pair_parser_result_t){NULL, NULL};
    }

    // Create a new arena allocator for the vectors
    ast_stream.vec_allocator = arena_new(50, sizeof(vector_generic_t));
    // Handle allocation failure
    if (!ast_stream.vec_allocator)
    {
        // Failed to create the arena allocator for vectors
        destroy_arena(ast_stream.allocator);
        return (pair_parser_result_t){NULL, NULL};
    }

    // Allocate the root AST node
    ast_stream.ast = ast_new(ast_stream.allocator, ast_stream.vec_allocator, TRUE);
    // Handle allocation failure
    if (!ast_stream.ast)
    {
        // Failed to allocate the root AST node
        destroy_arena(ast_stream.allocator);
        return (pair_parser_result_t){NULL, NULL};
    }

    // Initialize the root AST node
    ast_stream.ast->rule = AST_PROGRAM_RULE; // Set the rule to PROGRAM

    // Iterate over the token stream
    const token_t *current = token_stream_nth(stream, 0);
    while (current != NULL)
    {
        // Determine what we have to parse
        switch (current->type)
        {
            case TOKEN_IMPORT:
            {
                break;
            }

            case TOKEN_PUB:
            case TOKEN_FUNCTION:
            case TOKEN_MOD:
            {
                break;
            }

            default:
            {
                return (pair_parser_result_t){
                    ast_stream,
                     create_error(
                        current->line,
                        current->column,
                        current->col_start,
                        (ast_rule_t[]){AST_IMPORT, AST_FUNCTION, AST_MODULE},
                        3
                     )
                };
            }
        }

        // Move to the next token
        current = token_stream_next(stream);
    }

    // Return the parser result
    return (pair_parser_result_t){
        .first = ast_stream,
        .second = NULL
    };
}

#endif //FLUENT_PARSER_H
