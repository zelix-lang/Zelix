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
// Created by rodrigo on 8/2/25.
//

#pragma once

#include "fluent/container/stream.h"
#include "lexer/token.h"
#include "memory/allocator.h"
#include "parser/parser.h"
#include "parser/rule/expr/queue.h"

namespace fluent::parser::rule
{
    namespace arith
    {
        static inline int get_precedence(const lexer::token::t_type type)
        {
            switch (type)
            {
                case lexer::token::PLUS:
                case lexer::token::MINUS:
                    return 1; // Lowest precedence
                case lexer::token::MULTIPLY:
                case lexer::token::DIVIDE:
                    return 2; // Highest precedence
                default:
                    return 0; // No precedence
            }
        }
        inline ast::rule_t rule(
            const lexer::token &op
        )
        {
            switch (op.type)
            {
                case lexer::token::PLUS:
                    return ast::SUM;
                case lexer::token::MINUS:
                    return ast::SUB;
                case lexer::token::MULTIPLY:
                    return ast::MUL;
                case lexer::token::DIVIDE:
                    return ast::DIV;
                default:
                    return ast::ROOT; // Should not happen, as this is validated earlier
            }
        }

        template <typename Complex>
        static inline container::vector<lexer::token> collect(
            container::stream<lexer::token> &tokens
        )
        {
            // No special precedence, just add the next token as a child
            container::vector<lexer::token> vec;
            auto next_opt = tokens.peek();
            size_t nested_count = 0;

            while (next_opt.is_some())
            {
                const auto &next = next_opt.get();

                // Handle nested expressions
                if (next.type == lexer::token::OPEN_PAREN)
                {
                    // Do not use extract() directly, since it will create
                    // another stream
                    nested_count++;
                }
                else if (next.type == lexer::token::CLOSE_PAREN)
                {
                    if (nested_count == 0)
                    {
                        global_err.type = UNEXPECTED_TOKEN;
                        global_err.column = next.column;
                        global_err.line = next.line;
                        throw except::exception("Unexpected nested end delimiter");
                    }

                    nested_count--;
                }

                // Ignore nested expressions
                if (nested_count != 0)
                {
                    vec.push_back(next);
                    continue;
                }

                if constexpr (std::is_same_v<Complex, bool>)
                {
                    if (
                        // Break if we find arithmetic operators OR boolean operators
                        next.type == lexer::token::AND ||
                        next.type == lexer::token::OR ||
                        next.type == lexer::token::NOT ||
                        next.type == lexer::token::BOOL_EQ ||
                        next.type == lexer::token::BOOL_NEQ ||
                        next.type == lexer::token::BOOL_GT ||
                        next.type == lexer::token::BOOL_GTE ||
                        next.type == lexer::token::BOOL_LT ||
                        next.type == lexer::token::BOOL_LTE ||
                        next.type == lexer::token::PLUS ||
                        next.type == lexer::token::MINUS ||
                        next.type == lexer::token::MULTIPLY ||
                        next.type == lexer::token::DIVIDE
                    )
                    {
                        // If we encounter another arithmetic operator, break the loop
                        break;
                    }
                }
                else
                {
                    if (
                        // Break if we find boolean operators, which are
                        // the only possible tokens after an arithmetic operator
                        next.type == lexer::token::AND ||
                        next.type == lexer::token::OR ||
                        next.type == lexer::token::NOT ||
                        next.type == lexer::token::BOOL_EQ ||
                        next.type == lexer::token::BOOL_NEQ ||
                        next.type == lexer::token::BOOL_GT ||
                        next.type == lexer::token::BOOL_GTE ||
                        next.type == lexer::token::BOOL_LT ||
                        next.type == lexer::token::BOOL_LTE
                    )
                    {
                        // If we encounter another arithmetic operator, break the loop
                        break;
                    }
                }

                vec.push_back(next);
                next_opt = tokens.next(); // Get the next token
            }

            return vec;
        }
    }

    inline ast *arithmetic(
        ast *&candidate,
        container::stream<lexer::token> &tokens,
        memory::lazy_allocator<ast> &allocator,
        container::vector<expr::queue_node> &expr_queue
    )
    {
        // Get the next token
        auto next_opt = tokens.next();
        const auto &next = next_opt.get();
        // Create a new AST node for the arithmetic operation
        ast *arithmetic_node = allocator.alloc();
        arithmetic_node->rule = arith::rule(next); // Set the rule based on the operator type
        arithmetic_node->children.push_back(candidate); // Add the candidate as the first child

        // Since the arithmetic operator is validated earlier, possible tokens are:
        // lexer::token::PLUS
        // lexer::token::MINUS
        // lexer::token::MULTIPLY
        // lexer::token::DIVIDE
        // Check if we have an operator with special precedence
        if (
            next.type == lexer::token::MULTIPLY ||
            next.type == lexer::token::DIVIDE
        )
        {
            while (true)
            {
                // Collect the tokens until the next arithmetic operator
                auto tokens_group = arith::collect<float>(tokens);
                auto curr_opt = tokens.curr();
                if (tokens_group.empty() || curr_opt.is_none())
                {
                    break; // No more tokens, exit the loop
                }

                // Check if the current token is an arithmetic operator
                // with the same precedence
                if (
                    const auto &curr = curr_opt.get();
                    curr.type == lexer::token::MULTIPLY ||
                    curr.type == lexer::token::DIVIDE
                )
                {
                    // Create a new AST node for the arithmetic operation
                    ast *expr_node = allocator.alloc();
                    expr_node->rule = arith::rule(curr); // Set the rule based on the operator type
                    expr_node->children.push_back(arithmetic_node); // Add the previous node as a child
                    arithmetic_node = expr_node; // Update the arithmetic node to the new one

                    // Consume the operator token
                    tokens.next();
                }
                else
                {
                    // If we encounter a different operator, break the loop
                    break;
                }
            }
        }
        else
        {
            // Allocate a new expression for the arithmetic operation
            ast *expr_node = allocator.alloc();
            expr_node->rule = ast::EXPRESSION;
            arithmetic_node->children.push_back(expr_node);

            // Append the expression to the queue
            expr_queue.emplace_back(
                arith::collect<float>(tokens), // Collect the tokens until the next arithmetic operator
                expr_node
            );
        }

        return arithmetic_node;
    }
} // namespace fluent::parser::rule
