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

#include "expr/expr.h"
#include "zelix/container/stream.h"
#include "lexer/token.h"
#include "memory/allocator.h"

namespace zelix::parser::rule
{
    inline bool assignment(
        ast *&root,
        container::stream<lexer::token *> &tokens,
        memory::lazy_allocator<ast> &allocator,
        const lexer::token *const &trace
    )
    {
        // Peek into the next token to determine if it's an assignment
        const auto next_opt = tokens.peek();
        if (next_opt.is_none())
        {
            return false; // No more tokens to process
        }

        if (
            const auto &next = next_opt.get();
            next->type == lexer::token::EQUALS
        )
        {
            return false; // Not an assignment
        }

        tokens.next(); // Consume the '=' token

        // Create an assignment AST node
        ast *assign_node = allocator.alloc();
        assign_node->rule = ast::ASSIGNMENT;

        // Allocate a new node for the identifier
        ast *id_node = allocator.alloc();
        id_node->rule = ast::IDENTIFIER;
        id_node->value = trace->value; // Set the identifier value
        assign_node->children.push_back(id_node); // Add the identifier node as a child

        // Parse the expression on the right side of the assignment
        expression(assign_node, tokens, allocator, trace);
        root->children.push_back(assign_node); // Add the assignment node as a child of the root
        return true; // Successfully parsed an assignment
    }
}