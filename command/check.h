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

#include "../file/file_reader.h"
#include "../message/lexer_error_converter.h"
#include "../message/generator.h"
#include "../lexer/lexer.h"
#include "../logger/logger.h"

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

    const token_stream_t stream = result.first; // Get the token stream from the result
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
            error->column
        );
        printf("%s", msg);

        // Free the error message
        free(msg);
    }

    // Free the resources used by the lexer
    if (stream.allocator)
    {
        // Destroy the vector of tokens
        vec_destroy(stream.tokens, NULL);
        destroy_arena(stream.allocator);
    }

    // Free the allocated memory
    free(source);
    free(file_name);
    free(normalized_path);
}

#endif //FLUENT_COMMAND_CHECK_H
