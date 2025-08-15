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
#include <cstdlib>
#include <new>
#include "zelix/except/exception.h"

namespace zelix::memory
{
    template <typename T, size_t Capacity>
    class lazy_page
    {
        std::byte *buffer = nullptr;
        size_t offset = 0;

    public:
        explicit lazy_page()
        {
            buffer = static_cast<std::byte*>(malloc(Capacity * sizeof(T)));
            if (!buffer) throw std::bad_alloc();
        }

        T *alloc()
        {
            if (offset >= Capacity)
            {
                throw except::exception("Out of memory in lazy page allocator");
            }

            // Allocate the next object in the buffer
            T* ptr = reinterpret_cast<T*>(buffer + offset * sizeof(T));
            ++offset;
            new (ptr) T(); // Construct the object in place
            return ptr;
        }

        [[nodiscard]] bool full() const
        {
            return offset >= Capacity;
        }
    };

    template <typename T, size_t Capacity = 256>
    class lazy_allocator
    {
        container::vector<lazy_page<T, Capacity>> pages;
        container::vector<T *> free_list;

    public:
        template <typename... Args>
        T *alloc(Args &&...args)
        {
            // Check the free list first
            if (!free_list.empty())
            {
                T *ptr = free_list[free_list.size() - 1];
                free_list.pop_back();
                new (ptr) T(container::forward<Args>(args)...); // Placement new to construct the object
                return ptr;
            }

            // See if we have any pages available
            if (pages.empty())
            {
                pages.emplace_back();
            }

            auto &back = pages[pages.size() - 1];
            if (back.full())
            {
                // Allocate a new page
                pages.emplace_back();
                back = pages[pages.size() - 1];
                return back.alloc(container::forward<Args>(args)...);
            }

            return back.alloc(container::forward<Args>(args)...);
        }

        void dealloc(T *ptr)
        {
            ptr->~T(); // Call the destructor
            free_list.push_back(ptr);
        }

        ~lazy_allocator()
        {
            pages.clear(); // Clear the vector of pages
        }
    };
}