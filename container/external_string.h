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

//
// Created by rodrigo on 7/29/25.
//

#pragma once
#include <cstring>
#include <xxh3.h>
#include "except/exception.h"

namespace fluent::container
{
    class external_string
    {
        char *buffer;
        size_t len = 0;

    public:
        external_string(char *buffer, const size_t len)
            : buffer(buffer), len(len)
        {
            if (buffer == nullptr || len == 0)
            {
                throw except::exception("Buffer cannot be null or length zero");
            }
        }

        [[nodiscard]] char *ptr()
        const {
            return buffer;
        }

        [[nodiscard]] size_t size()
        const {
            return len;
        }
    };

    struct external_string_hash
    {
        using is_transparent = void;

        size_t operator()(const external_string &str) const
        {
            // Use xxHash
            return XXH3_64bits(str.ptr(), str.size());
        }

        size_t operator()(const char* c_str) const {
            const size_t len = strlen(c_str);
            return XXH64(c_str, len, len);
        }
    };
}