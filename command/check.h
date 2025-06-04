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
#include "../lexer/lexer.h"

static inline void check_command(const char *const path)
{
    // Read the file contents
    char *source = read_file(path);
    if (!source)
    {
        fprintf(stderr, "Error: Could not read file '%s'.\n", path);
        return;
    }

    // Lex the source code
    const pair_lex_result_t result = lexer_tokenize(source, path);

    const token_stream_t stream = result.first; // Get the token stream from the result
    const lexer_error_t *error = result.second; // Get the error if any

    // Check for lexer errors
    if (error)
    {
        fprintf(stderr, "Lexer error at line %zu, column %zu: %d\n",
                error->line, error->column, error->code);
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
}

#endif //FLUENT_COMMAND_CHECK_H
