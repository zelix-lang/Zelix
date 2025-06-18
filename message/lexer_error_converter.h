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
// Created by rodrigo on 6/5/25.
//

#ifndef FLUENT_LEXER_ERROR_CONVERTER_H
#define FLUENT_LEXER_ERROR_CONVERTER_H

#include <lexer/error.h>

static inline char *lexer_error_to_string(const lexer_error_t *error)
{
    if (!error)
    {
        return "No error";
    }

    switch (error->code)
    {
        case LEXER_ERROR_UNTERMINATED_STRING:
            return "Unterminated string literal";
        case LEXER_ERROR_UNTERMINATED_COMMENT:
            return "Unterminated block comment";
        case LEXER_ERROR_UNKNOWN_TOKEN:
            return "Unknown token encountered";
        case LEXER_ERROR_UNTERMINATED_DEC:
            return "Unterminated decimal number";
        default:
            return "Unknown lexer error";
    }
}

#endif //FLUENT_LEXER_ERROR_CONVERTER_H
