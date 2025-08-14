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
// Created by rodri on 8/13/25.
//

#pragma once
#include "declaration.h"
#include "file_code/function.h"
#include "file_code/mod.h"

namespace zelix::analyzer::tool
{
    class ref
    {
        code::function *fun = nullptr;
        code::mod *mod = nullptr;
        declaration *decl = nullptr;

        template <typename T>
        T *impl_get()
        {
            if constexpr (std::is_same_v<T, declaration>)
            {
                return decl;
            }

            else if constexpr (std::is_same_v<T, code::function>)
            {
                return fun;
            }

            else if constexpr (std::is_same_v<T, code::mod>)
            {
                return mod;
            }

            else
            {
                // Invalid type
                static_assert(
                    false,
                    "Invalid type for ref::get()"
                );

                return nullptr;
            }
        }
    public:
        enum type
        {
            NONE,
            DECLARATION,
            FUN,
            MOD,
        };

        type kind = NONE;

        template <typename T>
        explicit ref(T *ptr)
        {
            if constexpr (std::is_same_v<T, code::function>)
            {
                fun = ptr;
            }

            else if constexpr (std::is_same_v<T, code::mod>)
            {
                mod = ptr;
            }

            else if constexpr (std::is_same_v<T, declaration>)
            {
                decl = ptr;
            }

            else
            {
                // Invalid type
                static_assert(
                    false,
                    "Invalid type for ref constructor"
                );
            }
        }

        template <typename T>
        T *get()
        {
            const auto ptr = impl_get<T>();
            if (ptr == nullptr)
            {
                throw except::exception("Reference is null");
            }

            return ptr;
        }
    };
}