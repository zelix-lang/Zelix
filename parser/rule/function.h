/*
        ==== The Zelix Programming Language ====
---------------------------------------------------------
  - This file is part of the Zelix Programming Language
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
#include "block/block.h"
#include "zelix/container/stream.h"
#include "lexer/token.h"
#include "memory/allocator.h"
#include "parser/expect.h"
#include "parser/parser.h"
#include "parser/rule/type.h"

namespace zelix::parser::rule
{
    inline void function(
        ast *&root,
        container::stream<lexer::token *> &tokens,
        memory::lazy_allocator<ast> &allocator,
        const lexer::token *const &trace
    )
    {
        expect(tokens, lexer::token::IDENTIFIER);
        const auto name = tokens.next()
            .get();

        expect(tokens, lexer::token::OPEN_PAREN);
        tokens.next(); // Consume the open parenthesis
        ast *function = allocator.alloc(); // Create a new function AST node
        ast *name_ast = allocator.alloc();;
        name_ast->rule = ast::IDENTIFIER;
        name_ast->value = name->value;
        function->children.push_back(name_ast);

        // Check if we don't have any args
        auto peek = tokens.peek();
        if (peek.is_some() && peek.get()->type == lexer::token::CLOSE_PAREN)
        {
            tokens.next(); // Consume the close parenthesis
        }
        else
        {
            ast *args_node = allocator.alloc();
            args_node->rule = ast::ARGUMENTS;

            // Parse the arguments
            while (true)
            {
                // Expect the identifier for the argument
                expect(tokens, lexer::token::IDENTIFIER);
                const auto arg = tokens.next()
                    .get();
                expect(tokens, lexer::token::COLON); // Expect the colon after the identifier
                tokens.next(); // Consume the colon

                // Create a name AST
                ast *arg_name = allocator.alloc();
                arg_name->rule = ast::IDENTIFIER;
                arg_name->value = arg->value;

                ast *arg_node = allocator.alloc();
                arg_node->rule = ast::ARGUMENT;
                arg_node->children.push_back(arg_name);

                // Parse the type
                type(arg_node, tokens, allocator, arg);
                args_node->children.push_back(arg_node);

                // Consume the next token
                auto next_opt = tokens.next();
                if (next_opt.is_none())
                {
                    global_err.type = UNEXPECTED_TOKEN;
                    global_err.column = arg->column;
                    global_err.line = arg->line;
                }

                // Get the next token
                const auto &next = next_opt.get();

                // Check for more params
                if (next->type == lexer::token::COMMA)
                {
                    continue;
                }

                // Check if we have reached the end
                if (next->type == lexer::token::CLOSE_PAREN)
                {
                    break;
                }

                global_err.type = UNEXPECTED_TOKEN;
                global_err.column = arg->column;
                global_err.line = arg->line;
                throw except::exception("Invalid function signature");
            }

            function->children.push_back(args_node);
        }

        // Get the next token
        auto peek_opt = tokens.next();
        if (peek_opt.is_none())
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = trace->column;
            global_err.line = trace->line;
            throw except::exception("Invalid function signature");
        }

        // Check if we have a return type
        if (
            const auto type_peek = peek_opt.get();
            type_peek->type == lexer::token::ARROW
        )
        {
            // Consume the arrow token
            tokens.next();

            // Parse the return type
            type(function, tokens, allocator, type_peek);
        }
        else if (type_peek->type != lexer::token::OPEN_CURLY)
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = type_peek->column;
            global_err.line = type_peek->line;
            throw except::exception("Invalid function signature");
        }

        // Parse the block
        block(function, tokens, allocator, trace);
        root->children.push_back(function);
    }
}