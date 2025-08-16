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
// Created by rodri on 8/16/25.
//

#include <cstddef>
#include <cstdlib>
#include <new>
#include "zelix/except/exception.h"
#include "zelix/container/forward.h"
#include "allocator.h"

#include "file/converter.h"
#include "parser/parser.h"
#include "parser/rule/expr/queue.h"

template<typename T, size_t Capacity, bool CallDestructors>
zelix::memory::lazy_page<T, Capacity, CallDestructors>::lazy_page()
{
    buffer = static_cast<std::byte*>(malloc(Capacity * sizeof(T)));
    if (!buffer) throw std::bad_alloc();
}

template<typename T, size_t Capacity, bool CallDestructors>
T *zelix::memory::lazy_page<T, Capacity, CallDestructors>::alloc(auto&&... args)
{
    if (offset >= Capacity)
    {
        throw except::exception("Out of memory in lazy page allocator");
    }

    // Allocate the next object in the buffer
    T* ptr = reinterpret_cast<T*>(buffer + offset * sizeof(T));
    ++offset;
    new (ptr) T(container::forward<decltype(args)>(args)...); // Construct the object in place
    return ptr;
}

template<typename T, size_t Capacity, bool CallDestructors>
zelix::memory::lazy_page<T, Capacity, CallDestructors>::~lazy_page()
{
    if constexpr (CallDestructors)
    {
        // Call the destructor of all allocated objects
        for (size_t i = 0; i < offset; ++i)
        {
            T *ptr = reinterpret_cast<T*>(buffer + i * sizeof(T));
            ptr->~T(); // Call the destructor
        }
    }

    free(buffer);
}

template<typename T, size_t Capacity, bool CallDestructors>
T *zelix::memory::lazy_allocator<T, Capacity, CallDestructors>::alloc(auto &&...args)
{
    // Check the free list first
    if (!free_list.empty())
    {
        T *ptr = free_list[free_list.size() - 1];
        free_list.pop_back();
        new (ptr) T(container::forward<decltype(args)>(args)...); // Placement new to construct the object
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
        return back.alloc(container::forward<decltype(args)>(args)...);
    }

    return back.alloc(container::forward<decltype(args)>(args)...);
}

template<typename T, size_t Capacity, bool CallDestructors>
void zelix::memory::lazy_allocator<T, Capacity, CallDestructors>::dealloc(T *ptr)
{
    ptr->~T(); // Call the destructor
    free_list.push_back(ptr);
}

template<typename T, size_t Capacity, bool CallDestructors>
zelix::memory::lazy_allocator<T, Capacity, CallDestructors>::~lazy_allocator()
{
    pages.clear(); // Clear the vector of pages
}

template struct zelix::memory::lazy_page<zelix::code::symbol, 256ul, false>;
template struct zelix::memory::lazy_page<zelix::code::mod, 256ul, false>;
template struct zelix::memory::lazy_page<zelix::code::function, 256ul, false>;
template struct zelix::memory::lazy_page<zelix::code::declaration, 256ul, false>;
template class zelix::memory::lazy_page<zelix::lexer::token, 256ul, false>;
template class zelix::memory::lazy_page<zelix::parser::ast, 256ul, false>;
template class zelix::memory::lazy_page<zelix::parser::rule::expr::queue_node, 256ul, false>;

// Page templates
template
zelix::lexer::token*
    zelix::memory::lazy_page<zelix::lexer::token, 256ul, false>::alloc<>();

template
zelix::parser::ast*
    zelix::memory::lazy_page<zelix::parser::ast, 256ul, false>::alloc<>();

template
zelix::parser::rule::expr::queue_node*
    zelix::memory::lazy_page<zelix::parser::rule::expr::queue_node, 256ul, false>::alloc<>();
