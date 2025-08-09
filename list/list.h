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
// Created by rodri on 8/9/25.
//

#pragma once

#include "memory/allocator.h"

// Zelix STL extension for a linked list
namespace zelix::container
{
    template <typename T>
    class list
    {
        struct node
        {
            T value;
            node *next = nullptr;
        };

        memory::lazy_allocator<node> allocator;
        node *head = nullptr;
        size_t size_ = 0;
    public:
        class iterator
        {
            node *current = nullptr;
        public:
            T next()
            {
                if (current == nullptr)
                {
                    throw except::exception("Iterator out of range");
                }

                T value = current->value;
                current = current->next;
                return value;
            }

            bool has_next()
            {
                return current != nullptr;
            }
        };

        void push_front(const T &value)
        {
            node *new_node = allocator.alloc();
            new_node->value = value;
            new_node->next = head;
            head = new_node;
            size_++;
        }

        [[nodiscard]] size_t size() const
        {
            return size_;
        }

        [[nodiscard]] bool empty() const
        {
            return size_ == 0;
        }

        iterator it()
        {
            return iterator{head};
        }
    };
}