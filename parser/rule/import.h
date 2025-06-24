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

#ifndef FLUENT_AST_RULE_IMPORT_H
#define FLUENT_AST_RULE_IMPORT_H

// ============== FLUENT LIB C =============
#include <fluent/std_bool/std_bool.h> // fluent_lib
#include <fluent/arena/arena.h> // fluent_lib

// ============= INCLUDES =============
#include <parser/error.h>
#include "ast/ast.h"
#include <lexer/stream.h>

static inline bool parse_import(
    const ast_t *const root,
    token_stream_t *const stream,
    arena_allocator_t *const arena,
    arena_allocator_t *const vec_arena
)
{
    // Get the next token
    const token_t *token = token_stream_next(stream);

    // Get the next token to validate semantics
    const token_t *semicolon = token_stream_next(stream);

    // Make sure we didn't get NULL and the token is a string literal
    if (!token || token->type != TOKEN_STRING_LITERAL || !semicolon || semicolon->type != TOKEN_SEMICOLON)
    {
        // Create an error for unexpected end of stream
        create_error(
            stream,
            (ast_rule_t[]){AST_STRING_LITERAL},
            1
        );
        return FALSE;
    }

    // Allocate a new AST node for the import rule
    ast_t *import_node = ast_new(arena, vec_arena, TRUE);
    if (!import_node)
    {
        // Failed to allocate memory for the import node
        return FALSE;
    }

    // Set import rule
    import_node->rule = AST_IMPORT; // Set the rule type to AST_IMPORT

    // Allocate a new AST for the string literal
    ast_t *string_literal_node = ast_new(arena, vec_arena, FALSE);
    if (!string_literal_node)
    {
        // Failed to allocate memory for the string literal node
        return FALSE;
    }

    // Set the string literal rule
    string_literal_node->rule = AST_STRING_LITERAL;
    string_literal_node->value = token->value;

    // Append the string literal node to the import node's children
    vec_ast_push(import_node->children, string_literal_node);

    // Append the import node to the root's children
    vec_ast_push(root->children, import_node);

    return TRUE;
}

#endif //FLUENT_AST_RULE_IMPORT_H
