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
// Created by rodrigo on 6/6/25.
//

#ifndef FLUENT_AST_STREAM_H
#define FLUENT_AST_STREAM_H

// ============= FLUENT LIB C =============
#include <fluent/arena/arena.h>

// ============= INCLUDES =============
#include <ast/ast.h>

typedef struct
{
    ast_t *ast; /**< Pointer to the root AST node. */
    arena_allocator_t *allocator; /**< Arena allocator for memory management. */
    arena_allocator_t *vec_allocator; /**< Arena allocator for vector memory management. */
} ast_stream_t;

#endif //FLUENT_AST_STREAM_H
