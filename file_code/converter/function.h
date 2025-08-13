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
// Created by rodri on 8/11/25.
//

#pragma once
#include "file_code/file_code.h"
#include "file_code/function.h"
#include "memory/allocator.h"
#include "type/type.h"

namespace zelix::code::converter
{
    inline void function(
        ankerl::unordered_dense::map<
            container::external_string,
            function *,
            container::external_string_hash
        > &map,
        const parser::ast *root,
        memory::lazy_allocator<function> &fun_allocator
    )
    {
        // Allocate a new function
        code::function *fun = fun_allocator.alloc();

        // Get the children
        const auto &children = root->children;

        // Honor visibility
        auto &first = children[0];
        auto &second = children[1];
        if (first->rule == parser::ast::PUBLIC)
        {
            fun->pub = true;
            first = children[1]; // Move to the next child
            second = children[2]; // Move to the next child
        }
        else
        {
            fun->pub = false;
        }

        // Get the function name
        const auto &name = first->value.get();

        // Honor arguments too
        if (second->rule == parser::ast::ARGUMENTS)
        {
            // Convert arguments
            for (const auto &arg: second->children)
            {
                // Get the argument's children
                const auto &arg_children = arg->children;

                // Get the argument's info
                const auto &arg_name = arg_children[0]->value.get();
                const auto &arg_type = arg_children[1];

                // Convert the argument type
                auto wrapped_type = type(arg_type);

                // Append the argument
                fun->args[arg_name] = container::move(wrapped_type);
            }

            if (fun->pub)
            {
                second = children[3]; // Move to the next child
            }
            else
            {
                second = children[2]; // Move to the next child
            }
        }

        // Set the function's block
        fun->body = second;

        // Insert the function into the map
        map[name] = fun;
    }
}