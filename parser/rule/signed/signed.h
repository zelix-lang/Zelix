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

#include "parser/ast.h"
#include "zelix/container/stream.h"
#include "lexer/token.h"
#include "memory/allocator.h"
#include "parser/rule/expr/queue.h"

namespace zelix::parser::rule
{
    template <bool Arithmetic = true>
    ast *signed_op(
        ast *&candidate,
        container::stream<lexer::token *> &tokens,
        memory::lazy_allocator<ast> &allocator,
        memory::lazy_allocator<expr::queue_node> &q_allocator,
        container::vector<expr::queue_node *> &expr_queue
    );
} // namespace zelix::parser::rule