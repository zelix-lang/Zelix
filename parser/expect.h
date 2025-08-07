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
#include "zelix/container/stream.h"
#include "parser.h"

namespace zelix::parser
{
    inline void expect(
        container::stream<lexer::token *> &tokens,
        const lexer::token::t_type expected_type
    )
    {
        auto token_opt = tokens.peek();
        if (token_opt.is_none())
        {
            // Get the current token
            container::optional<lexer::token *> trace = tokens.curr();
            if (trace.is_none())
            {
                global_err.line = 0;
                global_err.column = 0;
            }
            else
            {
                const auto current_token = trace.get();
                global_err.type = UNEXPECTED_TOKEN;
                global_err.column = current_token->column; // Use the column from the current token
                global_err.line = current_token->line; // Use the line from the current token
            }

            throw except::exception("Assertion failed");
        }

        // Get the token
        // Check if the token type matches the expected type
        if (
            const auto &token = token_opt.get();
            token->type != expected_type
        )
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = token->column; // Use the column from the token
            global_err.line = token->line; // Use the line from the token

            // If it doesn't match, return false
            throw except::exception("Assertion failed");
        }
    }
}