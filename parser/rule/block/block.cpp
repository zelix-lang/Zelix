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

#include "block.h"
#include "parser/expect.h"
#include "parser/rule/assignment.h"
#include "parser/rule/declaration.h"
#include "parser/rule/expr/expr.h"
#include "parser/rule/for.h"
#include "parser/rule/if.h"
#include "parser/rule/ret.h"

using namespace zelix;

void parser::rule::block(
    ast *&root,
    container::stream<lexer::token *> &tokens,
    memory::lazy_allocator<ast> &allocator,
    const lexer::token *const &trace
)
{
    expect(tokens, lexer::token::OPEN_CURLY);
    tokens.next(); // Consume the open curly brace

    // Create a new block AST node
    ast *root_block = allocator.alloc();
    root_block->rule = ast::BLOCK;
    root->children.push_back(root_block);

    // Use a queue for the children of the block
    container::vector<ast *> block_queue;
    block_queue.push_back(root_block);

    // Save the current conditional for if/else blocks
    ast *current_conditional = nullptr;

    // Get the next token
    while (true)
    {
        auto next_opt = tokens.peek();
        if (next_opt.is_none() && !block_queue.empty())
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.column = trace->column;
            global_err.line = trace->line;
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
                if (block_queue.empty())
                {
                    global_err.type = UNEXPECTED_TOKEN;
                    global_err.column = next->column;
                    global_err.line = next->line;
                    throw except::exception("Unexpected close curly brace without an open block");
                }

                block_queue.pop_back(); // Close the current block
                tokens.next(); // Consume the close curly brace
                if (block_queue.empty())
                {
                    return; // If the block queue is empty, we are done parsing the block
                }

                break;
            }

            case lexer::token::OPEN_CURLY:
            {
                // Create a new block node for the nested block
                ast *nested_block = allocator.alloc();
                nested_block->rule = ast::BLOCK;
                current_block->children.push_back(nested_block);
                tokens.next(); // Consume the open curly brace

                block_queue.push_back(nested_block);
                break;
            }

            // Handle variables
            case lexer::token::LET:
            {
                tokens.next(); // Consume the let token

                declaration<false>(
                    current_block,
                    tokens,
                    allocator
                );

                break;
            }

            case lexer::token::CONST:
            {
                tokens.next(); // Consume the const token

                declaration<true>(
                    current_block,
                    tokens,
                    allocator
                );

                break;
            }

            case lexer::token::IF:
            {
                tokens.next(); // Consume the token

                conditional<true, false, false, false>(
                    current_block,
                    current_conditional,
                    next,
                    tokens,
                    allocator
                );

                break;
            }

            case lexer::token::ELSEIF:
            {
                tokens.next(); // Consume the token

                conditional<false, true, false, false>(
                    current_block,
                    current_conditional,
                    next,
                    tokens,
                    allocator
                );

                break;
            }

            case lexer::token::ELSE:
            {
                tokens.next(); // Consume the token

                conditional<false, false, true, false>(
                    current_block,
                    current_conditional,
                    next,
                    tokens,
                    allocator
                );

                break;
            }

            case lexer::token::WHILE:
            {
                tokens.next(); // Consume the token

                conditional<false, false, false, true>(
                    current_block,
                    current_conditional,
                    next,
                    tokens,
                    allocator
                );
                break;
            }

            case lexer::token::FOR:
            {
                tokens.next(); // Consume the token

                for_loop(
                    current_block,
                    tokens,
                    allocator,
                    next
                );

                break;
            }

            case lexer::token::RETURN:
            {
                tokens.next(); // Consume the token
                ret(
                    tokens,
                    allocator,
                    next,
                    current_block
                );

                break;
            }

            case lexer::token::IDENTIFIER:
            {
                if (
                    assignment(
                        current_block,
                        tokens,
                        allocator,
                        next
                    )
                ) break; // Successfully parsed an assignment
            }

            default:
            {
                // Pass the expression to the expression parser
                expression(
                    current_block,
                    tokens,
                    allocator,
                    trace
                );
            }
        }
    }
}

