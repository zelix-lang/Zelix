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
// Created by rodri on 8/12/25.
//

#pragma once
#include "file_code/file_code.h"
#include "file_code/function.h"
#include "function.h"
#include "memory/allocator.h"
#include "type/type.h"

namespace zelix::code::converter
{
    inline void mod(
        ankerl::unordered_dense::map<
            container::external_string,
            mod *,
            container::external_string_hash
        > &map,
        const parser::ast *root,
        memory::lazy_allocator<code::function> &fun_allocator,
        memory::lazy_allocator<mod> &mod_allocator
    )
    {
        // Allocate a new mod
        code::mod *m = mod_allocator.alloc();

        // Get the children
        const auto &children = root->children;

        // Honor visibility
        size_t start = 0;
        auto &first = children[start];
        if (first->rule == parser::ast::PUBLIC)
        {
            m->pub = true;
            first = children[1]; // Move to the next child
            start += 2; // Skip the public keyword and the name
        }
        else
        {
            m->pub = false;
            start++; // Skip the name
        }

        // Get the module name
        const auto &name = first->value.get();
        parser::ast *derive = nullptr; // Save the last derive node

        // Iterate over the children and fill the mod
        for (size_t i = start; i < children.size(); i++)
        {
            switch (
                const auto &child = children[i];
                child->rule
            )
            {
                case parser::ast::CONST_DECLARATION:
                case parser::ast::DECLARATION:
                {
                    // Create a new declaration
                    declaration decl;
                    decl.derive = derive; // Set the derive node if it exists
                    decl.is_const = child->rule == parser::ast::CONST_DECLARATION;

                    // Get the declaration's children
                    const auto &decl_children = child->children;
                    decl.value = decl_children[2];

                    // Convert the type
                    decl.decl_type = type(decl_children[1]);

                    // Append the declaration
                    m->declarations.emplace(
                        decl_children[0]->value.get(), // The identifier
                        container::move(decl) // Move the declaration
                    );

                    derive = nullptr; // Reset the derive node
                    break;
                }

                case parser::ast::DERIVE:
                {
                    // Store the last derive
                    derive = child;
                    break;
                }

                case parser::ast::FUNCTION:
                {
                    // Convert the function
                    function(m->functions, child, fun_allocator);
                    break;
                }

                default:
                {
                    // Impossible case
                }
            }
        }

        // Append the mod to the map
        map.emplace(name, m);
    }
}