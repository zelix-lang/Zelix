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
// Created by rodrigo on 8/1/25.
//

#pragma once
#include "fluent/container/optional.h"
#include "fluent/container/external_string.h"
#include "fluent/container/vector.h"

namespace fluent::parser
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
            CALL,
            PROP_ACCESS,
            IF,
            ELSEIF,
            ELSE,
            FOR,
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
            IDENTIFIER
        };

        rule_t rule = ROOT;
        container::optional<container::external_string> value = container::optional<container::external_string>::none();
        container::vector<ast *> children;
    };
}