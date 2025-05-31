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

#ifndef FLUENT_LEXER_H
#define FLUENT_LEXER_H

// ============= FLUENT LIB C =============
#include <fluent/vector/vector.h> // fluent_libc
#include <fluent/string_builder/string_builder.h> // fluent_libc
#include <fluent/std_bool/std_bool.h> // fluent_libc
#include <fluent/pair/pair.h> // fluent_libc

// ============= INCLUDES =============
#include <fluent/arena/arena.h>
#include <fluent/heap_guard/heap_guard.h>

#include "error.h"
#include "stream.h"
#include "../token/token_map.h"

// ============= MACROS =============
#ifndef FLUENT_PAIR_LEXER
#   define FLUENT_PAIR_LEXER 1
    DEFINE_PAIR_T(token_stream_t, lexer_error_t *, lex_result);
#endif // FLUENT_PAIR_LEXER

// ============= INCLUDES =============
#include "../token/token.h"

lexer_error_t global_error_state;

static inline pair_lex_result_t lexer_tokenize(
    const char *source,
    const char *path
)
{
    // Reset the global error state
    global_error_state.code = LEXER_ERROR_UNKNOWN;

    // Initialize an arena allocator for token allocation
    arena_allocator_t *allocator = arena_new(25, sizeof(token_t));

    // Initialize a vector to hold tokens
    const heap_guard_t *tokens_guard = heap_alloc(sizeof(vector_t), FALSE, FALSE, NULL);
    vec_init(tokens_guard->ptr, 256, sizeof(token_t), 1.5);
    vector_t *tokens = (vector_t *)tokens_guard->ptr;

    // Initialize the token stream
    token_stream_t stream;
    stream.tokens = (vector_t *)tokens_guard->ptr;
    stream.current = 0;
    stream.allocator = allocator;

    // Use a string builder to build the current token
    string_builder_t current;
    init_string_builder(&current, 64, 1.5);

    // Track the current position in the source code
    size_t line = 1;
    size_t column = 1;

    // Track lexing state
    bool in_string = FALSE; // Whether we are inside a string literal
    bool in_comment = FALSE; // Whether we are inside a comment
    bool in_block_comment = FALSE; // Whether we are inside a block comment
    bool in_str_escape = FALSE; // Whether we are inside a string escape sequence
    bool is_identifier = FALSE; // Whether the current token is an identifier
    bool is_number = FALSE; // Whether the current token is a number
    bool is_decimal = FALSE; // Whether the current token is a decimal number
    size_t token_idx = 0; // Index for the current token

    // Iterate over the source
    for (size_t i = 0; source[i] != '\0'; i++)
    {
        // Get the current character
        const char c = source[i];

        // Check for newlines
        if (c == '\n')
        {
            // Check if we are in a string
            if (in_string)
            {
                global_error_state.code = LEXER_ERROR_UNTERMINATED_STRING;
                global_error_state.column = column;
                global_error_state.line = line;
                return pair_lex_result_new(stream, &global_error_state);
            }

            line++;
            column = 1; // Reset column on new line
            in_comment = FALSE; // Exit comment state on newline

            continue;
        }

        // Check for whitespace
        if (c == ' ')
        {
            continue;
        }

        // Check for comments
        if (c == '/' && source[i + 1] == '/')
        {
            // Single-line comment
            in_comment = TRUE;
            i++; // Skip the next character
            continue;
        }

        // Handle block comments
        if (c == '/' && source[i + 1] == '*')
        {
            // Block comment start
            in_block_comment = TRUE;
            i++; // Skip the next character
            continue;
        }

        // Handle block comment end
        if (c == '*' && source[i + 1] == '/')
        {
            // Block comment end
            in_block_comment = FALSE;
            i++; // Skip the next character
            continue;
        }

        // Check if we have a punctuation character
        if (hashmap_btoken_get(&fluent_punctuation_map, c))
        {
            //
            column++;
            continue;
        }

        // Write the current character to the string builder
        write_char_string_builder(&current, c);

        column++; // Increment column for other characters
    }

    // Destroy the string builder
    destroy_string_builder(&current);

    return pair_lex_result_new(stream, NULL);
}

#endif //FLUENT_LEXER_H
