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
// Created by rodrigo on 6/16/25.
//

#ifndef FLUENT_PARSER_RULE_EXPRESSION_H
#define FLUENT_PARSER_RULE_EXPRESSION_H

// ============= FLUENT LIB C =============
#include <fluent/arena/arena.h> // fluent_libc
#include <fluent/std_bool/std_bool.h> // fluent_libc

// ============= INCLUDES =============
#include "../../ast/ast.h"
#include "../../token/token.h"

static inline bool parse_expression(
    ast_t *const root,
    token_t **body,
    const size_t start,
    const size_t body_len,
    arena_allocator_t *const arena,
    arena_allocator_t *const vec_arena
)
{
    return TRUE;
}

#endif //FLUENT_PARSER_RULE_EXPRESSION_H
