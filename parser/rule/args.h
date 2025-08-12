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
// Created by rodrigo on 8/2/25.
//

#pragma once

#include "extractor.h"
#include "zelix/container/stream.h"
#include "lexer/token.h"
#include "memory/allocator.h"
#include "parser/rule/expr/queue.h"

namespace zelix::parser::rule
{
    inline void args(
        ast *&root,
        container::stream<lexer::token *> &tokens,
        memory::lazy_allocator<ast> &allocator,
        memory::lazy_allocator<expr::queue_node> &q_allocator,
        container::vector<expr::queue_node *> &expr_queue
    )
    {
        // Create a new AST node for the arguments
        ast *args_node = allocator.alloc();
        args_node->rule = ast::ARGUMENTS;

        // Extract all tokens until the end of args
        auto args_group = extract(tokens);

        // Iterate over the extracted tokens
        while (true)
        {
            // Allocate a new AST node for the argument
            ast *arg_node = allocator.alloc();
            arg_node->rule = ast::ARGUMENT;
            args_node->children.push_back(arg_node);

            try
            {
                // Extract until the next comma or end of args
                auto arg_group = extract(
                    args_group,
                    lexer::token::COMMA,
                    lexer::token::CLOSE_PAREN,
                    lexer::token::OPEN_PAREN,
                    true, // Handle nested delimiters
                    false // Do not exclude the first delimiter
                );

                if (arg_group.empty())
                {
                    continue;
                }

                // Push the argument group to the args node
                auto q_el = q_allocator.alloc();
                q_el->tokens = container::move(arg_group);
                q_el->node = arg_node;
                expr_queue.emplace_back(q_el);
            }
            catch (const except::exception &_)
            {
                // If the delimiter was not found, we are at the last argument
                // and the position was restored
                auto q_el = q_allocator.alloc();
                q_el->tokens = container::move(args_group);
                q_el->node = arg_node;
                expr_queue.emplace_back(q_el);

                break; // Exit the loop if we reach the end of args
            }
        }

        // Append the args node to the root
        root->children.push_back(args_node);
    }
}