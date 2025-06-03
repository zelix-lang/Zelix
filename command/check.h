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

#ifndef FLUENT_COMMAND_CHECK_H
#define FLUENT_COMMAND_CHECK_H

#include <fluent/clock/clock.h>
#include "../file/file_reader.h"
#include "../lexer/lexer.h"

static inline void check_command(const char *const path)
{
    // Read the file contents
    char *source = read_file(path);
    if (!source)
    {
        fprintf(stderr, "Error: Could not read file '%s'.\n", path);
        return;
    }

    hr_clock_t clock;
    hr_clock_tick(&clock); // Start the clock to measure performance
    // Lex the source code
    const pair_lex_result_t result = lexer_tokenize(source);

    // Stop the clock and print the elapsed time
    long long elapsed_time = hr_clock_distance_from_now(&clock, CLOCK_MICROSECONDS);
    printf("%lld microseconds parsing time\n", elapsed_time);
    const token_stream_t stream = result.first; // Get the token stream from the result
    const lexer_error_t *error = result.second; // Get the error if any

    // Check for lexer errors
    if (error)
    {
        fprintf(stderr, "Lexer error at line %zu, column %zu: %d\n",
                error->line, error->column, error->code);
    }
    else
    {
        printf("%lld\n", result.first.tokens->length); // Print the number of tokens
        for (size_t i = 0; i < result.first.tokens->length; i++)
        {
            const token_t *token = vec_get(result.first.tokens, i);
            printf("Token %zu: ", i + 1);
            if (token->value != NULL)
            {
                printf("%s\n", (char *)token->value->ptr);
            } else
            {
                switch (token->type)
                {
                    case TOKEN_UNKNOWN: printf("UNKNOWN\n"); break;
                    case TOKEN_FUNCTION: printf("FUNCTION\n"); break;
                    case TOKEN_LET: printf("LET\n"); break;
                    case TOKEN_CONST: printf("CONST\n"); break;
                    case TOKEN_IF: printf("IF\n"); break;
                    case TOKEN_ELSE: printf("ELSE\n"); break;
                    case TOKEN_ELSE_IF: printf("ELSE_IF\n"); break;
                    case TOKEN_MOD: printf("MOD\n"); break;
                    case TOKEN_RETURN: printf("RETURN\n"); break;
                    case TOKEN_ASSIGN: printf("ASSIGN\n"); break;
                    case TOKEN_PLUS: printf("PLUS\n"); break;
                    case TOKEN_MINUS: printf("MINUS\n"); break;
                    case TOKEN_ASTERISK: printf("ASTERISK\n"); break;
                    case TOKEN_SLASH: printf("SLASH\n"); break;
                    case TOKEN_LESS_THAN: printf("LESS_THAN\n"); break;
                    case TOKEN_GREATER_THAN: printf("GREATER_THAN\n"); break;
                    case TOKEN_EQUAL: printf("EQUAL\n"); break;
                    case TOKEN_NOT_EQUAL: printf("NOT_EQUAL\n"); break;
                    case TOKEN_GREATER_THAN_OR_EQUAL: printf("GREATER_THAN_OR_EQUAL\n"); break;
                    case TOKEN_LESS_THAN_OR_EQUAL: printf("LESS_THAN_OR_EQUAL\n"); break;
                    case TOKEN_ARROW: printf("ARROW\n"); break;
                    case TOKEN_COMMA: printf("COMMA\n"); break;
                    case TOKEN_SEMICOLON: printf("SEMICOLON\n"); break;
                    case TOKEN_OPEN_PAREN: printf("OPEN_PAREN\n"); break;
                    case TOKEN_CLOSE_PAREN: printf("CLOSE_PAREN\n"); break;
                    case TOKEN_OPEN_CURLY: printf("OPEN_CURLY\n"); break;
                    case TOKEN_CLOSE_CURLY: printf("CLOSE_CURLY\n"); break;
                    case TOKEN_COLON: printf("COLON\n"); break;
                    case TOKEN_NOT: printf("NOT\n"); break;
                    case TOKEN_OR: printf("OR\n"); break;
                    case TOKEN_AND: printf("AND\n"); break;
                    case TOKEN_OPEN_BRACKET: printf("OPEN_BRACKET\n"); break;
                    case TOKEN_CLOSE_BRACKET: printf("CLOSE_BRACKET\n"); break;
                    case TOKEN_DOT: printf("DOT\n"); break;
                    case TOKEN_STRING: printf("STRING\n"); break;
                    case TOKEN_NUM: printf("NUM\n"); break;
                    case TOKEN_DEC: printf("DEC\n"); break;
                    case TOKEN_NOTHING: printf("NOTHING\n"); break;
                    case TOKEN_BOOL: printf("BOOL\n"); break;
                    case TOKEN_STRING_LITERAL: printf("STRING_LITERAL\n"); break;
                    case TOKEN_NUM_LITERAL: printf("NUM_LITERAL\n"); break;
                    case TOKEN_DECIMAL_LITERAL: printf("DECIMAL_LITERAL\n"); break;
                    case TOKEN_BOOL_LITERAL: printf("BOOL_LITERAL\n"); break;
                    case TOKEN_WHILE: printf("WHILE\n"); break;
                    case TOKEN_FOR: printf("FOR\n"); break;
                    case TOKEN_NEW: printf("NEW\n"); break;
                    case TOKEN_IN: printf("IN\n"); break;
                    case TOKEN_TO: printf("TO\n"); break;
                    case TOKEN_BREAK: printf("BREAK\n"); break;
                    case TOKEN_CONTINUE: printf("CONTINUE\n"); break;
                    case TOKEN_PUB: printf("PUB\n"); break;
                    case TOKEN_AMPERSAND: printf("AMPERSAND\n"); break;
                    case TOKEN_BAR: printf("BAR\n"); break;
                    case TOKEN_IMPORT: printf("IMPORT\n"); break;
                    case TOKEN_IDENTIFIER: printf("IDENTIFIER\n"); break;
                    default: printf("INVALID_TOKEN_TYPE\n"); break;
                }
            }
        }
    }

    // Free the resources used by the lexer
    if (stream.allocator)
    {
        // Destroy the vector of tokens
        vec_destroy(stream.tokens, NULL);
        destroy_arena(stream.allocator);
    }

    // Free the allocated memory
    free(source);
}

#endif //FLUENT_COMMAND_CHECK_H
