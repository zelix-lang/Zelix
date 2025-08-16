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
// Created by rodri on 8/16/25.
//

#include "parser/ast.h"
#include "parser/rule/expr/expr.h"
#include "lexer/token.h"
#include "memory/allocator.h"
#include "ret.h"

void zelix::parser::rule::ret(
    container::stream<lexer::token *> &tokens,
    memory::lazy_allocator<ast> &allocator,
    const lexer::token *const &trace,
    ast *&current_block
)
{
    // Allocate a new AST
    auto *tree = allocator.alloc();
    tree->rule = ast::RETURN;

    // Parse the expression
    expression(
        tree,
        tokens,
        allocator,
        trace
    );

    // Add the tree to the current block
    current_block->children.push_back(tree);
}
