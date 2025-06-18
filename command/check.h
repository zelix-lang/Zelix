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
// Created by rodrigo on 5/29/25.
//

#ifndef FLUENT_COMMAND_CHECK_H
#define FLUENT_COMMAND_CHECK_H

// ============= INCLUDES =============
#include <file/file_reader.h>
#include <message/lexer_error_converter.h>
#include <message/generator.h>
#include <lexer/lexer.h>
#include <logger/logger.h>
#include <parser/parser.h>

static inline void check_command(const char *const path)
{
    // Normalize the path
    char *normalized_path = get_real_path(path);

    // Read the file contents
    char *source = read_file(normalized_path, FALSE);
    if (!source)
    {
        fprintf(stderr, "Error: Could not read file '%s'.\n", path);
        return;
    }

    // Get the file name from the path
    char *file_name = get_file_name(path);

    // Lex the source code
    const pair_lex_result_t result = lexer_tokenize(source, file_name);

    token_stream_t stream = result.first; // Get the token stream from the result
    const lexer_error_t *error = result.second; // Get the error if any

    // Check for lexer errors
    if (error)
    {
        // Emit failed state
        timer_failed();

        // Log the error
        log_error(lexer_error_to_string(error));
        log_info("Full details:");

        // Build the error message
        char *msg = build_error_message(
            source,
            normalized_path,
            ANSI_BOLD_BRIGHT_RED,
            error->line,
            error->column,
            error->col_start
        );
        printf("%s", msg);

        // Free the error message
        free(msg);
    }

    // Emit done state
    timer_done();

    // Parse the code
    const pair_parser_result_t parser_result = parser_parse(&stream, file_name);
    const ast_stream_t ast_stream = parser_result.first; // Get the AST stream from the result
    const ast_error_t *parser_error = parser_result.second; // Get the parser error if any

    // Handle parser errors
    if (parser_error)
    {
        // Emit failed state
        timer_failed();

        // Log the error
        log_error("Invalid syntax in the source code");
        log_info("Full details:");

        // Build the error message
        char *msg = build_error_message(
            source,
            normalized_path,
            ANSI_BOLD_BRIGHT_RED,
            parser_error->line,
            parser_error->column,
            parser_error->col_start
        );
        printf("%s", msg);

        // Free the error message
        free(msg);
    }
    else
    {
        // Emit success state
        timer_done();
    }

    // Free the memory used by the parser
    if (ast_stream.allocator)
    {
        // Destroy the AST stream
        destroy_arena(ast_stream.allocator);
        destroy_arena(ast_stream.vec_allocator);
    }

    // Free the resources used by the lexer
    if (stream.allocator)
    {
        // Destroy the vector of tokens
        vec_token_destroy(stream.tokens, NULL);
        destroy_arena(stream.allocator);
    }

    // Free the allocated memory
    free(source);
    free(file_name);
    free(normalized_path);
}

#endif //FLUENT_COMMAND_CHECK_H
