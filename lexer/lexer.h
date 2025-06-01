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
#include <fluent/str_conv/str_conv.h> // fluent_libc

// ============= INCLUDES =============
#include <ctype.h>
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

static bool push_token(
    vector_t *tokens,
    arena_allocator_t *allocator,
    string_builder_t *current,
    bool *in_string_ptr,
    bool *is_identifier_ptr,
    bool *is_number_ptr,
    bool *is_decimal_ptr,
    size_t *token_idx_ptr,
    const size_t line,
    const size_t column
) {
    // Ignore empty tokens
    if (current->idx == 0)
    {
        return TRUE; // No token to push
    }

    // Get the current token from the string builder
    char *curr = collect_string_builder_no_copy(current);

    // Dereference the pointers for easier access
    const bool is_identifier = *is_identifier_ptr;
    const token_type_t *type_ptr = hashmap_token_get(&fluent_token_map, curr);

    // Check if we have a valid token in the string builder
    if (
        !is_identifier &&
        type_ptr == NULL
    )
    {
        // Destroy the string builder
        destroy_string_builder(current);

        // Set the global error state
        global_error_state.code = LEXER_ERROR_UNKNOWN_TOKEN;
        global_error_state.column = column;
        global_error_state.line = line;
        return FALSE;
    }

    // Create a new token
    token_t *token = arena_malloc(allocator);
    if (!token)
    {
        destroy_string_builder(current);
        global_error_state.code = LEXER_ERROR_UNKNOWN;
        return FALSE; // Memory allocation failed
    }

    // Dereference the flags
    const bool is_number = *is_number_ptr;
    const bool is_decimal = *is_decimal_ptr;
    const bool in_string = *in_string_ptr;
    bool copy_value = FALSE;

    // Set the token properties
    token->line = line;
    token->column = column;

    // Add the token depending on the flags
    if (is_identifier)
    {
        // Initialize the token
        copy_value = TRUE; // We need to copy the value for identifiers
        token->type = TOKEN_IDENTIFIER;
    }
    else if (is_decimal)
    {
        // Initialize the token for a decimal number
        copy_value = TRUE; // We need to copy the value for decimal numbers
        token->type = TOKEN_DECIMAL_LITERAL;
    }
    else if (is_number)
    {
        // Initialize the token for a number
        copy_value = TRUE; // We need to copy the value for numbers
        token->type = TOKEN_NUM_LITERAL;
    }
    else if (in_string)
    {
        // Initialize the token for a string literal
        copy_value = TRUE; // We need to copy the value for string literals
        token->type = TOKEN_STRING_LITERAL;
    }
    else
    {
        // Initialize the token for a known type
        token->type = *type_ptr;
    }

    // Check if we need to copy the value
    if (copy_value)
    {
        char *value = collect_string_builder(current); // Copy the builder value
        // Move the value to a heap_guard to take advantage of
        // automatic freeing
        token->value = heap_alloc(
            0, // We are not allocating new memory, no need to specify size
            FALSE, // Value does not need to be concurrent
            FALSE, // Internal insertion does not need to be concurrent
            NULL, // No destructor
            value // The value to guard
        );
    }
    else
    {
        token->value = NULL; // No value for known types
    }

    // Push the token to the vector
    vec_push(tokens, token);

    // Reset the string builder for the next token
    reset_string_builder(current);

    // Reset the flags
    *in_string_ptr = FALSE;
    *is_identifier_ptr = TRUE;
    *is_number_ptr = FALSE;
    *is_decimal_ptr = FALSE;
    *token_idx_ptr = 0; // Reset the token index
    return TRUE;
}

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
    const heap_guard_t *tokens_guard = heap_alloc(sizeof(vector_t), FALSE, FALSE, NULL, NULL);
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
    bool in_unicode_escape = FALSE; // Whether we are inside a Unicode escape sequence
    bool in_surrogate_pair = FALSE; // Whether we are inside a surrogate pair
    bool unicode_escape_initialized = FALSE; // Whether the Unicode escape sequence has been initialized
    string_builder_t unicode_escape_sequence; // String builder for Unicode escape sequences
    size_t token_idx = 0; // Index for the current token

    // Iterate over the source
    for (size_t i = 0; source[i] != '\0'; i++)
    {
        // Get the current character
        const char c = source[i];

        // Check for comments
        if (!in_string && c == '/' && source[i + 1] == '/')
        {
            // Single-line comment
            in_comment = TRUE;
            i++; // Skip the next character
            continue;
        }

        // Handle block comments
        if (!in_string && c == '/' && source[i + 1] == '*')
        {
            // Block comment start
            in_block_comment = TRUE;
            i++; // Skip the next character
            continue;
        }

        // Handle block comment end
        if (!in_string && c == '*' && source[i + 1] == '/')
        {
            // Block comment end
            in_block_comment = FALSE;
            i++; // Skip the next character
            continue;
        }

        // Ignore comments
        if (in_comment || in_block_comment)
        {
            // Check for newlines to exit comment state
            if (c == '\n')
            {
                in_comment = FALSE; // Exit comment state on newline
                line++;
                column = 1; // Reset column on new line
            }
            continue;
        }

        // Check for newlines
        if (c == '\n')
        {
            // Check if we are in a string
            if (in_string)
            {
                // Destroy the string builder
                destroy_string_builder(&current);

                global_error_state.code = LEXER_ERROR_UNTERMINATED_STRING;
                global_error_state.column = column;
                global_error_state.line = line;
                return pair_lex_result_new(stream, &global_error_state);
            }

            // Push the current token if it exists
            if (!push_token(
                tokens,
                allocator,
                &current,
                &in_string,
                &is_identifier,
                &is_number,
                &is_decimal,
                &token_idx,
                line,
                column
            ))
            {
                // If pushing the token failed, return the error state
                return pair_lex_result_new(stream, &global_error_state);
            }

            line++;
            column = 1; // Reset column on new line

            continue;
        }

        // Check for whitespace
        if (c == ' ' && !in_string)
        {
            // Push the current token if it exists
            if (!push_token(
                tokens,
                allocator,
                &current,
                &in_string,
                &is_identifier,
                &is_number,
                &is_decimal,
                &token_idx,
                line,
                column
            ))
            {
                // If pushing the token failed, return the error state
                return pair_lex_result_new(stream, &global_error_state);
            }

            continue;
        }

        // Check if we have a punctuation character
        if (!in_string && hashmap_btoken_get(&fluent_punctuation_map, c))
        {
            // Push the current token if it exists
            if (!push_token(
                tokens,
                allocator,
                &current,
                &in_string,
                &is_identifier,
                &is_number,
                &is_decimal,
                &token_idx,
                line,
                column
            ))
            {
                // If pushing the token failed, return the error state
                return pair_lex_result_new(stream, &global_error_state);
            }

            // Increment the column for the punctuation character
            column++;
            is_identifier = FALSE; // Reset identifier state

            // Write the punctuation character to the string builder
            write_char_string_builder(&current, c);

            // Push the punctuation token
            // Push the current token if it exists
            if (!push_token(
                tokens,
                allocator,
                &current,
                &in_string,
                &is_identifier,
                &is_number,
                &is_decimal,
                &token_idx,
                line,
                column
            ))
            {
                // If pushing the token failed, return the error state
                return pair_lex_result_new(stream, &global_error_state);
            }

            // Increment the column for the next character
            column++;
            continue;
        }

        // Handle string escapes
        if (in_string && c == '\\' && !in_str_escape)
        {
            // If we are in a string and encounter an escape character
            in_str_escape = TRUE; // Enter escape state
            column++; // Increment column for the escape character
            write_char_string_builder(&unicode_escape_sequence, c); // Write the escape character to the string builder
            continue;
        }

        // Handle escape sequences
        if (in_str_escape)
        {
            // Initialize Unicode escape sequence builder
            if (!unicode_escape_initialized)
            {
                // Initialize the Unicode escape sequence string builder
                init_string_builder(&unicode_escape_sequence, 8, 1.5);
                unicode_escape_initialized = TRUE; // Set the flag to true
                write_char_string_builder(&unicode_escape_sequence, '\\'); // Write the escape character
            }

            // Handle Unicode escape sequences
            if (c == 'u' && !in_unicode_escape)
            {
                in_unicode_escape = TRUE; // Enter Unicode escape state
                column++; // Increment column for the 'u'
                write_char_string_builder(&unicode_escape_sequence, c); // Write the 'u' character to the string builder
                continue;
            }

            // Handle Unicode surrogate pairs
            if (c == 'U' && !in_surrogate_pair)
            {
                in_surrogate_pair = TRUE; // Enter surrogate pair state
                column++; // Increment column for the 'U'
                continue;
            }

            // Check if we have ended a Unicode escape sequence
            if (
                (
                    !in_unicode_escape &&
                    !in_surrogate_pair
                ) ||
                (
                    (in_unicode_escape && unicode_escape_sequence.idx + 1 == 4) ||
                    in_surrogate_pair && unicode_escape_sequence.idx + 1 == 8
                )
            )
            {
                // Write the character to the string builder
                write_char_string_builder(&unicode_escape_sequence, c);

                // Convert the Unicode escape sequence to a character
                const char *unicode_char = convert_escapes_to_utf8_sb(collect_string_builder_no_copy(&unicode_escape_sequence));

                // Reset the Unicode escape sequence string builder
                reset_string_builder(&unicode_escape_sequence);

                // Write the Unicode character to the current token
                write_string_builder(&current, unicode_char);

                // Reset flags
                in_unicode_escape = FALSE; // Reset Unicode escape state
                in_surrogate_pair = FALSE; // Reset surrogate pair state
                in_str_escape = FALSE; // Reset escape state
                continue;
            }

            // Write the escape character to the string builder
            if (in_unicode_escape || in_surrogate_pair)
            {
                // If we are in a Unicode escape or surrogate pair, write the character
                write_char_string_builder(&unicode_escape_sequence, c);
                column++; // Increment column for the escape character
            }

            continue;
        }

        // Handle string literals
        if (c == '"')
        {
            // Push the current token if it exists
            if (!push_token(
                tokens,
                allocator,
                &current,
                &in_string,
                &is_identifier,
                &is_number,
                &is_decimal,
                &token_idx,
                line,
                column
            ))
            {
                // If pushing the token failed, return the error state
                return pair_lex_result_new(stream, &global_error_state);
            }

            // If we are not already in a string, start a new string
            if (!in_string)
            {
                in_string = TRUE; // Enter string state
                is_identifier = FALSE; // Reset identifier state
                is_number = FALSE; // Reset number state
                is_decimal = FALSE; // Reset decimal state
            }
            else
            {
                // If we are already in a string, end the string
                in_string = FALSE; // Exit string state
            }

            column++; // Increment column for the string character
            continue;
        }

        // Recognize identifiers without regex
        if (token_idx == 0)
        {
            // Check if the character is a valid identifier start
            if (isalpha(c) || c == '_')
            {
                is_identifier = TRUE; // Start an identifier
                token_idx++; // Increment token index
            }
            else if (isdigit(c))
            {
                is_identifier = FALSE; // Reset identifier state
                is_number = TRUE; // Start a number
                token_idx++; // Increment token index
            }

            token_idx = 1; // Prevent processing in the next iteration
        }

        // Handle decimal literals
        if (is_number && c == '.')
        {
            // If we are in a number and encounter a dot, it might be a decimal
            if (is_decimal)
            {
                // If we already have a decimal, it's an error
                destroy_string_builder(&current);
                global_error_state.code = LEXER_ERROR_UNKNOWN_TOKEN;
                global_error_state.column = column;
                global_error_state.line = line;
                return pair_lex_result_new(stream, &global_error_state);
            }

            is_decimal = TRUE; // Start a decimal number
            token_idx++; // Increment token index
            column++; // Increment column for the dot character
            write_char_string_builder(&current, c);
            continue;
        }

        // Handle chainable tokens
        if (source[i + 1] == '=' && hashmap_btoken_get(&fluent_chainable_map, c))
        {
            // Write the current token if it exists
            if (!push_token(
                tokens,
                allocator,
                &current,
                &in_string,
                &is_identifier,
                &is_number,
                &is_decimal,
                &token_idx,
                line,
                column
            ))
            {
                // If pushing the token failed, return the error state
                return pair_lex_result_new(stream, &global_error_state);
            }

            // Write the chainable token to the string builder
            write_char_string_builder(&current, c);
            write_char_string_builder(&current, '=');

            // Write the chainable token to the tokens vector
            if (!push_token(
                tokens,
                allocator,
                &current,
                &in_string,
                &is_identifier,
                &is_number,
                &is_decimal,
                &token_idx,
                line,
                column
            ))
            {
                // If pushing the token failed, return the error state
                return pair_lex_result_new(stream, &global_error_state);
            }

            // Increment the column for both characters
            column += 2;
            i++; // Skip the next character
            continue;
        }

        // Write the current character to the string builder
        write_char_string_builder(&current, c);

        column++; // Increment column for other characters
    }

    // Destroy the string builders
    destroy_string_builder(&current);
    if (unicode_escape_initialized)
    {
        destroy_string_builder(&unicode_escape_sequence);
    }

    return pair_lex_result_new(stream, NULL);
}

#endif //FLUENT_LEXER_H
