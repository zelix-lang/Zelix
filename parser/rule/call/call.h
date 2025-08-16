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
// Created by rodrigo on 8/2/25.
//

#pragma once

#include "../args/args.h"
#include "lexer/token.h"
#include "memory/allocator.h"
#include "parser/parser.h"
#include "parser/rule/expr/queue.h"
#include "zelix/container/stream.h"

namespace zelix::parser::rule
{
    inline ast *call(
        ast *&candidate,
        container::stream<lexer::token *> &tokens,
        memory::lazy_allocator<ast> &allocator,
        memory::lazy_allocator<expr::queue_node> &q_allocator,
        container::vector<expr::queue_node *> &expr_queue
    )
    {
        // Get the vector under the tokens
        auto &vec = tokens.ptr();

        // Make sure we have at least 2 tokens
        if (
            const auto pos = tokens.pos();
            vec.size() <= pos + 1
        )
        {
            const auto &trace = vec.ref_at(pos);
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = trace->column;
            global_err.line = trace->line;
            throw except::exception("Not enough tokens to form a call expression");
        }

        // Create a new AST node for the call
        ast *call_node = allocator.alloc();
        call_node->rule = ast::CALL;
        call_node->line = candidate->line;
        call_node->column = candidate->column;
        call_node->children.push_back(candidate); // Push the candidate as the first child (function name)

        // Parse the arguments
        args(call_node, tokens, allocator, q_allocator, expr_queue);

        return call_node; // Return the call node to update the candidate
    }
}