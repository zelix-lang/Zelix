/*
        ==== The Zelix Programming Language ====
---------------------------------------------------------
  - This file is part of the Zelix Programming Language
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
#include "ast.h"
#include "zelix/container/stream.h"
#include "lexer/token.h"
#include "memory/allocator.h"

namespace zelix::parser
{
    enum error_type
    {
        NONE,
        ILLEGAL_IMPORT,
        UNEXPECTED_TOKEN
    };

    struct error
    {
        error_type type = NONE;
        size_t line = 0;
        size_t column = 0;
    };

    inline error global_err;

    inline bool is_err() noexcept
    {
        return global_err.type != NONE;
    }

    ast *parse(
        container::stream<lexer::token *> &,
        memory::lazy_allocator<ast> &
    );
}