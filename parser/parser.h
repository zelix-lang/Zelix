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
#include <parser/error.h>
#include <lexer/stream.h>
#include <parser/stream.h>
#include <parser/rule/import.h>
#include <parser/rule/function.h>

DEFINE_PAIR_T(ast_stream_t, ast_error_t *, parser_result);

static pair_parser_result_t create_failed_result(const ast_stream_t ast_stream)
{
    // Create a pair with the AST stream and the error
    return (pair_parser_result_t){
        .first = ast_stream,
        .second = &global_parser_error
    };
}

static inline pair_parser_result_t parser_parse(
    token_stream_t *const stream,
    const char *const file_name
)
{
    // Emit parsing state
    new_timer(file_name, STATE_PARSING);
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
                if (!parse_import(
                    ast_stream.ast,
                    stream,
                    ast_stream.allocator,
                    ast_stream.vec_allocator,
                    current
                ))
                {
                    // Failed to parse the import statement
                    return create_failed_result(ast_stream);
                }

                break;
            }

            case TOKEN_FUNCTION:
            {
                parse_function(
                    ast_stream.ast,
                    stream,
                    ast_stream.allocator,
                    ast_stream.vec_allocator,
                    current
                );
                break;
            }

            default:
            {
                create_error(
                    current->line,
                    current->column,
                    current->col_start,
                    (ast_rule_t[]){AST_IMPORT, AST_FUNCTION, AST_MODULE},
                    3
                );

                return create_failed_result(ast_stream);
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
