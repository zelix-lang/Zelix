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

#include "block.h"
#include "parser/expect.h"

using namespace fluent;

void parser::rule::block(
    ast *&root,
    container::stream<lexer::token*> &tokens,
    memory::lazy_allocator<ast> &allocator,
    const lexer::token &trace
)
{
    expect(tokens, lexer::token::OPEN_CURLY);
    tokens.next(); // Consume the open curly brace

    // Create a new block AST node
    ast *root_block = allocator.alloc();
    root_block->rule = ast::BLOCK;

    // Use a queue for the children of the block
    container::vector<ast *> block_queue;
    block_queue.push_back(root_block);

    // Get the next token
    auto next_opt = tokens.next();
    while (next_opt.is_some())
    {
        if (block_queue.empty())
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = trace.column;
            global_err.line = trace.line;
            throw except::exception("Unexpected end of block");
        }

        // Get the current block node
        ast *current_block = block_queue.ref_at(block_queue.size() - 1);

        switch (
            const auto &next = next_opt.get();
            next->type
        )
        {
            case lexer::token::CLOSE_CURLY:
            {
                block_queue.pop_back(); // Close the current block
                break;
            }

            case lexer::token::OPEN_CURLY:
            {
                // Create a new block node for the nested block
                ast *nested_block = allocator.alloc();
                nested_block->rule = ast::BLOCK;
                current_block->children.push_back(nested_block);

                block_queue.push_back(nested_block);
                break;
            }

            default:
            {

            }
        }
    }
}

