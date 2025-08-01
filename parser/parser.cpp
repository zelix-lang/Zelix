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

#include "parser.h"

#include "import/import.h"
using namespace fluent;

parser::ast parse(container::stream<lexer::token> &tokens)
{
    // Create the root AST node
    parser::ast root;

    bool top_level = true; // Flag to track if we are at the top level of the AST

    // Iterate over the tokens
    auto current_opt = tokens.next();
    while (current_opt.is_some())
    {
        switch (
            const auto &current = current_opt.get();
            current.type
        )
        {
            case lexer::token::IMPORT:
            {
                parser::rule::imp(root, tokens, top_level, current);
                break;
            }

            default:
            {
                parser::global_err.type = parser::UNEXPECTED_TOKEN;
                parser::global_err.line = current.line;
                parser::global_err.column = current.column;
                throw except::exception("Unexpected token encountered during parsing");
            }
        }

        current_opt = tokens.next();
    }

    return root;
}