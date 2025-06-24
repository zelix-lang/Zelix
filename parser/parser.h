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
#include <parser/extractor.h>
#include <parser/rule/expression.h>

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

    // Track the blocks being parsed
    alinked_queue_ast_t blocks;
    alinked_queue_ast_init(&blocks, 10);
    bool in_mod = FALSE;

    while (current != NULL)
    {
        // Determine what we have to parse
        switch (current->type)
        {
            case TOKEN_IMPORT:
            {
                if (
                    !parse_import(
                        ast_stream.ast,
                        stream,
                        ast_stream.allocator,
                        ast_stream.vec_allocator
                    )
                )
                {
                    // Failed to parse the import statement
                    return create_failed_result(ast_stream);
                }

                break;
            }

            case TOKEN_FUNCTION:
            {
                // Make sure the declaration is valid
                if (!in_mod && blocks.len != 0)
                {
                    create_error(
                        stream,
                        (ast_rule_t[]){AST_IMPORT, AST_FUNCTION, AST_MODULE},
                        3
                    );

                    return create_failed_result(ast_stream);
                }

                // Parse the function
                ast_t *new_block = parse_function(
                    ast_stream.ast,
                    stream,
                    ast_stream.allocator,
                    ast_stream.vec_allocator
                );

                if (!new_block)
                {
                    return create_failed_result(ast_stream);
                }

                // Append the new block to the queue
                alinked_queue_ast_prepend(&blocks, new_block);
                break;
            }

            case TOKEN_OPEN_CURLY:
            {
                // Make sure we are in a block
                if (blocks.len == 0)
                {
                    create_error(
                        stream,
                        (ast_rule_t[]){AST_IMPORT, AST_FUNCTION, AST_MODULE},
                        3
                    );

                    return create_failed_result(ast_stream);
                }

                // Create a new block
                ast_t *block = ast_new(ast_stream.allocator, ast_stream.vec_allocator, TRUE);
                if (!block)
                {
                    create_error(
                        stream,
                        (ast_rule_t[]){AST_IMPORT, AST_FUNCTION, AST_MODULE},
                        3
                    );

                    return create_failed_result(ast_stream);
                }

                // Append the new block to the queue
                alinked_queue_ast_prepend(&blocks, block);
            }

            case TOKEN_CLOSE_CURLY:
            {
                // Make sure we are in a block
                if (blocks.len == 0)
                {
                    create_error(
                        stream,
                        (ast_rule_t[]){AST_IMPORT, AST_FUNCTION, AST_MODULE},
                        3
                    );

                    return create_failed_result(ast_stream);
                }

                // Delete the first element from the queue
                alinked_queue_ast_shift(&blocks);
                break;
            }

            default:
            {
                // Make sure we are in a block
                if (blocks.len == 0)
                {
                    create_error(
                        stream,
                        (ast_rule_t[]){AST_IMPORT, AST_FUNCTION, AST_MODULE},
                        3
                    );

                    return create_failed_result(ast_stream);
                }

                // Extract all tokens before the next semicolon
                const pair_extract_t extract = extract_tokens(
                    stream->tokens->data + stream->current,
                    stream->tokens->length - 1,
                    TOKEN_SEMICOLON,
                    TOKEN_SEMICOLON,
                    0,
                    FALSE
                );

                // Get the extracted range
                token_t **range = extract.first;
                const size_t len = extract.second;

                // Handle failure
                if (!range)
                {
                    create_error(
                        stream,
                        (ast_rule_t[]){AST_IMPORT, AST_FUNCTION, AST_MODULE},
                        3
                    );

                    return create_failed_result(ast_stream);
                }

                // Get the first block in the queue
                const ast_t *block = blocks.head->data;

                // Pass the range to the expression parser
                if (!
                    parse_expression(
                        block,
                        range,
                        0,
                        len,
                        ast_stream.allocator,
                        ast_stream.vec_allocator
                    )
                )
                {
                    return create_failed_result(ast_stream);
                }

                // Skip the extracted range
                stream->current += len;
                break;
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
