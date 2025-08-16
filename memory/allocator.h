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

#include <cstddef>
#include "zelix/container/vector.h"

namespace zelix::memory
{
    template <typename T, size_t Capacity, bool CallDestructors>
    class lazy_page
    {
        std::byte *buffer = nullptr;
        size_t offset = 0;

    public:
        lazy_page();
        T *alloc(auto&&... args);

        [[nodiscard]] bool full() const
        {
            return offset >= Capacity;
        }

        ~lazy_page();
    };

    template <typename T, size_t Capacity = 256, bool CallDestructors = false>
    class lazy_allocator
    {
        container::vector<lazy_page<T, Capacity, CallDestructors>> pages;
        container::vector<T *> free_list;

    public:
        T *alloc(auto&&... args);
        void dealloc(T *ptr);
        ~lazy_allocator();
    };
}