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
// Created by rodrigo on 8/1/25.
//

#pragma once

namespace zelix::parser
{
    container::stream<lexer::token *> extract(
        container::stream<lexer::token*> &tokens,
        lexer::token::t_type end_delim = lexer::token::CLOSE_PAREN,
        lexer::token::t_type nested_end_delim = lexer::token::CLOSE_PAREN,
        lexer::token::t_type start_delim = lexer::token::OPEN_PAREN,
        bool handle_nested = true,
        bool exclude_first_delim = true
    );
}