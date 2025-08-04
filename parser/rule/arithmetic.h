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
#include "parser/rule/expr/queue.h"

namespace fluent::parser::rule
{
    namespace arith
    {
        static inline ast::rule_t rule(const lexer::token::t_type type)
        {
            switch (type)
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
                    throw except::exception("Invalid arithmetic operator");
            }
        }

        static inline bool is_boolean_op(lexer::token::t_type type)
        {
            return type == lexer::token::BOOL_EQ || type == lexer::token::BOOL_NEQ || type == lexer::token::BOOL_LT ||
                   type == lexer::token::BOOL_GT || type == lexer::token::BOOL_LTE || type == lexer::token::BOOL_GTE;
        }

        inline ast *parse_term(
            ast *left,
            container::stream<lexer::token> &term_tokens,
            container::vector<expr::queue_node> &expr_queue,
            memory::lazy_allocator<ast> &allocator
        )
        {
            ast *result = left;
            while (!term_tokens.empty())
            {
                auto current_opt = term_tokens.peek();
                if (current_opt.is_none())
                {
                    break; // No more tokens, exit the loop
                }

                const auto &current = current_opt.get();
                if (current.type != lexer::token::MULTIPLY && current.type != lexer::token::DIVIDE)
                {
                    break; // Exit if not * or /
                }
                term_tokens.next(); // Consume operator

                // Get the next operand
                if (term_tokens.empty())
                {
                    throw except::exception("Expected operand after operator");
                }

                auto next_opt = term_tokens.peek();
                if (next_opt.is_none())
                {
                    throw except::exception("Expected operand after operator");
                }

                const auto &next = next_opt.get();
                ast *right = nullptr;
                if (next.type == lexer::token::OPEN_PAREN)
                {
                    // Handle subexpression: push to expr_queue
                    term_tokens.next(); // Consume '('
                    container::vector<lexer::token> sub_tokens;
                    int paren_count = 1;
                    while (!term_tokens.empty() && paren_count > 0)
                    {
                        auto token_opt = term_tokens.peek();
                        if (token_opt.is_none())
                        {
                            throw std::runtime_error("Mismatched parentheses");
                        }

                        const auto &t = token_opt.get();
                        term_tokens.next();
                        if (t.type == lexer::token::OPEN_PAREN)
                            paren_count++;
                        else if (t.type == lexer::token::CLOSE_PAREN)
                            paren_count--;
                        if (paren_count > 0 || t.type != lexer::token::CLOSE_PAREN)
                        {
                            sub_tokens.push_back(t);
                        }
                    }

                    if (paren_count != 0)
                    {
                        throw except::exception("Mismatched parentheses");
                    }

                    right = allocator.alloc();
                    right->rule = ast::EXPRESSION;
                    expr_queue.emplace_back(container::stream(sub_tokens), right);
                }
                else if (next.type == lexer::token::NUMBER || next.type == lexer::token::DECIMAL ||
                         next.type == lexer::token::NUMBER_LITERAL || next.type == lexer::token::DECIMAL_LITERAL)
                {
                    // Handle numeric literal
                    right = allocator.alloc();
                    right->rule = (next.type == lexer::token::NUMBER || next.type == lexer::token::NUMBER_LITERAL)
                                          ? ast::NUMBER_LITERAL
                                          : ast::DECIMAL_LITERAL;
                    right->value = next.value;
                    term_tokens.next();
                }
                else if (next.type == lexer::token::IDENTIFIER)
                {
                    // Handle identifier
                    right = allocator.alloc();
                    right->rule = ast::IDENTIFIER;
                    right->value = next.value;
                    term_tokens.next();
                }
                else
                {
                    throw except::exception("Expected number, identifier, or subexpression");
                }

                // Create new AST node for the operator
                ast *op_node = allocator.alloc();
                op_node->rule = rule(current.type);
                op_node->children.push_back(result);
                op_node->children.push_back(right);
                result = op_node;
            }
            return result;
        }
    } // namespace arith

    inline ast *arithmetic(
        ast *&candidate,
        container::stream<lexer::token> &tokens,
        memory::lazy_allocator<ast> &allocator,
        container::vector<expr::queue_node> &expr_queue
    )
    {
        ast *result = arith::parse_term(candidate, tokens, expr_queue, allocator);
        while (!tokens.empty())
        {
            auto current_opt = tokens.peek();
            if (current_opt.is_none())
            {
                break; // No more tokens, exit the loop
            }

            const auto &current = current_opt.get();
            if (arith::is_boolean_op(current.type))
            {
                break; // Stop at boolean operator
            }
            if (current.type != lexer::token::PLUS && current.type != lexer::token::MINUS)
            {
                break; // Exit if not + or -
            }
            tokens.next(); // Consume operator

            // Get the next operand
            if (tokens.empty())
            {
                throw except::exception("Expected operand after operator");
            }

            auto next_opt = tokens.peek();
            if (next_opt.is_none())
            {
                throw except::exception("Expected operand after operator");
            }

            const auto &next = next_opt.get();
            ast *right = nullptr;
            if (next.type == lexer::token::OPEN_PAREN)
            {
                // Handle subexpression: push to expr_queue
                tokens.next(); // Consume '('
                container::vector<lexer::token> sub_tokens;
                int paren_count = 1;
                while (!tokens.empty() && paren_count > 0)
                {
                    auto t_opt = tokens.peek();
                    if (t_opt.is_none())
                    {
                        throw except::exception("Mismatched parentheses");
                    }

                    const auto &t = t_opt.get();
                    tokens.next();
                    if (t.type == lexer::token::OPEN_PAREN)
                        paren_count++;
                    else if (t.type == lexer::token::CLOSE_PAREN)
                        paren_count--;
                    if (paren_count > 0 || t.type != lexer::token::CLOSE_PAREN)
                    {
                        sub_tokens.push_back(t);
                    }
                }
                if (paren_count != 0)
                {
                    throw std::runtime_error("Mismatched parentheses");
                }
                right = allocator.alloc();
                right->rule = ast::EXPRESSION;
                expr_queue.emplace_back(container::stream(sub_tokens), right);
            }
            else if (next.type == lexer::token::NUMBER || next.type == lexer::token::DECIMAL ||
                     next.type == lexer::token::NUMBER_LITERAL || next.type == lexer::token::DECIMAL_LITERAL)
            {
                // Handle numeric literal
                right = allocator.alloc();
                right->rule = (next.type == lexer::token::NUMBER || next.type == lexer::token::NUMBER_LITERAL)
                                      ? ast::NUMBER_LITERAL
                                      : ast::DECIMAL_LITERAL;
                right->value = next.value;
                tokens.next();
            }
            else if (next.type == lexer::token::IDENTIFIER)
            {
                // Handle identifier
                right = allocator.alloc();
                right->rule = ast::IDENTIFIER;
                right->value = next.value;
                tokens.next();
            }
            else
            {
                throw except::exception("Expected number, identifier, or subexpression");
            }

            // Parse any following * or / operators with the right operand
            right = arith::parse_term(right, tokens, expr_queue, allocator);

            // Create new AST node for the operator
            ast *op_node = allocator.alloc();
            op_node->rule = arith::rule(current.type);
            op_node->children.push_back(result);
            op_node->children.push_back(right);
            result = op_node;
        }

        return result;
    }
} // namespace fluent::parser::rule
