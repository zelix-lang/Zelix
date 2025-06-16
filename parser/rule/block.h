//
// Created by rodrigo on 6/16/25.
//

#ifndef FLUENT_PARSER_RULE_BLOCK_H
#define FLUENT_PARSER_RULE_BLOCK_H

// ============= FLUENT LIB C =============
#include <fluent/arena/arena.h> // fluent_libc
#include <fluent/std_bool/std_bool.h> // fluent_libc

// ============= INCLUDES =============
#include "../../ast/ast.h"
#include "../../lexer/stream.h"

static inline bool parse_block(
    const ast_t *const root,
    token_stream_t *const stream,
    arena_allocator_t *const arena,
    arena_allocator_t *const vec_arena
)
{
    token_t *token = token_stream_next(stream);
    printf("%d\n", token->type);
    return TRUE;
}

#endif //FLUENT_PARSER_RULE_BLOCK_H
