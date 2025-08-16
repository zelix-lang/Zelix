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

#pragma once
#include "file/program.h"

namespace zelix::code::rule
{
    inline void package(parser::ast *ast, program &pro)
    {
        container::string str;
        str.reserve(ast->children.size() * 2 - 1);

        // Build the string
        const auto size = ast->children.size() - 1;
        for (size_t i = 0; i <= size; i++)
        {
            const auto &child = ast->children[i];
            const auto &value = child->value.get();
            str.push(value.ptr(), value.size());

            if (i != size)
            {
                str.push('.');
            }
        }

        // Insert the new package
        pro.new_pkg(str);
    }
}