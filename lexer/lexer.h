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
#include "fluent/container/stream.h"
#include "memory/allocator.h"
#include "token.h"

namespace zelix::lexer
{
    container::stream<token *> lex(
        const container::external_string &source,
        memory::lazy_allocator<token> &allocator
    );

    enum error_type
    {
        NONE,
        UNKNOWN_TOKEN,
        UNCLOSED_COMMENT,
        UNCLOSED_STRING,
    };

    struct global_err_t
    {
        error_type type = NONE;
        size_t line = 0;
        size_t column = 0;
    };

    inline global_err_t global_err;

    bool is_err() noexcept;
}
