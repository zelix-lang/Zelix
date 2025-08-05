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
        template <bool Arithmetic>
        inline ast::rule_t rule(lexer::token *&token)
        {
            if constexpr (Arithmetic)
            {
                switch (token->type)
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
                        return ast::ROOT; // Default case if no arithmetic operation matches
                }
            }

            switch (token->type)
            {
                case lexer::token::BOOL_EQ:
                    return ast::EQ;
                case lexer::token::BOOL_GT:
                    return ast::GT;
                case lexer::token::BOOL_GTE:
                    return ast::GTE;
                case lexer::token::BOOL_LT:
                    return ast::LT;
                case lexer::token::BOOL_LTE:
                    return ast::LTE;
                case lexer::token::BOOL_NEQ:
                    return ast::NEQ;
                default:
                    return ast::ROOT; // Default case if no arithmetic operation matches
            }
        }

        template <bool Arithmetic>
        inline void process_op(
            lexer::token *&next,
            memory::lazy_allocator<ast> &allocator,
            ast *arithmetic_node
        )
        {
            // Create a new AST node for the arithmetic operation
            ast *arithmetic_op_node = allocator.alloc();
            arithmetic_op_node->rule = rule<Arithmetic>(next);
            arithmetic_node->children.push_back(arithmetic_op_node); // Add the operation node
        }

        template <bool Arithmetic>
        inline void process_sub(
            container::vector<lexer::token *> &current_tokens,
            lexer::token *&next,
            ast *arithmetic_node,
            ast *candidate,
            memory::lazy_allocator<ast> &allocator,
            container::vector<expr::queue_node> &expr_queue,
            bool &first_iteration,
            const bool append_sign = true
        )
        {
            // Handle first iteration
            if (first_iteration)
            {
                first_iteration = false; // Set the flag to false for subsequent iterations
                arithmetic_node->children.push_back(candidate); // Add the candidate as the first child

                if (append_sign)
                {
                    process_op<Arithmetic>(
                        next,
                        allocator,
                        arithmetic_node
                    ); // Process the arithmetic operator
                }

                return;
            }

            // Append the current tokens to the last nested node
            if (current_tokens.empty())
            {
                global_err.type = UNEXPECTED_TOKEN;
                global_err.column = next->column;
                global_err.line = next->line;
                throw except::exception("Unexpected empty arithmetic expression");
            }

            // Create a new expression node for the current tokens
            ast *expr_node = allocator.alloc();
            expr_node->rule = ast::EXPRESSION;
            arithmetic_node->children.push_back(expr_node); // Add the expression node

            // Queue the current tokens for processing
            expr_queue.push_back(
                expr::queue_node(
                    container::stream<lexer::token *>(current_tokens), // Clone the current tokens
                    expr_node
                )
            );
            current_tokens.clear(); // Clear the current tokens for the next operation

            if (!append_sign) return; // If we don't want to append the sign, return early
            process_op<Arithmetic>(
                next,
                allocator,
                arithmetic_node
            ); // Process the arithmetic operator
        }

        template <bool Arithmetic>
        inline void process_high_precedence(
            container::vector<lexer::token *> &current_tokens,
            lexer::token *&next,
            ast *&last_nested,
            ast *&arithmetic_node,
            ast *&candidate,
            memory::lazy_allocator<ast> &allocator,
            container::vector<expr::queue_node> &expr_queue,
            bool &first_iteration
        )
        {
            // Allocate the last nested node if it doesn't exist
            if (last_nested == nullptr)
            {
                last_nested = allocator.alloc();
                last_nested->rule = ast::ARITHMETIC;
                arithmetic_node->children.push_back(last_nested);
            }

            process_sub<Arithmetic>(
                current_tokens,
                next,
                last_nested,
                candidate,
                allocator,
                expr_queue,
                first_iteration
            ); // Process the subexpression
        }

        template <bool Arithmetic>
        inline void process_low_precedence(
            container::vector<lexer::token *> &current_tokens,
            lexer::token *&next,
            ast *&last_nested,
            ast *&arithmetic_node,
            ast *&candidate,
            memory::lazy_allocator<ast> &allocator,
            container::vector<expr::queue_node> &expr_queue,
            bool &first_iteration
        )
        {
            // Check if we have a last nested node
            if (last_nested != nullptr)
            {
                // Process the last nested node
                process_sub<Arithmetic>(
                    current_tokens,
                    next,
                    last_nested,
                    candidate,
                    allocator,
                    expr_queue,
                    first_iteration,
                    false
                ); // Process the subexpression
            }
            else
            {
                process_sub<Arithmetic>(
                    current_tokens,
                    next,
                    arithmetic_node,
                    candidate,
                    allocator,
                    expr_queue,
                    first_iteration,
                    false
                ); // Process the subexpression
            }

            // Set the last nested node to nullptr
            last_nested = nullptr; // Reset the last nested node

            process_op<Arithmetic>(
                next,
                allocator,
                arithmetic_node
            ); // Process the arithmetic operator
        }
    }

    template <bool Arithmetic = true>
    inline ast *signed_op(
        ast *&candidate,
        container::stream<lexer::token *> &tokens,
        memory::lazy_allocator<ast> &allocator,
        container::vector<expr::queue_node> &expr_queue
    )
    {
        // Create the arithmetic AST node
        ast *arithmetic_node = allocator.alloc();
        if constexpr (Arithmetic)
        {
            arithmetic_node->rule = ast::ARITHMETIC;
        }
        else
        {
            arithmetic_node->rule = ast::BOOLEAN;
        }

        ast *last_nested = nullptr; // Last nested arithmetic node
        auto next_opt = tokens.next();
        auto &next = next_opt.get();
        // Tokens for the current subexpression
        container::vector<lexer::token *> current_tokens;
        // Track nested parentheses
        size_t nested_count = 0;
        bool first_iteration = true; // Flag for the first iteration
        bool last_is_mul = false; // Track if the last operation was multiplication/division

        while (next_opt.is_some())
        {
            next = next_opt.get();

            if constexpr (Arithmetic)
            {
                // Check if we have to break
                if (
                    next->type >= lexer::token::BOOL_EQ &&
                    next->type <= lexer::token::BOOL_GTE
                )
                {
                    break;
                }
            }

            // Check if we have an arithmetic operation with high precedence
            constexpr bool high_precedence = Arithmetic ?
                    next->type == lexer::token::MULTIPLY ||
                    next->type == lexer::token::DIVIDE
                :
                    next->type == lexer::token::OR ||
                    next->type == lexer::token::AND;

            constexpr bool low_precedence = Arithmetic ?
                    next->type == lexer::token::PLUS ||
                    next->type == lexer::token::MINUS
                :
                    next->type >= lexer::token::BOOL_EQ &&
                    next->type <= lexer::token::BOOL_GTE;

            if (nested_count == 0 && high_precedence)
            {
                last_is_mul = true;
                arith::process_high_precedence<Arithmetic>(
                    current_tokens,
                    next,
                    last_nested,
                    arithmetic_node,
                    candidate,
                    allocator,
                    expr_queue,
                    first_iteration
                ); // Process the multiplication/division operation
            }
            // Check if we have an arithmetic operator with low precedence
            else if (nested_count == 0 && low_precedence)
            {
                last_is_mul = false;
                arith::process_low_precedence<Arithmetic>(
                    current_tokens,
                    next,
                    last_nested,
                    arithmetic_node,
                    candidate,
                    allocator,
                    expr_queue,
                    first_iteration
                ); // Process the addition/subtraction operation
            }
            // Push the token to the current vector
            else
            {
                // Check if we have a nested expression
                if (next->type == lexer::token::OPEN_PAREN)
                {
                    nested_count++;

                    // Make sure we don't include the opening parenthesis in the current tokens
                    if (nested_count == 1)
                    {
                        next_opt = tokens.next();
                        continue;
                    }
                }
                else if (next->type == lexer::token::CLOSE_PAREN)
                {
                    // Make sure we don't overflow
                    if (nested_count == 0)
                    {
                        global_err.type = UNEXPECTED_TOKEN;
                        global_err.column = next->column;
                        global_err.line = next->line;
                        throw except::exception("Unexpected closing parenthesis in arithmetic expression");
                    }

                    nested_count--;
                    // Make sure we don't include the closing parenthesis in the current tokens
                    if (nested_count == 0)
                    {
                        next_opt = tokens.next();
                        continue;
                    }
                }

                current_tokens.push_back(next); // Add the token to the current tokens
            }

            next_opt = tokens.next();
        }

        // In a valid expression, the last chunk should not be empty
        if (current_tokens.empty())
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = next->column;
            global_err.line = next->line;
            throw except::exception("Unexpected empty arithmetic expression");
        }

        // Process the last arithmetic operation
        const auto last_sub = allocator.alloc();
        last_sub->rule = ast::EXPRESSION;
        if (last_is_mul)
        {
            last_nested->children.push_back(last_sub); // Add the last subexpression to the last nested node
        }
        else
        {
            arithmetic_node->children.push_back(last_sub); // Add the last subexpression to the arithmetic node
        }

        // Queue the current tokens for processing
        expr_queue.push_back(
            expr::queue_node(
                container::stream<lexer::token *>(current_tokens), // Clone the current tokens
                last_sub
            )
        );

        // Unwrap if we have a nested expression
        if (arithmetic_node->children.size() == 1)
        {
            const auto first_child = arithmetic_node->children[0];
            allocator.dealloc(arithmetic_node); // Deallocate the arithmetic node
            arithmetic_node = first_child; // Set the first child as the new node
        }

        return arithmetic_node;

    }
} // namespace fluent::parser::rule