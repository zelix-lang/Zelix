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
    /**
     * @brief A small-string-optimized, owned string class.
     *
     * Uses stack allocation for small strings and switches to heap allocation for larger strings.
     * Provides basic push and reserve operations, and exposes a C-style string interface.
     */
    class string
    {
        char stack[256]{}; ///< Stack-allocated buffer for small strings
        char* heap = nullptr; ///< Pointer to heap-allocated string for larger strings
        size_t len = 0; ///< Length of the string
        size_t max_capacity = sizeof(stack) - 1; ///< Maximum capacity of the stack buffer, -1 for null terminator
        size_t capacity = max_capacity + 1; ///< Capacity of the string
        bool stack_mem = true; ///< Flag to indicate if the string is using stack memory
        double growth_factor = 1.8; ///< Growth factor for heap allocation

        /**
         * @brief Initializes the heap-allocated buffer and copies stack data.
         *
         * Called when the string grows beyond the stack buffer.
         */
        void heap_init()
        {
            // Initialize the heap-allocated string if needed
            max_capacity = static_cast<size_t>(sizeof(stack) * growth_factor - 1); // Adjust max capacity based on growth factor
            capacity = max_capacity + 1;
            heap = new char[capacity];
            stack_mem = false;
            memcpy(heap, stack, sizeof(stack));
        }

        /**
         * @brief Reallocates the heap buffer to a larger size.
         *
         * Copies existing heap data to the new buffer and updates capacity.
         */
        void reallocate()
        {
            // Initialize the heap-allocated string if needed
            capacity = max_capacity + 1;
            const auto new_heap = new char[capacity];
            memcpy(new_heap, heap, len); // Copy existing data to new heap
            delete[] heap; // Free old heap memory
            heap = new_heap; // Update heap pointer
        }

    public:
        /**
         * @brief Default constructor. Initializes an empty string using stack memory.
         */
        explicit string() = default;

        /**
         * @brief Constructs a string with a specified capacity.
         * @param capacity The initial capacity to reserve.
         *
         * Uses stack memory if capacity is small, otherwise allocates on the heap.
         */
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

        /**
         * @brief Returns a pointer to a null-terminated C-style string.
         * @return Pointer to the string data.
         *
         * Ensures the string is null-terminated before returning.
         */
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

        /**
         * @brief Ensures the string has enough capacity for additional data.
         * @param required The additional capacity to reserve.
         *
         * Switches to heap memory or grows the heap buffer if needed.
         */
        void reserve(const size_t required)
        {
            if (stack_mem)
            {
                // Check if we need to switch to heap memory
                if (len + required > max_capacity)
                {
                    heap_init(); // Switch to heap memory
                }
            }

            if (!stack_mem && len + required > max_capacity)
            {
                // Check if we need to grow the heap memory
                auto new_capacity = static_cast<size_t>(capacity * growth_factor);
                const size_t count_to = len + required;

                while (new_capacity < count_to)
                {
                    new_capacity *= growth_factor; // Increase capacity by growth factor
                }

                max_capacity = new_capacity;
                reallocate();
            }
        }

        /**
         * @brief Appends a single character to the string.
         * @param c The character to append.
         */
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

        /**
         * @brief Appends a C-style string to the string.
         * @param c The null-terminated string to append.
         * @param c_len The length of the string to append. Defaults to strlen(c).
         */
        void push(const char *c, const size_t c_len)
        {
            reserve(c_len); // Reserve space for one character
            if (stack_mem)
            {
                memcpy(stack + len, c, c_len);
            }
            else
            {
                memcpy(heap + len, c, c_len);
            }

            len += c_len; // Update the length of the string
        }

        void push(const char *c)
        {
            push(c, strlen(c)); // Push with length
        }

        /**
         * @brief Returns the current length of the string.
         * @return The number of characters in the string.
         */
        [[nodiscard]] size_t size()
        const {
            return len;
        }

        /**
         * @brief Destructor. Releases heap memory if allocated.
         */
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