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
// Created by rodri on 8/14/25.
//

#pragma once

#include "memory/allocator.h"
#include "parser/ast.h"
#include "parser/expect.h"
#include "zelix/container/stream.h"

namespace zelix::parser::rule
{
    inline void package(
        ast *&root,
        container::stream<lexer::token *> &tokens,
        memory::lazy_allocator<ast> &allocator
    )
    {
        // Expect the very first token to be a package
        expect(tokens, lexer::token::PACKAGE);
        const auto &trace = tokens.next().get(); // Consume the package token
        bool id = true; // Flag to track if we are expecting an identifier

        // Allocate a new AST node for the package
        ast *package_node = allocator.alloc();
        package_node->rule = ast::PACKAGE;
        root->children.push_back(package_node); // Add the package node to the root

        // Get the next token
        auto next_opt = tokens.next();
        if (next_opt.is_none())
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = trace->column;
            global_err.line = trace->line;
            throw except::exception("Unexpected end of input while parsing package");
        }

        while (next_opt.is_some())
        {
            const auto &next = next_opt.get();

            // Break if we encounter a semicolon
            if (next->type == lexer::token::SEMICOLON)
            {
                break;
            }

            if (id)
            {
                if (next->type != lexer::token::IDENTIFIER)
                {
                    global_err.type = UNEXPECTED_TOKEN;
                    global_err.column = next->column;
                    global_err.line = next->line;
                    throw except::exception("Unexpected end of input while parsing package");
                }

                // Allocate a new AST node for the identifier
                ast *id_node = allocator.alloc();
                id_node->rule = ast::IDENTIFIER;
                id_node->value = next->value; // Set the identifier value
                id_node->line = next->line; // Set the line number
                id_node->column = next->column; // Set the column number
                package_node->children.push_back(id_node); // Add the identifier node to the package node
                id = false; // Switch to expecting a dot next
                next_opt = tokens.next(); // Get the next token
                continue;
            }

            if (next->type != lexer::token::DOT)
            {
                global_err.type = UNEXPECTED_TOKEN;
                global_err.column = next->column;
                global_err.line = next->line;
                throw except::exception("Expected '.' after package identifier");
            }

            id = true; // Switch to expecting an identifier next
            next_opt = tokens.next(); // Get the next token
        }

        // Make sure we don't end up expecting an identifier
        if (id)
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = trace->column;
            global_err.line = trace->line;
            throw except::exception("Unexpected end of input while parsing package");
        }
    }
}