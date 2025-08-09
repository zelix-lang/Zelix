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
// Created by rodri on 8/8/25.
//

#pragma once

#include "lexer/token.h"
#include "memory/allocator.h"
#include "parser/ast.h"
#include "parser/parser.h"
#include "zelix/container/stream.h"

namespace zelix::parser::rule
{
    inline void derive(
        ast *&root,
        container::stream<lexer::token *> &tokens,
        memory::lazy_allocator<ast> &allocator
    )
    {
        // Get the next token
        auto next_opt = tokens.next();
        if (next_opt.is_none())
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = 0;
            global_err.line = 0;
            throw except::exception("Unexpected end of input while parsing derive");
        }

        // Create a new AST node
        ast *derive_node = allocator.alloc();
        derive_node->rule = ast::DERIVE;

        auto &next = next_opt.get();
        bool expecting_comma = false;

        // Parse the derive list
        while (next->type != lexer::token::SEMICOLON)
        {
            if (expecting_comma)
            {
                if (next->type != lexer::token::COMMA)
                {
                    global_err.type = UNEXPECTED_TOKEN;
                    global_err.column = next->column;
                    global_err.line = next->line;
                    throw except::exception("Expected comma in derive list");
                }

                expecting_comma = false;
            }
            else
            {
                // Make sure the token is an identifier
                if (next->type != lexer::token::IDENTIFIER)
                {
                    global_err.type = UNEXPECTED_TOKEN;
                    global_err.column = next->column;
                    global_err.line = next->line;
                    throw except::exception("Expected identifier in derive list");
                }

                expecting_comma = true; // Next we expect a comma or semicolon
            }

            next_opt = tokens.next();
            if (next_opt.is_none())
            {
                global_err.type = UNEXPECTED_TOKEN;
                global_err.column = next->column;
                global_err.line = next->line;
                throw except::exception("Unexpected end of input while parsing derive");
            }

            next = next_opt.get();
        }

        // Make sure we ended correctly
        if (!expecting_comma)
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = next->column;
            global_err.line = next->line;
            throw except::exception("Trailing comma in derive list");
        }

        root->children.push_back(derive_node);
    }
}
