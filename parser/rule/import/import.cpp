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

#include "import.h"
#include "parser/expect.h"
#include "parser/parser.h"

using namespace zelix;

void parser::rule::imp(
    ast *&root,
    container::stream<lexer::token *> &tokens,
    const bool &top_level,
    memory::lazy_allocator<ast> &allocator,
    const lexer::token *const &trace
)
{
    // Make sure we are at the top level
    if (!top_level)
    {
        global_err.type = ILLEGAL_IMPORT;
        global_err.column = trace->column;
        global_err.line = trace->line;
        throw except::exception("Illegal import statement outside of top-level scope");
    }

    expect(tokens, lexer::token::STRING_LITERAL);
    const auto path = tokens.next()
        .get();

    // Expect a semicolon
    expect(tokens, lexer::token::SEMICOLON);
    tokens.next(); // Consume the semicolon

    // Create the import node
    ast *import_node = allocator.alloc();
    import_node->rule = ast::IMPORT;
    import_node->value = path->value;
    import_node->line = path->line;
    import_node->column = path->column;
    root->children.push_back(import_node);
}
