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
#include "fluent/container/stream.h"
#include "fluent/container/vector.h"
#include "lexer/token.h"
#include "parser/parser.h"

namespace fluent::parser
{
    inline container::vector<lexer::token> extract(
        container::vector<lexer::token> &tokens_vec,
        size_t start,
        const lexer::token::t_type end_delim = lexer::token::OPEN_PAREN,
        const lexer::token::t_type start_delim = lexer::token::OPEN_PAREN,
        const bool handle_nested = true,
        const bool exclude_first_delim = true
    )
    {
        container::vector<lexer::token> result;
        size_t nested_count = 0;

        // Iterate over the tokens
        while (start < tokens_vec.size())
        {
            const auto &current = tokens_vec.ref_at(start);

            if (current.type == end_delim)
            {
                // Handle nested delimiters
                if (handle_nested)
                {
                    if (nested_count == 0)
                    {
                        global_err.type = UNEXPECTED_TOKEN;
                        global_err.column = current.column;
                        global_err.line = current.line;
                        throw except::exception("Unexpected end delimiter");
                    }

                    nested_count--;
                }

                if (nested_count == 0)
                {
                    return result;
                }
            }
            else if (current.type == start_delim)
            {
                if (handle_nested)
                {
                    nested_count++;
                }

                if (exclude_first_delim)
                {
                    start++;
                    continue;
                }
            }

            result.push_back(current);
            start++;
        }

        // Delimiter never closed
        if (start >= tokens_vec.size())
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = 0;
            global_err.line = 0;
            throw except::exception("Unexpected end delimiter");
        }

        const auto &current = tokens_vec.ref_at(start);
        global_err.type = UNEXPECTED_TOKEN;
        global_err.column = current.column;
        global_err.line = current.line;
        throw except::exception("Unexpected end delimiter");
    }

    inline container::vector<lexer::token> extract(
        container::stream<lexer::token> &tokens,
        const lexer::token::t_type end_delim = lexer::token::OPEN_PAREN,
        const lexer::token::t_type start_delim = lexer::token::OPEN_PAREN,
        const bool handle_nested = true,
        const bool exclude_first_delim = true
    )
    {
        return extract(
            tokens.ptr(),
            tokens.pos(),
            end_delim,
            start_delim,
            handle_nested,
            exclude_first_delim
        );
    }
}