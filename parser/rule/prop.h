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
#include "zelix/container/stream.h"
#include "lexer/token.h"
#include "memory/allocator.h"
#include "parser/expect.h"
#include "parser/rule/call.h"

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

        while (true)
        {
            // Get the next token
            auto next_opt = tokens.next();
            if (next_opt.is_none())
            {
                global_err.type = UNEXPECTED_TOKEN;
                global_err.column = trace->column;
                global_err.line = trace->line;
                throw except::exception("Unexpected end of input while parsing property access");
            }

            if (next_opt.get()->type != lexer::token::DOT)
            {
                break; // If the next token is not a dot, exit the loop
            }

            // Expect an identifier token
            expect(tokens, lexer::token::IDENTIFIER);
            const auto prop_name = tokens.next().get();
            ast *prop_name_node = allocator.alloc();
            prop_name_node->rule = ast::IDENTIFIER;
            prop_name_node->value = prop_name->value;

            // Peek into the next token
            next_opt = tokens.peek();
            if (next_opt.is_none())
            {
                break; // No more tokens, exit the loop
            }

            if (
                const auto &next = next_opt.get();
                next->type == lexer::token::OPEN_PAREN
            )
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
                prop_node->children.push_back(prop_name_node);
            }
        }

        return prop_node; // Return the property access node
    }
}