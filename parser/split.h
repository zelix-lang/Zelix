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
// Created by rodri on 8/14/25.
//

#pragma once
#include "zelix/container/stream.h"
#include "lexer/token.h"

namespace zelix::parser
{
    inline container::vector<container::stream<lexer::token *>> split_args(
        container::stream<lexer::token*> &tokens
    )
    {
        tokens.next();
        container::vector<container::stream<lexer::token *>> res;
        container::vector<lexer::token *> current_group;
        size_t nested_count = 0;
        auto next_opt = tokens.next();

        while (next_opt.is_some())
        {
            const auto &next = next_opt.get();

            if (next->type == lexer::token::COMMA)
            {
                // Handle the end of a group
                if (nested_count == 0)
                {
                    res.emplace_back(current_group);
                    current_group.clear();
                }
                else
                {
                    current_group.push_back(next);
                }
            }

            else if (next->type == lexer::token::OPEN_PAREN)
            {
                nested_count++;

                if (nested_count > 1)
                {
                    current_group.push_back(next);
                }
            }

            else if (next->type == lexer::token::CLOSE_PAREN)
            {
                if (nested_count == 0)
                {
                    break;
                }

                nested_count--;

                if (nested_count == 1)
                {
                    break;
                }

                current_group.push_back(next);
            }

            else
            {
                current_group.push_back(next);
            }

            next_opt = tokens.next();
        }

        // Append the last group if it contains tokens
        if (!current_group.empty())
        {
            res.emplace_back(current_group);
        }

        return res;
    }
}