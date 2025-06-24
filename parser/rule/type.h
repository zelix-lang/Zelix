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
// Created by rodrigo on 6/9/25.
//

#ifndef FLUENT_PARSER_TYPE_H
#define FLUENT_PARSER_TYPE_H

// ============= FLUENT LIB C =============
#include <fluent/std_bool/std_bool.h> // fluent_libc
#include <fluent/alinked_queue/alinked_queue.h> // fluent_libc
#include <fluent/arena/arena.h> // fluent_libc

// ============= INCLUDES =============
#include "ast/ast.h"
#include <type/type.h>
#include <token/token.h>

DEFINE_ALINKED_NODE(ast_t *, ast);

static ast_rule_t parse_base_type(
    const token_t *token,
    const ast_t *root,
    arena_allocator_t *const arena,
    arena_allocator_t *const vec_arena
)
{
    token_type_t primitive_type = TOKEN_UNKNOWN;

    // Check if we have a primitive type
    switch (token->type)
    {
        case TOKEN_NUM:
        {
            // Set the rule to AST_NUMBER
            primitive_type = AST_NUMBER;
            break;
        }

        case TOKEN_BOOL:
        {
            // Set the rule to AST_BOOL
            primitive_type = AST_BOOL;
            break;
        }

        case TOKEN_STRING:
        {
            // Set the rule to AST_STRING
            primitive_type = AST_STRING;
            break;
        }

        case TOKEN_DEC:
        {
            // Set the rule to AST_DECIMAL
            primitive_type = AST_DECIMAL;
            break;
        }

        case TOKEN_NOTHING:
        {
            // Set the rule to AST_NOTHING
            primitive_type = AST_NOTHING;
            break;
        }

        case TOKEN_IDENTIFIER:
        {
            // Set the rule to AST_CUSTOM_TYPE
            primitive_type = AST_IDENTIFIER;
            break;
        }

        default:
        {
            break;
        }
    }

    // Return primitive types earlier
    if (primitive_type != TOKEN_UNKNOWN)
    {
        // Create a new AST node for the type
        ast_t *new_node = ast_new(arena, vec_arena, FALSE);
        if (!new_node)
        {
            // Failed to allocate the node
            return AST_PROGRAM_RULE;
        }

        // Set the rule for the node
        new_node->rule = primitive_type;

        // Append the value to the node
        vec_ast_push(root->children, new_node);

        // Check for identifiers
        if (primitive_type == AST_IDENTIFIER)
        {
            new_node->value = token->value; // Set the value to the identifier
        }

        // Return the parsed rule
        return primitive_type;
    }

    // Invalid type
    return AST_PROGRAM_RULE;
}

static inline ast_t * parse_type(
    token_stream_t *stream,
    arena_allocator_t *const arena,
    arena_allocator_t *const vec_arena
)
{
    // Allocate the root node
    ast_t *root = ast_new(arena, vec_arena, TRUE);
    if (!root)
    {
        // Failed to allocate the root node
        return NULL;
    }

    // Get the first token
    const token_t *token = token_stream_nth(stream, stream->current);

    // Parse the base type
    const token_type_t primitive_type = parse_base_type(
        token,
        root,
        arena,
        vec_arena
    );

    // Ignore non-generic types
    if (primitive_type == AST_PROGRAM_RULE)
    {
        // Invalid type definition
        return NULL;
    }

    if (primitive_type != AST_IDENTIFIER)
    {
        return root; // Return the root node if it's a primitive type
    }

    // Peek to see if we have generic types
    const token_t *peek = token_stream_peek(stream);
    if (!peek || peek->type != TOKEN_LESS_THAN)
    {
        // Return the root node if we don't have a generic type
        return root;
    }

    // Skip the less than token and position at the
    // token after it
    token_stream_next(stream);
    token = token_stream_next(stream);

    // Define flags to track nested types
    size_t nest_level = 1;
    bool expecting_comma = FALSE;
    bool generics_allowed = FALSE;
    alinked_queue_ast_t queue;
    alinked_queue_ast_init(&queue, 20);
    alinked_queue_ast_append(&queue, root);

    // Iterate over the tokens to parse the type
    while (token != NULL)
    {
        // Check for end of the type definitions
        if (token->type == TOKEN_GREATER_THAN)
        {
            // Get the first element from the queue
            const ast_t *current = alinked_queue_ast_shift(&queue);

            // Make sure the AST has children and nesting level is valid
            if (current->children->length == 0)
            {
                return NULL; // Invalid type definition
            }

            // Decrease the nest level
            nest_level--;

            // Bail out if we are at the root level
            if (nest_level == 0)
            {
                // Return the root node
                return root;
            }

            continue;
        }

        // Check for the start if a type definition
        if (token->type == TOKEN_LESS_THAN)
        {
            // Make sure that generics are allowed
            if (!generics_allowed)
            {
                return NULL;
            }

            // Increment the nesting level
            nest_level++;

            // Create a new type AST
            ast_t *new_type = ast_new(arena, vec_arena, TRUE);

            // Add the nested type to the queue
            alinked_queue_ast_prepend(&queue, new_type);

            // Set the expecting comma flag to TRUE for further processing
            expecting_comma = TRUE;

            // Continue
            continue;
        }

        // Handle commas
        if (token->type == TOKEN_COMMA)
        {
            // Check if we are expecting a comma
            if (!expecting_comma)
            {
                return NULL; // Invalid type definition
            }

            // Flip the flag
            expecting_comma = FALSE;
            continue;
        }

        // Get the first element of the queue
        const ast_t *current = queue.head->data;

        // Parse the base type
        const ast_rule_t type_rule = parse_base_type(
            token,
            current,
            arena,
            vec_arena
        );

        // Set the generics allowed flag
        generics_allowed = type_rule == AST_IDENTIFIER;
        expecting_comma = TRUE;

        // Move to the next token
        token = token_stream_next(stream);
    }

    // Destroy the queue when done
    alinked_queue_ast_destroy(&queue);

    // Make sure we don't end up with a nested level
    if (nest_level != 0)
    {
        return NULL; // Invalid type definition
    }

    return root;
}

#endif //FLUENT_PARSER_TYPE_H
