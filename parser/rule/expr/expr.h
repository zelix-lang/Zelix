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
#include "lexer/token.h"
#include "likely.h"
#include "memory/allocator.h"
#include "parser/parser.h"
#include "parser/rule/arithmetic.h"
#include "parser/rule/call.h"
#include "parser/rule/extractor.h"
#include "parser/rule/prop.h"
#include "queue.h"

namespace fluent::parser::rule
{
    static inline bool process_next(
        ast *&parent,
        ast *&candidate,
        const lexer::token &trace,
        container::optional<lexer::token> &first_opt,
        lexer::token &first
    )
    {
        if (first_opt.is_none())
        {
            // Push the candidate to the current node
            if (candidate != nullptr)
            {
                parent->children.push_back(candidate);
                return true;
            }

            // Set error state
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = trace.column;
            global_err.line = trace.line;
            throw except::exception("Unexpected end of expression");
        }

        first = first_opt.get();
        return false;
    }

    inline void expression(
        ast *&root,
        container::stream<lexer::token> &tokens,
        memory::lazy_allocator<ast> &allocator,
        const lexer::token &trace
    )
    {
        // Create a new AST node for the expression
        ast *expr_node = allocator.alloc();
        expr_node->rule = ast::EXPRESSION;
        root->children.push_back(expr_node);

        // Extract all tokens until the next semicolon
        auto expr_group = extract(
            tokens,
            lexer::token::SEMICOLON,
            lexer::token::UNKNOWN,
            lexer::token::UNKNOWN,
            false, // Do not handle nested delimiters
            false // Do not exclude the first delimiter
        );

        // Create a queue for nested expressions
        container::vector<expr::queue_node> expr_queue;
        expr_queue.emplace_back(expr_group, expr_node);

        // Process all expressions
        while (!expr_queue.empty())
        {
            auto &[expr_stream, node] = expr_queue.ref_at(expr_queue.size() - 1);
            expr_queue.pop_back(); // Remove the current expression from the queue

            auto first_opt = expr_stream.peek();
            // Check if the current element is empty
            if (first_opt.is_none())
            {
                global_err.type = UNEXPECTED_TOKEN;
                global_err.column = trace.column;
                global_err.line = trace.line;
                throw except::exception("Unexpected empty expression");
            }


            // Get the first token to determine the expression type
            ast *candidate = nullptr;
            // Define the likely operations
#       if defined(__x86_64__) || defined(_M_X64) || defined(__aarch64__)
            uint64_t likely = 0;
#       else
            uint32_t likely = 0;
#       endif

            auto &first = first_opt.get();
            switch (first.type)
            {
                case lexer::token::IDENTIFIER:
                {
                    likely |= expr::ALL_LIKELY; // All operations are likely

                    // Update the candidate
                    candidate = allocator.alloc();
                    candidate->rule = ast::IDENTIFIER;
                    candidate->value = first.value;
                    expr_stream.next(); // Consume the identifier token
                    break;
                }

                case lexer::token::NUMBER_LITERAL:
                {
                    // Arithmetic operation is likely
                    likely |= expr::ARITHMETIC_OP_LIKELY
                           | expr::BOOLEAN_OP_LIKELY; // Boolean ops are also likely (ex. 1 == 2)

                    // Update the candidate
                    candidate = allocator.alloc();
                    candidate->rule = ast::NUMBER_LITERAL;
                    candidate->value = first.value;
                    expr_stream.next(); // Consume the number literal token
                    break;
                }

                case lexer::token::DECIMAL_LITERAL:
                {
                    // Arithmetic operation is likely
                    likely |= expr::ARITHMETIC_OP_LIKELY
                           | expr::BOOLEAN_OP_LIKELY; // Boolean ops are also likely (ex. 1 == 2)

                    // Update the candidate
                    candidate = allocator.alloc();
                    candidate->rule = ast::DECIMAL_LITERAL;
                    candidate->value = first.value;
                    expr_stream.next(); // Consume the decimal literal token
                    break;
                }

                case lexer::token::STRING_LITERAL:
                {
                    // Boolean operation is likely
                    likely |= expr::BOOLEAN_OP_LIKELY;

                    // Update the candidate
                    candidate = allocator.alloc();
                    candidate->rule = ast::STRING_LITERAL;
                    candidate->value = first.value;
                    expr_stream.next(); // Consume the string literal token
                    break;
                }

                case lexer::token::TRUE:
                {
                    likely |= expr::BOOLEAN_OP_LIKELY; // Boolean operation is likely

                    // Update the candidate
                    candidate = allocator.alloc();
                    candidate->rule = ast::TRUE;
                    expr_stream.next(); // Consume the true token
                    break;
                }

                case lexer::token::FALSE:
                {
                    // Update the candidate
                    candidate = allocator.alloc();
                    candidate->rule = ast::FALSE;
                    expr_stream.next(); // Consume the false token
                    likely |= expr::BOOLEAN_OP_LIKELY; // Boolean operation is likely
                    break;
                }

                case lexer::token::OPEN_PAREN:
                {
                    likely |= expr::PROP_ACCESS_LIKELY // Property access is likely
                            | expr::ARITHMETIC_OP_LIKELY // Arithmetic operation is likely
                            | expr::BOOLEAN_OP_LIKELY; // Boolean ops are also likely (ex. 1 == 2)

                    // Handle nested expression
                    candidate = allocator.alloc();
                    candidate->rule = ast::EXPRESSION;

                    // Extract the nested expression
                    auto nested_expr = extract(expr_stream);
                    expr_queue.emplace_back(container::move(nested_expr), candidate);
                    expr_stream.next(); // Consume the close parenthesis
                    break;
                }

                default:
                {
                    global_err.type = UNEXPECTED_TOKEN;
                    global_err.column = trace.column;
                    global_err.line = trace.line;
                    throw except::exception("Unexpected token in expression");
                    break;
                }
            }

            // Start checking for likely ops
            first_opt = expr_stream.peek();
            if (
                process_next(
                    node,
                    candidate,
                    trace,
                    first_opt,
                    first
                ) // Process the next token
            ) continue;

            first = first_opt.get();
            if (likely & expr::CALL_LIKELY && first.type == lexer::token::OPEN_PAREN)
            {
                candidate = call(
                    candidate,
                    expr_stream,
                    allocator,
                    expr_queue
                ); // Call the function with the candidate

                first_opt = expr_stream.next(); // Peek the next token again
            }

            if (
                process_next(
                    node,
                    candidate,
                    trace,
                    first_opt,
                    first
                ) // Process the next token
            ) continue;

            first = first_opt.get();
            if (likely & expr::PROP_ACCESS_LIKELY && first.type == lexer::token::DOT)
            {
                candidate = prop(
                    candidate,
                    expr_stream,
                    allocator,
                    trace,
                    expr_queue
                ); // Call the property access with the candidate

                first_opt = expr_stream.next(); // Peek the next token again
            }

            if (
                process_next(
                    node,
                    candidate,
                    trace,
                    first_opt,
                    first
                ) // Process the next token
            ) continue;

            first = first_opt.get();
            if (
                likely & expr::ARITHMETIC_OP_LIKELY
                && (
                    first.type == lexer::token::PLUS ||
                    first.type == lexer::token::MINUS ||
                    first.type == lexer::token::MULTIPLY ||
                    first.type == lexer::token::DIVIDE
                )
            )
            {
                candidate = arithmetic(
                    candidate,
                    expr_stream,
                    allocator,
                    expr_queue
                );

                first_opt = expr_stream.curr(); // Get the current token
            }

            if (
                process_next(
                    node,
                    candidate,
                    trace,
                    first_opt,
                    first
                ) // Process the next token
            ) continue;

            first = first_opt.get();
            if (
                likely & expr::BOOLEAN_OP_LIKELY
                && (
                    first.type == lexer::token::BOOL_EQ ||
                    first.type == lexer::token::BOOL_GT ||
                    first.type == lexer::token::BOOL_GTE ||
                    first.type == lexer::token::BOOL_LT ||
                    first.type == lexer::token::BOOL_LTE ||
                    first.type == lexer::token::BOOL_NEQ
                )
            )
            {
                // TODO!
            }

            // Start checking for likely ops
            first_opt = expr_stream.next();
            if (
                process_next(
                    node,
                    candidate,
                    trace,
                    first_opt,
                    first
                ) // Process the next token
            ) continue;

            // Invalid token detected
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = trace.column;
            global_err.line = trace.line;
            throw except::exception("Unexpected token in expression");
        }
    }
}