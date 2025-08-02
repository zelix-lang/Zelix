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
    template <typename T>
    class lazy_page
    {
        T *ptr = nullptr;
        size_t size = 0;
        size_t used = 0;

    public:
        explicit lazy_page()
        {
            size = sizeof(T) * 20;
            ptr = new T[20];
        }

        explicit lazy_page(const size_t page_size)
        {
            size = sizeof(T) * page_size;
            ptr = new T[page_size];
        }

        T *alloc()
        {
            if (used >= size)
            {
                throw except::exception("Out of memory in lazy page allocator");
            }

            // Get the next available pointer
            T *next_ptr = reinterpret_cast<T *>(reinterpret_cast<char *>(ptr) + used);
            used += sizeof(T);

            // Initialize the allocated memory
            new (next_ptr) T(); // Placement new to construct the object in the allocated memory
            return next_ptr;
        }

        [[nodiscard]] bool full() const
        {
            return used >= size;
        }

        ~lazy_page()
        {
            delete[] ptr; // Free the allocated memory
            ptr = nullptr; // Avoid dangling pointer
        }
    };

    template <typename T>
    class lazy_allocator
    {
        container::vector<lazy_page<T>> pages;
        container::vector<T *> free_list;
        size_t page_size;
    public:
        explicit lazy_allocator(const size_t page_size)
            : page_size(page_size)
        {
            if (page_size == 0)
            {
                throw except::exception("Page size cannot be zero");
            }
        }

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
                pages.emplace_back(page_size);
            }

            auto &back = pages.ref_at(pages.size() - 1);
            if (back.full())
            {
                // Allocate a new page
                pages.emplace_back(page_size);
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