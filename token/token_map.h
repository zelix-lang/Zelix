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

#ifndef FLUENT_TOKEN_MAP_H
#define FLUENT_TOKEN_MAP_H

// ============= FLUENT LIB C =============
#include <fluent/hashmap/hashmap.h> // fluent_libc
#include <fluent/heap_guard/heap_guard.h> // fluent_libc
#include <fluent/arena/arena.h> // fluent_libc
#include <fluent/str_comp/str_comp.h> // fluent_libc
#include <fluent/std_bool/std_bool.h> // fluent_libc

// ============= INCLUDES =============
#include "token.h"

#ifndef FLUENT_LIBC_CLI_HASHMAP_TOKEN_VALUE
    DEFINE_HASHMAP(char *, token_type_t, token);
    DEFINE_HASHMAP(char, bool, btoken);
#   define FLUENT_LIBC_CLI_HASHMAP_TOKEN_VALUE 1
#endif

// Holds the mapping of token names to their corresponding token types
hashmap_token_t fluent_token_map;
hashmap_btoken_t fluent_punctuation_map;
hashmap_btoken_t fluent_chainable_map;

/**
 * Compares two characters for equality.
 *
 * \param a The first character.
 * \param b The second character.
 * \return true if the characters are equal, false otherwise.
 */
static inline bool char_cmp(const char a, const char b)
{
    return a == b;
}

static inline hashmap_btoken_t *get_punctuation_map()
{
    // Check if the punctuation map is already initialized
    if (fluent_punctuation_map.hash_fn)
    {
        // Return the existing punctuation map
        return &fluent_punctuation_map;
    }

    // Initialize the heap guard for the punctuation map
    hashmap_btoken_init(
        &fluent_punctuation_map,
        55,
        1.5,
        NULL,
        (hash_btoken_function_t)hash_char_key,
        (hash_btoken_cmp_t)char_cmp
    );

    hashmap_btoken_insert(&fluent_punctuation_map, ';', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, ',', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '(', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, ')', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '{', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '}', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, ':', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '+', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '-', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '>', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '<', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '%', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '*', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '.', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '=', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '!', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '/', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '&', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '|', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '^', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, '[', TRUE);
    hashmap_btoken_insert(&fluent_punctuation_map, ']', TRUE);

    return &fluent_punctuation_map;
}

static inline hashmap_token_t *get_token_map()
{
    // Check if the token map is already initialized
    if (fluent_token_map.hash_fn)
    {
        // Return the existing token map
        return &fluent_token_map;
    }

    // Initialize the heap guard for the token map
    hashmap_token_init(
        &fluent_token_map,
        55,
        1.5,
        NULL,
        (hash_token_function_t)hash_str_key,
        (hash_token_cmp_t)str_comp
    );

    hashmap_token_insert(&fluent_token_map, "fun", TOKEN_FUNCTION);
    hashmap_token_insert(&fluent_token_map, "let", TOKEN_LET);
    hashmap_token_insert(&fluent_token_map, "const", TOKEN_CONST);
    hashmap_token_insert(&fluent_token_map, "while", TOKEN_WHILE);
    hashmap_token_insert(&fluent_token_map, "for", TOKEN_FOR);
    hashmap_token_insert(&fluent_token_map, "break", TOKEN_BREAK);
    hashmap_token_insert(&fluent_token_map, "continue", TOKEN_CONTINUE);
    hashmap_token_insert(&fluent_token_map, "if", TOKEN_IF);
    hashmap_token_insert(&fluent_token_map, "else", TOKEN_ELSE);
    hashmap_token_insert(&fluent_token_map, "elseif", TOKEN_ELSE_IF);
    hashmap_token_insert(&fluent_token_map, "return", TOKEN_RETURN);
    hashmap_token_insert(&fluent_token_map, "mod", TOKEN_MOD);
    hashmap_token_insert(&fluent_token_map, "new", TOKEN_NEW);
    hashmap_token_insert(&fluent_token_map, "in", TOKEN_IN);
    hashmap_token_insert(&fluent_token_map, "to", TOKEN_TO);

    hashmap_token_insert(&fluent_token_map, "=", TOKEN_ASSIGN);
    hashmap_token_insert(&fluent_token_map, "+", TOKEN_PLUS);
    hashmap_token_insert(&fluent_token_map, "-", TOKEN_MINUS);
    hashmap_token_insert(&fluent_token_map, "*", TOKEN_ASTERISK);
    hashmap_token_insert(&fluent_token_map, "/", TOKEN_SLASH);
    hashmap_token_insert(&fluent_token_map, "<", TOKEN_LESS_THAN);
    hashmap_token_insert(&fluent_token_map, ">", TOKEN_GREATER_THAN);
    hashmap_token_insert(&fluent_token_map, "==", TOKEN_EQUAL);
    hashmap_token_insert(&fluent_token_map, "!=", TOKEN_NOT_EQUAL);
    hashmap_token_insert(&fluent_token_map, ">=", TOKEN_GREATER_THAN_OR_EQUAL);
    hashmap_token_insert(&fluent_token_map, "<=", TOKEN_LESS_THAN_OR_EQUAL);
    hashmap_token_insert(&fluent_token_map, "&", TOKEN_AMPERSAND);
    hashmap_token_insert(&fluent_token_map, "&&", TOKEN_AND);
    hashmap_token_insert(&fluent_token_map, "||", TOKEN_OR);
    hashmap_token_insert(&fluent_token_map, "|", TOKEN_BAR);
    hashmap_token_insert(&fluent_token_map, "!", TOKEN_NOT);
    hashmap_token_insert(&fluent_token_map, ",", TOKEN_COMMA);
    hashmap_token_insert(&fluent_token_map, ";", TOKEN_SEMICOLON);
    hashmap_token_insert(&fluent_token_map, "(", TOKEN_OPEN_PAREN);
    hashmap_token_insert(&fluent_token_map, ")", TOKEN_CLOSE_PAREN);
    hashmap_token_insert(&fluent_token_map, "{", TOKEN_OPEN_CURLY);
    hashmap_token_insert(&fluent_token_map, "}", TOKEN_CLOSE_CURLY);
    hashmap_token_insert(&fluent_token_map, ":", TOKEN_COLON);
    hashmap_token_insert(&fluent_token_map, "->", TOKEN_ARROW);
    hashmap_token_insert(&fluent_token_map, "[", TOKEN_OPEN_BRACKET);
    hashmap_token_insert(&fluent_token_map, "]", TOKEN_CLOSE_BRACKET);
    hashmap_token_insert(&fluent_token_map, ".", TOKEN_DOT);

    hashmap_token_insert(&fluent_token_map, "str", TOKEN_STRING);
    hashmap_token_insert(&fluent_token_map, "num", TOKEN_NUM);
    hashmap_token_insert(&fluent_token_map, "dec", TOKEN_DEC);
    hashmap_token_insert(&fluent_token_map, "bool", TOKEN_BOOL);
    hashmap_token_insert(&fluent_token_map, "nothing", TOKEN_NOTHING);

    hashmap_token_insert(&fluent_token_map, "pub", TOKEN_PUB);

    hashmap_token_insert(&fluent_token_map, "true", TOKEN_BOOL_LITERAL);
    hashmap_token_insert(&fluent_token_map, "false", TOKEN_BOOL_LITERAL);

    hashmap_token_insert(&fluent_token_map, "import", TOKEN_IMPORT);
    return &fluent_token_map;
}

#endif //FLUENT_TOKEN_MAP_H
