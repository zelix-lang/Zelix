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
// Created by rodri on 8/15/25.
//

#pragma once
#include "exception/unresolved_symbol.h"
#include "memory/allocator.h"
#include "symbol.h"

namespace zelix::code
{
    using package = ankerl::unordered_dense::map<
        container::external_string,
        symbol,
        container::external_string_hash
    >;
    class program
    {
        memory::lazy_allocator<symbol> symbol_alloc;
        memory::lazy_allocator<mod> mod_alloc;
        memory::lazy_allocator<function> function_alloc;
        memory::lazy_allocator<declaration> declaration_alloc;
        ankerl::unordered_dense::map<
            container::external_string,
            package,
            container::external_string_hash
        > context; // The global context

    public:
        template <typename T>
        T *alloc()
        {
            if constexpr (std::is_same_v<T, function>)
            {
                return function_alloc.alloc();
            }
            else if constexpr (std::is_same_v<T, mod>)
            {
                return mod_alloc.alloc();
            }
            else if constexpr (std::is_same_v<T, declaration>)
            {
                return declaration_alloc.alloc();
            }
            else
            {
                throw except::exception("Unknown type for program symbol");
            }
        }

        package &pkg(const container::external_string &str)
        {
            // See if the package exists
            if (!context.contains(str))
            {
                throw exception::unresolved_symbol(str.ptr());
            }

            return context[str];
        }

        package &new_pkg(const container::external_string &str)
        {
            // See if the package exists
            if (!context.contains(str))
            {
                // Insert the package directly
                context.try_emplace(str);
            }

            // Return the package
            return context[str];
        }
    };
}