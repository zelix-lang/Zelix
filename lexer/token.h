/*
        ==== The Fluent Programming Language ====
---------------------------------------------------------
  - This file is part of the Fluent Programming Language
    codebase. Fluent is a fast, statically-typed and
    memory-safe programming language that aims to
    match native speeds while staying highly performant.
---------------------------------------------------------
  - Fluent is categorized as free software; you can
    redistribute it and/or modify it under the terms of
    the GNU General Public License as published by the
    Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.
---------------------------------------------------------
  - Fluent is distributed in the hope that it will
    be useful, but WITHOUT ANY WARRANTY; without even
    the implied warranty of MERCHANTABILITY or FITNESS
    FOR A PARTICULAR PURPOSE. See the GNU General Public
    License for more details.
---------------------------------------------------------
  - You should have received a copy of the GNU General
    Public License along with Fluent. If not, see
    <https://www.gnu.org/licenses/>.
*/

//
// Created by rodrigo on 7/30/25.
//

#pragma once
#include "fluent/container/optional.h"
#include "fluent/container/external_string.h"

namespace fluent::lexer
{
    struct token
    {
        enum t_type
        {
            UNKNOWN,
            IMPORT,
            FUNCTION,
            MOD,
            ARROW,
            OPEN_CURLY,
            CLOSE_CURLY,
            OPEN_PAREN,
            CLOSE_PAREN,
            OPEN_BRACKET,
            CLOSE_BRACKET,
            IDENTIFIER,
            STRING,
            NUMBER,
            DECIMAL,
            NOTHING,
            STRING_LITERAL,
            NUMBER_LITERAL,
            DECIMAL_LITERAL,
            SEMICOLON,
            COMMA,
            COLON,
            EQUALS,
            BOOL_EQ,
            BOOL_NEQ,
            BOOL_LT,
            BOOL_GT,
            BOOL_LTE,
            BOOL_GTE,
            PLUS,
            MINUS,
            MULTIPLY,
            DIVIDE,
            AND,
            OR,
            NOT,
        };

        container::optional<container::external_string> value
            = container::optional<container::external_string>::none();
        t_type type = UNKNOWN;
        size_t line = 0;
        size_t column = 0;
    };
}