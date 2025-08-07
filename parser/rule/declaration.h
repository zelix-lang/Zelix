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
// Created by rodrigo on 8/4/25.
//

#pragma once

#include "expr/expr.h"
#include "zelix/container/stream.h"
#include "lexer/token.h"
#include "memory/allocator.h"
#include "parser/expect.h"
#include "type.h"

namespace zelix::parser::rule
{
    template <bool Const>
    inline void declaration(
        ast *root,
        container::stream<lexer::token *> &tokens,
        memory::lazy_allocator<ast> &allocator
    )
    {
        // Expect an identifier
        expect(tokens, lexer::token::IDENTIFIER);
        const auto identifier = tokens.next().get();

        // Allocate a new AST node for the declaration
        ast *decl_node = allocator.alloc();
        if constexpr (Const)
        {
            decl_node->rule = ast::CONST_DECLARATION;
        }
        else
        {
            decl_node->rule = ast::DECLARATION;
        }
        root->children.push_back(decl_node); // Add the declaration node as a child of the root

        // Allocate a new node for the identifier
        ast *id_node = allocator.alloc();
        id_node->rule = ast::IDENTIFIER;
        id_node->value = identifier->value; // Set the identifier value
        decl_node->children.push_back(id_node); // Add the identifier node as a child

        // Expect a colon for type declaration
        expect(tokens, lexer::token::COLON);
        tokens.next(); // Consume the colon

        // Parse the type
        type(decl_node, tokens, allocator, identifier);
        expect(tokens, lexer::token::EQUALS); // Expect an equals sign for assignment
        tokens.next(); // Consume the equals sign

        // Queue the rest of the tokens for expression parsing
        // Since expr() extracts until the next semicolon, we don't
        // need to manually expect it here.
        expression(
            decl_node,
            tokens,
            allocator,
            identifier
        );
    }
}