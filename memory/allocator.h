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

#include "fluent/container/vector.h"
#include "fluent/except/exception.h"

namespace fluent::memory
{
    template <typename T, size_t N>
    class lazy_page
    {
        union storage_u {
            alignas(T) char raw[sizeof(T)];
        };

        storage_u storage[N];
        container::vector<storage_u*> mem_list;

    public:
        explicit lazy_page()
        {
            for (size_t i = 0; i < N; ++i)
                mem_list.push_back(&storage[i]);
        }

        T *alloc()
        {
            if (mem_list.empty())
            {
                throw except::exception("Out of memory in lazy page allocator");
            }

            // Get the next available pointer
            T* ptr = reinterpret_cast<T*>(mem_list.ref_at(mem_list.size() - 1));
            mem_list.calibrate(mem_list.size() - 1);

            // Placement new to construct object
            new (ptr) T();
            return ptr;
        }

        [[nodiscard]] bool full() const
        {
            return mem_list.empty();
        }
    };

    template <typename T, size_t N = 50>
    class lazy_allocator
    {
        container::vector<lazy_page<T, N>> pages;
        container::vector<T *> free_list;
    public:
        explicit lazy_allocator()
        {}

        T *alloc()
        {
            // Check the free list first
            if (!free_list.empty())
            {
                T *ptr = free_list.ref_at(free_list.size() - 1);
                free_list.pop_back();
                return ptr;
            }

            // See if we have any pages available
            if (pages.empty())
            {
                pages.emplace_back();
            }

            auto &back = pages.ref_at(pages.size() - 1);
            if (back.full())
            {
                // Allocate a new page
                pages.emplace_back();
                back = pages.ref_at(pages.size() - 1);
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