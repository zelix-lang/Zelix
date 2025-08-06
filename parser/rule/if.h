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
// Created by rodrigo on 8/5/25.
//

#pragma once
#include "expr/expr.h"
#include "fluent/container/stream.h"
#include "lexer/token.h"
#include "memory/allocator.h"
#include "parser/parser.h"

namespace zelix::parser::rule
{
    template <bool If, bool ElseIf, bool Else, bool While>
    inline void conditional(
        ast *&root,
        ast *&current_conditional,
        const lexer::token *const &trace,
        container::stream<lexer::token *> &tokens,
        memory::lazy_allocator<ast> &allocator
    )
    {
        // Allocate a new AST node for the conditional
        ast *cond_node = allocator.alloc();
        if constexpr (If)
        {
            cond_node->rule = ast::IF;
        }
        else if constexpr (ElseIf)
        {
            cond_node->rule = ast::ELSEIF;
        }
        else if constexpr (Else)
        {
            cond_node->rule = ast::ELSE;
        }
        else if constexpr (While)
        {
            cond_node->rule = ast::WHILE;
        }

        if constexpr (!If)
        {
            // Make sure we have a current conditional
            if (current_conditional == nullptr)
            {
                global_err.type = UNEXPECTED_TOKEN;
                global_err.column = trace->column;
                global_err.line = trace->line;
                throw except::exception("Unexpected else without a preceding if");
            }
        }

        // Parse the condition
        if constexpr (If || ElseIf || While)
        {
            expression<true, false>(
                cond_node,
                tokens,
                allocator,
                trace
            );
        }
        else
        {
            // Expect an open curly brace for else
            expect(tokens, lexer::token::OPEN_CURLY);
            tokens.next(); // Consume the open curly brace
        }

        // Append the conditional node to the root
        root->children.push_back(cond_node);

        // Set the current conditional to the new node
        if constexpr ((If || ElseIf) && !While)
        {
            current_conditional = cond_node; // Set the current conditional to the new node
        }
        else
        {
            current_conditional = nullptr; // Reset the current conditional for else
        }
    }
}