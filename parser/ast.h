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
// Created by rodrigo on 8/1/25.
//

#pragma once
#include "zelix/container/optional.h"
#include "zelix/container/external_string.h"
#include "zelix/container/vector.h"

namespace zelix::parser
{
    struct ast
    {
        enum rule_t
        {
            ROOT,
            IMPORT,
            FUNCTION,
            MOD,
            TYPE,
            ARGUMENTS,
            ARGUMENT,
            BLOCK,
            DECLARATION,
            CONST_DECLARATION,
            BOOLEAN_OPERATION,
            SUM,
            SUB,
            MUL,
            DIV,
            EQ,
            NEQ,
            GT,
            GTE,
            LT,
            LTE,
            EXPRESSION,
            ARITHMETIC,
            BOOLEAN,
            CALL,
            PROP_ACCESS,
            IF,
            ELSEIF,
            ELSE,
            FOR,
            FROM,
            TO,
            IN,
            STEP,
            WHILE,
            STR,
            NUM,
            DEC,
            BOOL,
            NOTHING,
            STRING_LITERAL,
            NUMBER_LITERAL,
            DECIMAL_LITERAL,
            TRUE,
            FALSE,
            IDENTIFIER,
            PUBLIC,
            PTR,
            DEREF,
            DERIVE,
            ASSIGNMENT
        };

        rule_t rule = ROOT;
        container::optional<container::external_string> value = container::optional<container::external_string>::none();
        container::vector<ast *> children;
        size_t line = 0;
        size_t column = 0;
    };
}