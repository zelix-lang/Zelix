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
// Created by rodrigo on 8/1/25.
//

#pragma once

#include <cstddef>
#include <cstdlib>
#include <new>
#include "absl/container/inlined_vector.h"
#include "fluent/except/exception.h"

namespace fluent::memory
{
    template <typename T>
    class lazy_page
    {
        std::byte *buffer = nullptr;
        size_t capacity = 0;
        size_t offset = 0;

    public:
        explicit lazy_page(const size_t page_size = 512)
            : capacity(page_size)
        {
            buffer = static_cast<std::byte*>(malloc(page_size * sizeof(T)));
            if (!buffer) throw std::bad_alloc();
        }

        T *alloc()
        {
            if (offset >= capacity)
            {
                throw except::exception("Out of memory in lazy page allocator");
            }

            // Allocate the next object in the buffer
            T* ptr = reinterpret_cast<T*>(buffer + offset * sizeof(T));
            ++offset;

            // Placement new to construct object
            new (ptr) T();
            return ptr;
        }

        [[nodiscard]] bool full() const
        {
            return offset >= capacity;
        }
    };

    template <typename T>
    class lazy_allocator
    {
        absl::InlinedVector<lazy_page<T>, 4> pages;
        std::vector<T *> free_list;
        size_t page_size = 512;

    public:
        explicit lazy_allocator(const size_t page_size = 512)
            : page_size(page_size)
        {}

        T *alloc()
        {
            // Check the free list first
            if (!free_list.empty())
            {
                T *ptr = free_list[free_list.size() - 1];
                free_list.pop_back();
                return ptr;
            }

            // See if we have any pages available
            if (pages.empty())
            {
                pages.emplace_back(page_size);
            }

            auto &back = pages[pages.size() - 1];
            if (back.full())
            {
                // Allocate a new page
                pages.emplace_back(page_size);
                back = pages[pages.size() - 1];
                return back.alloc();
            }

            return back.alloc();
        }

        void dealloc(T *ptr)
        {
            ptr->~T(); // Call the destructor for the object
            free_list.push_back(ptr);
        }

        ~lazy_allocator()
        {
            pages.clear(); // Clear the vector of pages
        }
    };
}