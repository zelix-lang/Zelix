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
#include "function.h"
#include "mod.h"

namespace zelix::code
{
    // Forward declaration of symbol
    class symbol;

    using package = ankerl::unordered_dense::map<
        container::external_string,
        symbol *,
        container::external_string_hash
    >;

    class symbol
    {
        function *func = nullptr; // Pointer to the function symbol
        mod *mod = nullptr; // Pointer to the module symbol
        declaration *decl = nullptr; // Pointer to the declaration symbol
        ankerl::unordered_dense::map<
            container::external_string,
            symbol *,
            container::external_string_hash
        > *package = nullptr; // Pointer to the package symbol

    public:
        template <typename T>
        explicit symbol(T* ptr)
        {
            if constexpr (std::is_same_v<T, function>)
            {
                func = ptr;
            }
            else if constexpr (std::is_same_v<T, code::mod>)
            {
                mod = ptr;
            }
            else if constexpr (std::is_same_v<T, declaration>)
            {
                decl = ptr;
            }
            else if constexpr (std::is_same_v<T, code::package>)
            {
                package = ptr;
            }
            else
            {
                static_assert(
                    std::is_same_v<T, void>,
                    "Unsupported type for symbol"
                );
            }
        }

        template <typename T>
        [[nodiscard]] T *get() const
        {
            if constexpr (std::is_same_v<T, function>)
            {
                return func;
            }
            else if constexpr (std::is_same_v<T, code::mod>)
            {
                return mod;
            }
            else if constexpr (std::is_same_v<T, declaration>)
            {
                return decl;
            }
            else if constexpr (std::is_same_v<T, code::package>)
            {
                return package;
            }
            else
            {
                static_assert(
                    std::is_same_v<T, void>,
                    "Unsupported type for symbol"
                );
                return nullptr; // This line will never be reached
            }
        }

        template <typename T>
        [[nodiscard]] bool is() const
        {
            if constexpr (std::is_same_v<T, function>)
            {
                return func != nullptr;
            }
            else if constexpr (std::is_same_v<T, code::mod>)
            {
                return mod != nullptr;
            }
            else if constexpr (std::is_same_v<T, declaration>)
            {
                return decl != nullptr;
            }
            else if constexpr (std::is_same_v<T, code::package>)
            {
                return package != nullptr;
            }
            else
            {
                static_assert(
                    std::is_same_v<T, void>,
                    "Unsupported type for symbol"
                );
                return false; // This line will never be reached
            }
        }
    };
}