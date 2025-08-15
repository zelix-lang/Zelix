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
// Created by rodri on 8/11/25.
//

#include "type.h"

#include "parser/rule/package.h"

using namespace zelix;

parser::ast *parse_base(
    container::stream<lexer::token *> &tokens,
    memory::lazy_allocator<parser::ast> &allocator,
    lexer::token *const &trace,
    bool &allow_nested,
    parser::ast *&type_node
)
{
    // Get the next token
    auto next_opt = tokens.peek();
    if (next_opt.is_none())
    {
        parser::global_err.type = parser::UNEXPECTED_TOKEN;
        parser::global_err.column = trace->column;
        parser::global_err.line = trace->line;
        throw except::exception("Unexpected end of input while parsing type");
    }

    auto &next = next_opt.get();

    // Parse pointers
    while (
        next->type == lexer::token::AMPERSAND || // Pointer type
        next->type == lexer::token::AND // Double pointer (&&)
    )
    {
        parser::ast *pointer_node = allocator.alloc();
        pointer_node->rule = parser::ast::PTR;
        pointer_node->line = next->line;
        pointer_node->column = next->column;
        type_node->children.push_back(pointer_node);

        if (next->type == lexer::token::AND)
        {
            // Add the same node for a double pointer (&&)
            type_node->children.push_back(pointer_node);
        }

        // Get the next token
        next_opt = tokens.next();
        if (next_opt.is_none())
        {
            parser::global_err.type = parser::UNEXPECTED_TOKEN;
            parser::global_err.column = trace->column;
            parser::global_err.line = trace->line;
            throw except::exception("Unexpected end of input while parsing type");
        }
        next = next_opt.get();
    }

    parser::ast *node = allocator.alloc();
    node->line = next->line;
    node->column = next->column;
    switch (next->type)
    {
        case lexer::token::NOTHING:
        {
            node->rule = parser::ast::NOTHING;
            allow_nested = false; // No nested types allowed after NOTHING
            tokens.next(); // Consume the token
            break;
        }

        case lexer::token::STRING:
        {
            node->rule = parser::ast::STR;
            allow_nested = false; // No nested types allowed
            tokens.next(); // Consume the token
            break;
        }

        case lexer::token::NUMBER:
        {
            node->rule = parser::ast::NUM;
            allow_nested = false; // No nested types allowed
            tokens.next(); // Consume the token
            break;
        }

        case lexer::token::DECIMAL:
        {
            node->rule = parser::ast::DEC;
            allow_nested = false; // No nested types allowed
            tokens.next(); // Consume the token
            break;
        }

        case lexer::token::BOOL:
        {
            node->rule = parser::ast::BOOL;
            allow_nested = false; // No nested types allowed
            tokens.next(); // Consume the token
            break;
        }

        default:
        {
            // Parse a package
            parser::rule::package<false, lexer::token::UNKNOWN, true>(
                type_node,
                tokens,
                allocator
            );

            allow_nested = true;
            allocator.dealloc(node);
            return type_node; // Return the type node without the base type
        }
    }

    type_node->children.push_back(node);
    return type_node;
}

void parser::rule::type(
    ast *&root,
    container::stream<lexer::token *> &tokens,
    memory::lazy_allocator<ast> &allocator,
    lexer::token *const &trace
)
{
    // Parse the base type
    bool allow_nested = true; // Flag to allow nested types
    // Create a new AST node for the type
    ast *type_node = allocator.alloc();
    type_node->rule = ast::TYPE;
    ast *base = parse_base(tokens, allocator, trace, allow_nested, type_node);
    bool parsed_base = true; // Flag to indicate if we successfully parsed a base type
    container::vector<ast *> children; // Vector to hold the children of the base type

    auto peek_opt = tokens.peek();
    if (!allow_nested || peek_opt.is_none())
    {
        children.push_back(base);
    }
    else
    {
        auto nested = allocator.alloc();
        nested->rule = ast::TYPE;
        base->children.push_back(nested); // Add the nested type to the base type
        children.push_back(nested);
    }

    // Peek into the next token
    while (peek_opt.is_some())
    {
        const auto &curr = peek_opt.get();

        // Parse nested children
        if (curr->type == lexer::token::BOOL_LT)
        {
            if (!parsed_base || !allow_nested)
            {
                global_err.type = UNEXPECTED_TOKEN;
                global_err.column = curr->column;
                global_err.line = curr->line;
                throw except::exception("Unexpected '<' token in type declaration");
            }

            parsed_base = false;
            tokens.next(); // Consume the '<' token
            const auto &back = children.back(); // Get the last base type
            auto nested = allocator.alloc();
            nested->rule = ast::TYPE;
            nested->line = curr->line;
            nested->column = curr->column;
            back->children.push_back(nested);
            children.push_back(nested); // Add the nested type to the children vector
            peek_opt = tokens.peek(); // Peek the next token
            continue;
        }

        if (curr->type == lexer::token::BOOL_GT)
        {
            if (children.empty())
            {
                global_err.type = UNEXPECTED_TOKEN;
                global_err.column = curr->column;
                global_err.line = curr->line;
                throw except::exception("Unexpected '>' token in type declaration");
            }

            if (!parsed_base)
            {
                global_err.type = UNEXPECTED_TOKEN;
                global_err.column = curr->column;
                global_err.line = curr->line;
                throw except::exception("Unexpected '>' token in type declaration");
            }

            // Pop the last nested type
            children.pop_back();
            tokens.next(); // Consume the '>' token
            peek_opt = tokens.peek(); // Peek the next token
            continue;
        }

        if (curr->type == lexer::token::COMMA)
        {
            if (!parsed_base)
            {
                global_err.type = UNEXPECTED_TOKEN;
                global_err.column = curr->column;
                global_err.line = curr->line;
                throw except::exception("Unexpected ',' token in type declaration");
            }

            if (children.empty())
            {
                global_err.type = UNEXPECTED_TOKEN;
                global_err.column = curr->column;
                global_err.line = curr->line;
                throw except::exception("Unexpected ',' token in type declaration");
            }

            // Consume the comma
            tokens.next();

            // Go back to the container TYPE (just inside the < >)
            const auto container_type = children[children.size() - 2];
            // Create a new TYPE node for the next argument
            auto arg_type = allocator.alloc();
            arg_type->rule = ast::TYPE;
            container_type->children.push_back(arg_type);

            // Replace the last child in `children` with this new arg type
            children.back() = arg_type;

            parsed_base = false;
            allow_nested = true; // Let the new arg type have nested types
            peek_opt = tokens.peek();
            continue;
        }

        if (!parsed_base)
        {
            // Get the base type from the children vector
            auto &last_child = children.back();
            parse_base(tokens, allocator, trace, allow_nested, last_child);
            parsed_base = true; // Set the parsed base flag to true
            peek_opt = tokens.peek(); // Peek the next token
            continue;
        }

        // If we reach here, it means we have a valid base type
        break; // Exit the loop
    }

    // Make sure there are no unclosed nested types
    if (children.size() > 1) // 1 for the base type
    {
        global_err.type = UNEXPECTED_TOKEN;
        global_err.column = children.back()->column;
        global_err.line = children.back()->line;
        throw except::exception("Unclosed nested type in type declaration");
    }

    // Check for empty types at the end
    auto &back = base->children.back();
    while (base->children.size() > 1 && back->children.empty())
    {
        allocator.dealloc(back);
        base->children.pop_back();
        back = base->children.back();
    }

    root->children.push_back(container::move(base));
}