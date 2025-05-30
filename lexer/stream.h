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
// Created by rodrigo on 5/30/25.
//

#ifndef FLUENT_TOKEN_STREAM_H
#define FLUENT_TOKEN_STREAM_H

// ============= FLUENT LIB C =============
#include <fluent/vector/vector.h> // fluent_libc

#include "../token/token.h"

/**
 * @brief Represents a stream of tokens for the lexer.
 *
 * This structure holds a vector of tokens and an index indicating
 * the current position in the token stream.
 */
typedef struct
{
    vector_t *tokens; ///< Vector of tokens
    size_t current;  ///< Current index in the token stream
    arena_allocator_t *allocator; ///< Arena allocator for token memory management
} token_stream_t;

/**
 * @brief Returns the n-th token in the token stream.
 *
 * @param stream Pointer to the token stream.
 * @param n Index of the token to retrieve.
 * @return Pointer to the n-th token, or NULL if out of bounds.
 */
static inline token_t *token_stream_nth(const token_stream_t *const stream, const size_t n)
{
    if (n < stream->tokens->length)
    {
        return (token_t *)vec_get(stream->tokens, n);
    }

    return NULL;
}

/**
 * @brief Peeks at the next token in the stream without advancing the position.
 *
 * @param stream Pointer to the token stream.
 * @return Pointer to the next token, or NULL if at the end.
 */
static inline token_t *token_stream_peek(const token_stream_t *const stream)
{
    return token_stream_nth(stream, stream->current + 1);
}

/**
 * @brief Advances the stream and returns the next token.
 *
 * @param stream Pointer to the token stream.
 * @return Pointer to the next token, or NULL if at the end.
 */
static inline token_t *token_stream_next(token_stream_t *const stream)
{
    if (stream->current + 1 < stream->tokens->length)
    {
        stream->current++;
        return token_stream_nth(stream, stream->current);
    }

    return NULL;
}

#endif //FLUENT_TOKEN_STREAM_H
