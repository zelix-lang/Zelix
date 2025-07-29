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

namespace fluent::container
{
    class string
    {
        char stack[256]{}; // Stack-allocated buffer for small strings
        char* heap = nullptr; // Pointer to heap-allocated string for larger strings
        size_t len = 0; // Length of the string
        size_t max_capacity = sizeof(stack) - 1; // Maximum capacity of the stack buffer, -1 for null terminator
        size_t capacity = max_capacity + 1; // Capacity of the string
        bool stack_mem = true; // Flag to indicate if the string is using stack memory
        double growth_factor = 1.8; // Growth factor for heap allocation

        void heap_init()
        {
            // Initialize the heap-allocated string if needed
            max_capacity = static_cast<size_t>(sizeof(stack) * growth_factor - 1); // Adjust max capacity based on growth factor
            capacity = max_capacity + 1;
            heap = new char[capacity];
            stack_mem = false;
            memcpy(heap, stack, sizeof(stack));
        }

        void reallocate()
        {
            // Initialize the heap-allocated string if needed
            max_capacity = static_cast<size_t>(sizeof(stack) * growth_factor - 1); // Adjust max capacity based on growth factor
            capacity = max_capacity + 1;
            const auto new_heap = new char[capacity];
            memcpy(new_heap, heap, len); // Copy existing data to new heap
            delete[] heap; // Free old heap memory
            heap = new_heap; // Update heap pointer
        }

    public:
        explicit string() = default;

        explicit string(const size_t capacity)
        {
            // Determine if we have to use the stack or heap
            if (capacity <= max_capacity) // -1 to account for null terminator
            {
                this->capacity = sizeof(stack);
                stack_mem = true; // Use stack memory
            }
            else
            {
                this->capacity = capacity;
                heap_init(); // Initialize heap buffer
            }
        }

        char *c_str()
        {
            if (stack_mem)
            {
                stack[len] = '\0'; // Ensure null termination for stack memory
                return stack;
            }
            else
            {
                if (heap)
                {
                    heap[len] = '\0'; // Ensure null termination for heap memory
                    return heap;
                }

                return nullptr; // Return nullptr if no memory is allocated
            }
        }

        void reserve(const size_t capacity)
        {
            if (stack_mem)
            {
                // Check if we need to switch to heap memory
                if (len + capacity > max_capacity)
                {
                    heap_init(); // Switch to heap memory
                }
            }

            // Check if we need to grow the heap memory
            if (len + capacity > max_capacity)
            {
                reallocate();
            }
        }

        void push(const char c)
        {
            reserve(1); // Reserve space for one character
            if (stack_mem)
            {
                stack[len++] = c; // Add character to stack memory
            }
            else
            {
                heap[len++] = c; // Add character to heap memory
            }
        }

        void push(const char *c)
        {
            const size_t c_len = strlen(c); // Get the length of the input string
            reserve(c_len); // Reserve space for one character
            if (stack_mem)
            {
                memcpy(stack + len, c, c_len);
            }
            else
            {
                memcpy(heap + len, c, c_len);
            }
        }

        ~string()
        {
            if (!stack_mem && heap)
            {
                delete[] heap; // Free heap memory if allocated
                heap = nullptr;
            }
        }
    };
}
