/*
        ==== The Zelix Programming Language ====
---------------------------------------------------------
  - This file is part of the Zelix Programming Language
    codebase. Zelix is a fast, statically-typed and
    memory-safe programming language that aims to
    match native speeds while staying highly performant.
---------------------------------------------------------
  - Zelix is categorized as free software; you can
    redistribute it and/or modify it under the terms of
    the GNU General Public License as published by the
    Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.
---------------------------------------------------------
  - Zelix is distributed in the hope that it will
    be useful, but WITHOUT ANY WARRANTY; without even
    the implied warranty of MERCHANTABILITY or FITNESS
    FOR A PARTICULAR PURPOSE. See the GNU General Public
    License for more details.
---------------------------------------------------------
  - You should have received a copy of the GNU General
    Public License along with Zelix. If not, see
    <https://www.gnu.org/licenses/>.
*/

//
// Created by rodrigo on 7/30/25.
//

#pragma once
#include "zelix/container/optional.h"
#include "zelix/container/external_string.h"

namespace zelix::lexer
{
    struct token
    {
        enum t_type
        {
            UNKNOWN, // 0
            IMPORT, // 1
            FUNCTION, // 2
            MOD, // 3
            ARROW, // 4
            OPEN_CURLY, // 5
            CLOSE_CURLY, // 6
            OPEN_PAREN, // 7
            CLOSE_PAREN, // 8
            OPEN_BRACKET, // 9
            CLOSE_BRACKET, // 10
            IDENTIFIER, // 11
            STRING, // 12
            NUMBER, // 13
            DECIMAL, // 14
            NOTHING, // 15
            STRING_LITERAL, // 16
            NUMBER_LITERAL, // 17
            DECIMAL_LITERAL, // 18
            SEMICOLON, // 19
            COMMA, // 20
            COLON, // 21
            EQUALS, // 22
            BOOL_EQ, // 23
            BOOL_NEQ, // 24
            BOOL_LT, // 25
            BOOL_GT, // 26
            BOOL_LTE, // 27
            BOOL_GTE, // 28
            PLUS, // 29
            MINUS, // 30
            MULTIPLY, // 31
            DIVIDE, // 32
            AND, // 33
            OR, // 34
            NOT, // 35
            DOT, // 36
            STEP, // 37
            BOOL, // 38
            TRUE, // 39
            FALSE, // 40
            LET, // 41
            CONST, // 42
            PUB, // 43
            IF, // 44
            ELSEIF, // 45
            ELSE, // 46
            FOR, // 47
            WHILE, // 48
            RETURN, // 49
            TO, // 50
            IN, // 51
            AMPERSAND // 51
        };

        container::optional<container::external_string> value
            = container::optional<container::external_string>::none();
        t_type type = UNKNOWN;
        size_t line = 0;
        size_t column = 0;
    };
}