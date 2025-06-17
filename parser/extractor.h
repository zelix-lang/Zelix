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
// Created by rodrigo on 6/7/25.
//

#ifndef FLUENT_TOKEN_EXTRACTOR_H
#define FLUENT_TOKEN_EXTRACTOR_H

// ============= FLUENT LIB C =============
#include <fluent/pair/pair.h> // fluent_libc
#include <fluent/std_bool/std_bool.h> // fluent_libc

DEFINE_PAIR_T(token_t **, size_t, extract);

/**
 * Extracts a sequence of tokens from a token stream, delimited by specified start and end token types.
 *
 * This function searches the token stream starting at the given index (`start`) for a matching pair of
 * delimiters (`delim` and `end_delim`). It supports nested delimiters by maintaining a counter.
 * If a matching end delimiter is found at the correct nesting level, a pair containing a pointer to the
 * start of the extracted tokens and the number of tokens extracted is returned.
 * If delimiters are mismatched or not found, a pair with NULL and 0 is returned.
 *
 * @param tokens      Pointer to the token stream (array of token_t pointers).
 * @param length      The total number of tokens in the stream.
 * @param delim       The token type that marks the start of a delimited section.
 * @param end_delim   The token type that marks the end of a delimited section.
 * @param start       The index in the token stream to start extraction.
 * @param allow_nested Whether to allow nested delimiters (if TRUE, nested delimiters are counted).
 * @return            A pair containing a pointer to the extracted tokens and the count, or (NULL, 0) on error.
 */
static inline pair_extract_t extract_tokens(
    token_t **tokens,
    const size_t length,
    const token_type_t delim,
    const token_type_t end_delim,
    const size_t start,
    const bool allow_nested
)
{
    // Make sure the stream has enough space to extract tokens
    if (length < start)
    {
        return pair_extract_new(NULL, 0); // Not enough tokens to extract
    }

    // Get the new buffer without copying
    token_t **new_buffer = tokens + start;

    // Define a nested counter
    size_t counter = 0;

    // Counter to know where the counting stopped at
    size_t end = 0;

    // Boolean to know if we have met an end delimiter
    bool has_met_delim = FALSE;

    // Iterate over the buffer
    for (size_t i = start; i < length; i++)
    {
        // Increment the end index
        end++;

        const token_t *token = tokens[i];
        if (token->type == end_delim)
        {
            // Bail out if nesting is not allowed
            if (!allow_nested && counter > 0)
            {
                break;
            }

            // Check if we have a 0 counter
            if (counter == 0)
            {
                // Mismatched delimiters, return NULL
                return pair_extract_new(NULL, 0);
            }

            counter--;

            // Check if we have reached a 0 counter
            if (counter == 0)
            {
                has_met_delim = TRUE; // We have met the end delimiter
                break;
            }
        }
        else if (token->type == delim)
        {
            counter++;
        }
    }

    // Check if we have met the end delimiter
    if (!has_met_delim)
    {
        // We have not met the end delimiter, return NULL
        return pair_extract_new(NULL, 0);
    }

    // Return the pair
    return pair_extract_new(new_buffer, end);
}

#endif //FLUENT_TOKEN_EXTRACTOR_H
