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
#include "zelix/container/stream.h"
#include "lexer/token.h"
#include "parser/parser.h"

namespace zelix::parser
{
    inline container::stream<lexer::token *> extract(
        container::stream<lexer::token*> &tokens,
        const lexer::token::t_type end_delim = lexer::token::CLOSE_PAREN,
        const lexer::token::t_type nested_end_delim = lexer::token::CLOSE_PAREN,
        const lexer::token::t_type start_delim = lexer::token::OPEN_PAREN,
        const bool handle_nested = true,
        const bool exclude_first_delim = true
    )
    {
        container::vector<lexer::token *> vec;
        container::stream result(container::move(vec));
        size_t nested_count = 0;
        const size_t start_pos = tokens.pos();

        // Iterate over the tokens
        auto next_opt = tokens.next();
        while (next_opt.is_some())
        {
            const auto &current = next_opt.get();

            if (current->type == nested_end_delim)
            {
                // Handle nested delimiters
                if (handle_nested)
                {
                    if (nested_count == 0)
                    {
                        tokens.set_pos(start_pos);
                        global_err.type = UNEXPECTED_TOKEN;
                        global_err.column = current->column;
                        global_err.line = current->line;
                        throw except::exception("Unexpected nested end delimiter");
                    }

                    nested_count--;

                    if (nested_count == 0 && nested_end_delim == end_delim)
                    {
                        return result;
                    }
                }
            }

            else if (current->type == end_delim)
            {
                // Handle nested delimiters
                if (handle_nested)
                {
                    if (nested_count != 0)
                    {
                        result.push(current);
                        next_opt = tokens.next();
                        continue; // Continue to the next token
                    }
                }

                return result;
            }

            if (current->type == start_delim)
            {
                if (handle_nested)
                {
                    nested_count++;
                }

                if (exclude_first_delim)
                {
                    next_opt = tokens.next();
                    continue;
                }
            }

            result.push(current);
            next_opt = tokens.next();
        }

        // Restore the original position of the tokens stream
        tokens.set_pos(start_pos);

        // Delimiter never closed
        if (next_opt.is_none())
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = 0;
            global_err.line = 0;
            throw except::exception("Unexpected end delimiter");
        }

        const auto &current = next_opt.get();
        global_err.type = UNEXPECTED_TOKEN;
        global_err.column = current->column;
        global_err.line = current->line;
        throw except::exception("Unexpected end delimiter");
    }
}