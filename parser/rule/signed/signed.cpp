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
// Created by rodri on 8/16/25.
//

#include "signed.h"
#include "parser/parser.h"
using namespace zelix;

template<bool Arithmetic>
inline parser::ast::rule_t rule(lexer::token *&token)
{
    if constexpr (Arithmetic)
    {
        switch (token->type)
        {
            case lexer::token::PLUS:
                return parser::ast::SUM;
            case lexer::token::MINUS:
                return parser::ast::SUB;
            case lexer::token::MULTIPLY:
                return parser::ast::MUL;
            case lexer::token::DIVIDE:
                return parser::ast::DIV;
            default:
                return parser::ast::ROOT; // Default case if no arithmetic operation matches
        }
    }

    switch (token->type)
    {
        case lexer::token::BOOL_EQ:
            return parser::ast::EQ;
        case lexer::token::BOOL_GT:
            return parser::ast::GT;
        case lexer::token::BOOL_GTE:
            return parser::ast::GTE;
        case lexer::token::BOOL_LT:
            return parser::ast::LT;
        case lexer::token::BOOL_LTE:
            return parser::ast::LTE;
        case lexer::token::BOOL_NEQ:
            return parser::ast::NEQ;
        default:
            return parser::ast::ROOT; // Default case if no arithmetic operation matches
    }
}

template<bool Arithmetic>
inline void process_op(
    lexer::token *&next,
    memory::lazy_allocator<parser::ast> &allocator,
    parser::ast *arithmetic_node
)
{
    // Create a new AST node for the arithmetic operation
    parser::ast *arithmetic_op_node = allocator.alloc();
    arithmetic_op_node->rule = rule<Arithmetic>(next);
    arithmetic_node->children.push_back(arithmetic_op_node); // Add the operation node
}

template<bool Arithmetic>
inline void process_sub(
    container::vector<lexer::token *> &current_tokens,
    lexer::token *&next,
    parser::ast *arithmetic_node,
    parser::ast *candidate,
    memory::lazy_allocator<parser::ast> &allocator,
    memory::lazy_allocator<parser::rule::expr::queue_node> &q_allocator,
    container::vector<parser::rule::expr::queue_node *>
    &expr_queue, bool &first_iteration,
    const bool append_sign = true
)
{
    // Handle first iteration
    if (first_iteration)
    {
        first_iteration = false; // Set the flag to false for subsequent iterations

        if (append_sign)
        {
            process_op<Arithmetic>(next, allocator,
                                   arithmetic_node); // Process the arithmetic operator
        }

        return;
    }

    // Append the current tokens to the last nested node
    if (current_tokens.empty())
    {
        parser::global_err.type = parser::UNEXPECTED_TOKEN;
        parser::global_err.column = next->column;
        parser::global_err.line = next->line;
        throw except::exception("Unexpected empty arithmetic expression");
    }

    // Create a new expression node for the current tokens
    parser::ast *expr_node = allocator.alloc();
    expr_node->rule = parser::ast::EXPRESSION;
    arithmetic_node->children.push_back(expr_node); // Add the expression node

    // Queue the current tokens for processing
    const auto q_el = q_allocator.alloc();
    q_el->tokens = container::stream<lexer::token *>(current_tokens);
    q_el->node = expr_node;
    expr_queue.push_back(q_el);
    current_tokens.clear(); // Clear the current tokens for the next operation

    if (!append_sign)
        return; // If we don't want to append the sign, return early
    process_op<Arithmetic>(next, allocator,
                           arithmetic_node); // Process the arithmetic operator
}

template<bool Arithmetic>
inline void process_high_precedence(
    container::vector<lexer::token *> &current_tokens,
    lexer::token *&next,
    parser::ast *&last_nested,
    parser::ast *&arithmetic_node,
    parser::ast *&candidate,
    memory::lazy_allocator<parser::ast> &allocator,
    memory::lazy_allocator<parser::rule::expr::queue_node> &q_allocator,
    container::vector<parser::rule::expr::queue_node *> &expr_queue,
    bool &first_iteration
)
{
    // Allocate the last nested node if it doesn't exist
    if (last_nested == nullptr)
    {
        last_nested = allocator.alloc();
        last_nested->rule = parser::ast::ARITHMETIC;
        arithmetic_node->children.push_back(last_nested);
    }

    process_sub<Arithmetic>(
        current_tokens,
        next,
        last_nested,
        candidate,
        allocator,
        q_allocator,
        expr_queue,
        first_iteration
    ); // Process the subexpression
}

template<bool Arithmetic>
inline void process_low_precedence(
    container::vector<lexer::token *> &current_tokens,
    lexer::token *&next,
    parser::ast *&last_nested,
    parser::ast *&arithmetic_node,
    parser::ast *&candidate,
    memory::lazy_allocator<parser::ast> &allocator,
    memory::lazy_allocator<parser::rule::expr::queue_node> &q_allocator,
    container::vector<parser::rule::expr::queue_node *> &expr_queue,
    bool &first_iteration
)
{
    // Check if we have a last nested node
    if (last_nested != nullptr)
    {
        // Process the last nested node
        process_sub<Arithmetic>(current_tokens, next, last_nested, candidate, allocator, q_allocator, expr_queue,
                                first_iteration,
                                false); // Process the subexpression
    }
    else
    {
        process_sub<Arithmetic>(current_tokens, next, arithmetic_node, candidate, allocator, q_allocator, expr_queue,
                                first_iteration,
                                false); // Process the subexpression
    }

    // Set the last nested node to nullptr
    last_nested = nullptr; // Reset the last nested node

    process_op<Arithmetic>(next, allocator,
                           arithmetic_node); // Process the arithmetic operator
}

template<bool Arithmetic>
parser::ast *parser::rule::signed_op(
    ast *&candidate,
    container::stream<lexer::token *> &tokens,
    memory::lazy_allocator<ast> &allocator,
    memory::lazy_allocator<expr::queue_node> &q_allocator,
    container::vector<expr::queue_node *> &expr_queue
)
{
    // Create the arithmetic AST node
    ast *arithmetic_node = allocator.alloc();
    arithmetic_node->children.push_back(candidate); // Add the candidate as the first child
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
            if (next->type >= lexer::token::BOOL_EQ && next->type <= lexer::token::BOOL_GTE)
            {
                break;
            }
        }

        // Check if we have an arithmetic operation with high precedence
        bool high_precedence;
        if constexpr (Arithmetic)
        {
            high_precedence = next->type == lexer::token::MULTIPLY || next->type == lexer::token::DIVIDE;
        }
        else
        {
            high_precedence = next->type == lexer::token::OR;
        }

        bool low_precedence;
        if constexpr (Arithmetic)
        {
            low_precedence = next->type == lexer::token::PLUS || next->type == lexer::token::MINUS;
        }
        else
        {
            low_precedence = (next->type >= lexer::token::BOOL_EQ && next->type <= lexer::token::BOOL_GTE) ||
                             next->type == lexer::token::AND;
        }

        if (nested_count == 0 && high_precedence)
        {
            last_is_mul = true;
            process_high_precedence<Arithmetic>(
                current_tokens,
                next,
                last_nested,
                arithmetic_node,
                candidate,
                allocator,
                q_allocator,
                expr_queue,
                first_iteration
            ); // Process the multiplication/division operation
        }
        // Check if we have an arithmetic operator with low precedence
        else if (nested_count == 0 && low_precedence)
        {
            last_is_mul = false;
            process_low_precedence<Arithmetic>(
                current_tokens,
                next,
                last_nested,
                arithmetic_node,
                candidate,
                allocator,
                q_allocator,
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
    const auto q_el = q_allocator.alloc();
    q_el->tokens = container::stream<lexer::token *>(current_tokens); // Clone
    q_el->node = last_sub;
    expr_queue.push_back(q_el);

    // Unwrap if we have a nested expression
    if constexpr (Arithmetic)
    {
        if (arithmetic_node->children.size() == 1)
        {
            const auto first_child = arithmetic_node->children[0];
            allocator.dealloc(arithmetic_node); // Deallocate the arithmetic node
            arithmetic_node = first_child; // Set the first child as the new node
        }
    }

    return arithmetic_node;
}

// signed_op templates
template
parser::ast*
    parser::rule::signed_op<true>(
        ast*&,
        container::stream<lexer::token*>&,
        memory::lazy_allocator<ast, 256ul, false>&,
        memory::lazy_allocator<expr::queue_node, 256ul, false>&,
        container::vector<expr::queue_node*>&
    );

template
parser::ast*
    parser::rule::signed_op<false>(
        ast*&,
        container::stream<lexer::token*>&,
        memory::lazy_allocator<ast, 256ul, false>&,
        memory::lazy_allocator<expr::queue_node, 256ul, false>&,
        container::vector<expr::queue_node*>&
    );