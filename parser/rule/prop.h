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
#include "call/call.h"
#include "lexer/token.h"
#include "memory/allocator.h"
#include "zelix/container/stream.h"

namespace zelix::parser::rule
{
    inline ast *prop(
        ast *&candidate,
        container::stream<lexer::token *> &tokens,
        memory::lazy_allocator<ast> &allocator,
        const lexer::token *const &trace,
        memory::lazy_allocator<expr::queue_node> &q_allocator,
        container::vector<expr::queue_node *> &expr_queue
    )
    {
        // Create a new AST node for the property access
        ast *prop_node = allocator.alloc();
        prop_node->rule = ast::PROP_ACCESS;
        prop_node->children.push_back(candidate); // Add the candidate as the first child
        bool allow_dot = true;

        while (true)
        {
            // Get the next token
            auto next_opt = tokens.next();
            if (next_opt.is_none())
            {
                if (!allow_dot)
                {
                    global_err.type = UNEXPECTED_TOKEN;
                    global_err.column = trace->column;
                    global_err.line = trace->line;
                    throw except::exception("Unexpected dot in property access");
                }

                break;
            }

            const auto &next = next_opt.get();
            if (next->type == lexer::token::DOT)
            {
                if (!allow_dot)
                {
                    global_err.type = UNEXPECTED_TOKEN;
                    global_err.column = next->column;
                    global_err.line = next->line;
                    throw except::exception("Unexpected dot in property access");
                }

                allow_dot = false;
                continue;
            }

            // Expect an identifier token
            if (next->type != lexer::token::IDENTIFIER)
            {
                global_err.type = UNEXPECTED_TOKEN;
                global_err.column = next->column;
                global_err.line = next->line;
                throw except::exception("Expected identifier in property access");
            }

            ast *prop_name_node = allocator.alloc();
            prop_name_node->rule = ast::IDENTIFIER;
            prop_name_node->line = next->line;
            prop_name_node->column = next->column;
            prop_name_node->value = next->value;

            // Peek into the next token
            next_opt = tokens.peek();

            if (next_opt.is_some() && next_opt.get()->type == lexer::token::OPEN_PAREN)
            {
                // Create a call token
                ast *prop_call_node = allocator.alloc();
                prop_call_node->rule = ast::CALL;
                prop_call_node->children.push_back(prop_name_node);

                // Parse the arguments
                args(prop_call_node, tokens, allocator, q_allocator, expr_queue);
                prop_node->children.push_back(prop_call_node);
            }
            else
            {
                tokens.next(); // Consume the identifier token
                prop_node->children.push_back(prop_name_node);
            }

            allow_dot = true;
        }

        return prop_node; // Return the property access node
    }
}