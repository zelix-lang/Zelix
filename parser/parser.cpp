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

#include "parser.h"

#include "memory/allocator.h"
#include "rule/function.h"
#include "rule/import.h"
#include "rule/mod.h"
#include "rule/package.h"
using namespace zelix;

parser::ast *parser::parse(
    container::stream<lexer::token *> &tokens,
    memory::lazy_allocator<ast> &allocator
)
{
    // Create the root AST node
    ast *root = allocator.alloc();

    bool top_level = true; // Flag to track if we are at the top level of the AST
    bool pub = false; // Flag to track if the next declaration is public

    // Parse the package
    rule::package(root, tokens, allocator);

    // Iterate over the tokens
    auto current_opt = tokens.next();
    while (current_opt.is_some())
    {
        switch (
            const auto &current = current_opt.get();
            current->type
        )
        {
            case lexer::token::PUB:
            {
                pub = true; // Next declaration is public
                break;
            }

            case lexer::token::IMPORT:
            {
                if (pub)
                {
                    global_err.type = UNEXPECTED_TOKEN;
                    global_err.line = current->line;
                    global_err.column = current->column;
                    throw except::exception("The 'pub' modifier cannot be applied to import statements");
                }

                rule::imp(root, tokens, top_level, allocator, current);
                break;
            }

            case lexer::token::FUNCTION:
            {
                top_level = false; // We are no longer at the top level after a function declaration
                rule::function(root, tokens, allocator, current, pub);
                break;
            }

            case lexer::token::MOD:
            {
                top_level = false; // We are no longer at the top level after a mod declaration
                rule::mod(root, tokens, allocator, current, pub);
                break;
            }

            default:
            {
                global_err.type = UNEXPECTED_TOKEN;
                global_err.line = current->line;
                global_err.column = current->column;
                throw except::exception("Unexpected token encountered during parsing");
            }
        }

        current_opt = tokens.next();
    }

    // Pub must be false at the end of parsing
    if (pub)
    {
        global_err.type = UNEXPECTED_TOKEN;
        global_err.line = 0; // No specific line for the error
        global_err.column = 0; // No specific column for the error
        throw except::exception("The 'pub' modifier must be followed by a declaration");
    }

    return root;
}