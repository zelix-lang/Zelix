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
// Created by rodri on 8/11/25.
//

#pragma once
#include "../file_code.h"
#include "lexer/token.h"
#include "memory/allocator.h"
#include "zelix/container/vector.h"

namespace zelix::code
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

    namespace converter
    {
        struct queue_el
        {
            parser::ast *root;
            container::string path;
            container::string content;
        };
    }

    container::vector<file_code *> convert(
        memory::lazy_allocator<file_code> &allocator,
        memory::lazy_allocator<parser::ast> &ast_allocator,
        memory::lazy_allocator<lexer::token> &token_allocator,
        memory::lazy_allocator<function> &fun_allocator,
        memory::lazy_allocator<mod> &mod_allocator,
        parser::ast *const &root,
        container::string &root_path,
        container::string &root_content
    );
}