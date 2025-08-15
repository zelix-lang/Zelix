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
// Created by rodrigo on 8/6/25.
//

#pragma once

#include "expr/expr.h"
#include "extractor.h"
#include "zelix/container/stream.h"
#include "lexer/token.h"
#include "memory/allocator.h"
#include "parser/expect.h"
#include "parser/parser.h"

namespace zelix::parser::rule
{
    inline void for_loop(
        ast *&root,
        container::stream<lexer::token *> &tokens,
        memory::lazy_allocator<ast> &allocator,
        const lexer::token *const &trace
    )
    {
        // Create a new AST node for the for loop
        ast *for_node = allocator.alloc();
        for_node->rule = ast::FOR;
        root->children.push_back(for_node);

        // Expect an identifier for the loop variable
        expect(tokens, lexer::token::IDENTIFIER);
        const auto loop_var = tokens.next().get();
        ast *loop_var_ast = allocator.alloc();
        loop_var_ast->rule = ast::IDENTIFIER;
        loop_var_ast->value = loop_var->value;
        for_node->children.push_back(loop_var_ast); // Add the loop variable as a child of the for node

        // Expect the 'in' keyword
        expect(tokens, lexer::token::IN);
        tokens.next(); // Consume the 'in' keyword

        // Extract all tokens until the next "to"
        auto range_group = extract(
            tokens,
            lexer::token::TO,
            lexer::token::UNKNOWN,
            lexer::token::UNKNOWN,
            false,
            false
        );

        // Create the AST node for the range
        ast *from_range_node = allocator.alloc();
        from_range_node->rule = ast::FROM; // Set the rule to FROM
        for_node->children.push_back(from_range_node); // Add the range node as a child of the for node

        rule::expression<false, true>(
            from_range_node,
            range_group,
            allocator,
            trace
        ); // Parse the range expression

        // Collect the remaining expression
        container::vector<lexer::token *> expr_tokens;
        auto next_opt = tokens.next();
        if (next_opt.is_none())
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = trace->column;
            global_err.line = trace->line;
            throw except::exception("Unexpected end of input while parsing for loop range");
        }

        // The node to queue the expression with
        ast *queue_node = allocator.alloc();
        queue_node->rule = ast::TO; // Set the rule to "TO"
        queue_node->children.push_back(from_range_node); // Add the "from" range node as a child
        for_node->children.push_back(queue_node); // Add the queue node to the for node

        while (next_opt.is_some())
        {
            auto &next = next_opt.get();

            // Parse step
            if (next->type == lexer::token::STEP)
            {
                auto st = container::stream<lexer::token *>(expr_tokens);

                rule::expression<false, true>(
                    queue_node,
                    st,
                    allocator,
                    trace
                ); // Parse the range expression

                expr_tokens.clear(); // Clear the expression tokens

                // Create the AST node for the step
                ast *step_node = allocator.alloc();
                step_node->rule = ast::STEP; // Set the rule to STEP
                for_node->children.push_back(step_node); // Add the step node as a child of the for node
                queue_node = step_node; // Update the queue node to the step node
                continue;
            }

            expr_tokens.push_back(next); // Add the token to the expression tokens

            // Peek into the next token
            auto peek_opt = tokens.peek();
            if (peek_opt.is_some() && peek_opt.get()->type == lexer::token::OPEN_CURLY)
            {
                // We have reached the end of the range expression
                break;
            }

            next_opt = tokens.next();
        }

        if (!expr_tokens.empty())
        {
            auto st = container::stream<lexer::token *>(expr_tokens);
            rule::expression<false, true>(
                queue_node,
                st,
                allocator,
                trace
            ); // Parse the range expression
        }
    }
}