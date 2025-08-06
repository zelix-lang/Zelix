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
#include "parser/parser.h"

namespace zelix::parser::rule
{
    inline void type(
        ast *&root,
        container::stream<lexer::token *> &tokens,
        memory::lazy_allocator<ast> &allocator,
        lexer::token *const& trace
    )
    {
        // Get the next token
        const auto next_opt = tokens.next();
        if (next_opt.is_none())
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = trace->column;
            global_err.line = trace->line;
            throw except::exception("Unexpected end of input while parsing type");
        }

        // Create a new AST node for the type
        ast *type_node = allocator.alloc();
        type_node->rule = ast::TYPE;

        auto &next = next_opt.get();
        ast *node = allocator.alloc();
        switch (next->type)
        {
            case lexer::token::NOTHING:
            {
                node->rule = ast::NOTHING;
                break;
            }

            case lexer::token::STRING:
            {
                node->rule = ast::STR;
                break;
            }

            case lexer::token::NUMBER:
            {
                node->rule = ast::NUM;
                break;
            }

            case lexer::token::DECIMAL:
            {
                node->rule = ast::DEC;
                break;
            }

            case lexer::token::IDENTIFIER:
            {
                node->rule = ast::IDENTIFIER;
                break;
            }

            default:
            {
                global_err.type = UNEXPECTED_TOKEN;
                global_err.column = next->column;
                global_err.line = next->line;
                throw except::exception("Unexpected token while parsing type");
            }
        }

        type_node->children.push_back(node);
        root->children.push_back(container::move(type_node));
    }
}