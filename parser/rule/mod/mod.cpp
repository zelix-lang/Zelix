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

#include "mod.h"
#include "../declaration/declaration.h"
#include "../derive/derive.h"
#include "../function/function.h"
#include "parser/expect.h"
#include "parser/parser.h"

using namespace zelix;

void parser::rule::mod(
    ast *&root,
    container::stream<lexer::token *> &tokens,
    memory::lazy_allocator<ast> &allocator,
    const lexer::token *const &trace,
    bool &pub
)
{
        expect(tokens, lexer::token::IDENTIFIER); // Expect an identifier for the module name
        const auto name = tokens.next()
            .get();
        expect(tokens, lexer::token::OPEN_CURLY); // Expect an open curly brace
        tokens.next(); // Consume the open curly brace

        // Create a new AST node for the module
        ast *module = allocator.alloc();
        module->rule = ast::MOD;

        // Honor the public flag
        if (pub)
        {
            ast *public_ast = allocator.alloc();
            public_ast->rule = ast::PUBLIC;
            module->children.push_back(public_ast);
        }

        ast *name_ast = allocator.alloc();
        name_ast->rule = ast::IDENTIFIER;
        name_ast->line = name->line;
        name_ast->column = name->column;
        name_ast->value = name->value;
        module->children.push_back(name_ast);

        // Start parsing the module block
        pub = false; // Reset the pub flag for the module
        auto next_opt = tokens.next();
        bool expecting_declaration = false;

        while (next_opt.is_some())
        {
            const auto &next = next_opt.get();

            // Check if we have to break
            if (next->type == lexer::token::CLOSE_CURLY)
            {
                // Nested braces are handled by individual parsers
                break; // Exit the loop when we reach the end of the module block
            }

            switch (next->type)
            {
                case lexer::token::PUB:
                {
                    pub = true;
                    break; // Skip to the next iteration
                }

                case lexer::token::FUNCTION:
                {
                    // Parse the function declaration
                    function(module, tokens, allocator, next, pub);
                    break;
                }

                case lexer::token::LET:
                {
                    declaration<false>(
                        module,
                        tokens,
                        allocator
                    );
                    expecting_declaration = false;
                    break;
                }

                case lexer::token::CONST:
                {
                    declaration<true>(
                        module,
                        tokens,
                        allocator
                    );

                    expecting_declaration = false;
                    break;
                }

                case lexer::token::DERIVE:
                {
                    derive(
                        module,
                        tokens,
                        allocator
                    );
                    expecting_declaration = true; // Next we expect a declaration
                    break;
                }

                default:
                {
                    // Unexpected token in module block
                    global_err.type = UNEXPECTED_TOKEN;
                    global_err.line = next->line;
                    global_err.column = next->column;
                    throw except::exception("Unexpected token in module block");
                }
            }

            next_opt = tokens.next(); // Get the next token
        }

        // Make sure we don't expect a declaration
        if (expecting_declaration)
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.line = trace->line;
            global_err.column = trace->column;
            throw except::exception("Expected a declaration after 'derive'");
        }

        // pub must never be true at this point
        if (pub)
        {
            global_err.type = UNEXPECTED_TOKEN;
            global_err.line = trace->line;
            global_err.column = trace->column;
            throw except::exception("The 'pub' modifier cannot be applied here");
        }

        // Append the module to the root
        root->children.push_back(module);
    }
