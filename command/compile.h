/*
        ==== The Zelix Programming Language ====
---------------------------------------------------------
  - This file is part of the Zelix Programming Language
    codebase. Fluent is a fast, statically-typed and
    memory-safe programming language that aims to
    match native speeds while staying highly performant.
---------------------------------------------------------
  - Zelix is categorized as free software; you can
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
// Created by rodrigo on 8/4/25.
//

#pragma once
#include "zelix/cli/app.h"
#include "lexer/lexer.h"
#include "memory/allocator.h"
#include "parser/parser.h"
#include "parser/rule/expr/expr.h"
#include "time/timed_task.h"
#include "util/read_file.h"

namespace zelix::command
{
    inline void compile(cli::args &args)
    {
        try
        {
            time::post("Reading", 1);
            auto f = util::read_file(
                args.command<container::external_string>(
                    container::external_string("compile", 7)
                ).ptr()
            );

            time::complete();
            memory::lazy_allocator<lexer::token> token_allocator;
            time::post("Lexing", 1);
            auto tokens = lexer::lex(
                container::external_string(
                    f.c_str(),
                    f.size()
                ),
                token_allocator
            );

            time::complete();
            memory::lazy_allocator<parser::ast> allocator;
            time::post("Parsing", 1);
            parser::parse(
                tokens,
                allocator
            );
            time::complete();
        }
        catch (const except::exception &e)
        {
            time::fail(e.what());
        }
    }
}