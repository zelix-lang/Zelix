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
#include "memory/allocator.h"
#include "symbol.h"

namespace zelix::code
{
    class program
    {
        memory::lazy_allocator<package, 64, true> package_alloc;
        memory::lazy_allocator<symbol> symbol_alloc;
        memory::lazy_allocator<mod> mod_alloc;
        memory::lazy_allocator<function> function_alloc;
        memory::lazy_allocator<declaration> declaration_alloc;
        package *package = package_alloc.alloc();

        static void throw_mismatch()
        {
            throw except::exception("Symbol type mismatch");
        }

        static void throw_does_not_exist()
        {
            throw except::exception("Symbol does not exist");
        }

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
            else if constexpr (std::is_same_v<T, code::package>)
            {
                return package_alloc.alloc();
            }
            else
            {
                throw except::exception("Unknown type for program symbol");
            }
        }

        template <typename T>
        void set(container::external_string &name, T *ptr)
        {
            package->try_emplace(name, ptr);
        }

        template <typename T>
        T *set(container::external_string &name)
        {
            // Do not allocate if the package is already there
            if (package->contains(name))
            {
                return static_cast<T *>(package->at(name));
            }

            // Allocate a new instance of T
            T *ptr = alloc<T>();

            // Allocate a new symbol
            auto symbol = symbol_alloc.alloc(ptr);
            package->try_emplace(name, symbol);
            return ptr;
        }

        template <typename T>
        T *resolve(container::external_string &name)
        {
            const auto it = package->find(name);
            if (it == package->end())
            {
                throw_does_not_exist();
            }

            const auto &symbol = it->second;
            if (!symbol->is<T>())
            {
                throw_mismatch();
            }

            return symbol->get<T>();
        }

        template <typename T = code::package, bool RetrievePackage = false>
        T *resolve(parser::ast *package_node)
        {
            static_assert(
                std::is_same_v<T, code::package> ||
                std::is_same_v<T, mod>,
                "Unsupported package type"
            );

            // Get the children
            const auto &children = package_node->children;
            const auto size = children.size() - 1;
            code::package *last_pkg = package;

            for (size_t i = 0; i <= size; ++i)
            {
                const auto &child = children[i];

                // Get the child's value
                const auto &child_val = child->value.get();

                // Make sure that the value exists
                if (!last_pkg->contains(child_val))
                {
                    throw_does_not_exist();
                }

                // Get the symbol
                const auto symbol = last_pkg->at(child_val);
                if (symbol->is<mod>())
                {
                    if (i != size)
                    {
                        throw_mismatch();
                    }

                    // Make sure the types match
                    if constexpr (std::is_same_v<T, mod>)
                    {
                        throw_mismatch();
                    }

                    return symbol->get<mod>();
                }

                if (symbol->is<code::package>())
                {
                    if constexpr (RetrievePackage)
                    {
                        if (i == size)
                        {
                            throw_mismatch();
                        }
                    }

                    return symbol->get<code::package>();
                }

                throw_mismatch();
            }

            if constexpr (std::is_same_v<T, code::package>)
            {
                return last_pkg;
            }

            // Unreachable code
            return nullptr;
        }
    };
}