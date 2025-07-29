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
#include <optional>
#include <stdexcept>

namespace fluent::container
{
    /**
     * @brief A hybrid vector implementation with stack and heap allocation.
     *
     * Uses stack allocation for small vectors (up to 100 elements), and switches to heap allocation for larger sizes.
     * Provides basic vector operations similar to std::vector, with custom growth factor and capacity management.
     *
     * @tparam T Type of elements stored in the vector.
     */
    template <typename T>
    class vector
    {
        size_t len = 0; ///< Length of the vector
        size_t capacity = 25; ///< Initial capacity of the vector
        double growth_factor = 2; ///< Growth factor for resizing
        T stack[100]; ///< Stack-allocated array for small vectors
        T *data = nullptr; ///< Pointer to the data
        bool stack_mem = true; ///< Flag to indicate if we are using stack memory

        /**
         * @brief Initializes heap memory for the vector.
         *
         * Allocates memory for the vector on the heap and copies existing elements from the stack if necessary.
         * Throws std::bad_alloc if memory allocation fails.
         */
        void init()
        {
            // Throws std::bad_alloc if memory allocation fails
            data = new T[capacity];

            if (len > 0)
            {
                // Copy elements from stack to heap if we have already used the stack
                for (size_t i = 0; i < len; ++i)
                {
                    data[i] = stack[i]; // Move elements from stack to heap
                }
            }
        }

        /**
         * @brief Reallocates heap memory with the current capacity.
         *
         * Moves existing elements to a new heap-allocated array and updates the data pointer.
         * Throws std::bad_alloc if memory allocation fails.
         */
        void realloc_()
        {
            // Create a new array
            T *new_data = new T[capacity]; // Assume capacity is already set
            if (!new_data)
            {
                throw std::bad_alloc(); // Memory allocation failed
            }

            // Copy existing elements to the new array
            for (size_t i = 0; i < len; ++i)
            {
                new_data[i] = std::move(data[i]); // Move elements to the new array
            }

            delete[] data;
            data = new_data; // Update the data pointer to the new array
        }

        /**
         * @brief Sets the value at the specified index in heap memory.
         *
         * @param idx Index to set.
         * @param value Value to assign.
         */
        void set(const size_t idx, const T &value)
        {
            data[idx] = value;
        }

        /**
         * @brief Sets the value at the specified index in stack memory.
         *
         * Throws std::out_of_range if the index is out of bounds.
         *
         * @param idx Index to set.
         * @param value Value to assign.
         */
        void set_stack(const size_t idx, const T &value)
        {
            if (idx >= 100)
            {
                throw std::out_of_range("Index out of bounds for stack-allocated array");
            }

            stack[idx] = value;
        }

    public:
        /**
         * @brief Default constructor. Uses stack allocation.
         */
        explicit vector()
        {};
        /*{
            // <init();>
            // Do not allocate right away
        }*/

        /**
         * @brief Constructs a vector with a specified capacity.
         *
         * Allocates heap memory if capacity is specified.
         * Throws std::invalid_argument if capacity is zero.
         *
         * @param capacity Initial capacity.
         */
        explicit vector(size_t capacity)
        {
            if (capacity == 0)
            {
                throw std::invalid_argument("Capacity must be greater than zero");
            }

            this->capacity = capacity;
            // Do not initialize heap memory if the capacity is less than 100
            if (capacity < 100)
            {
                stack_mem = true; // Use stack memory
                return; // No need to allocate heap memory
            }

            init(); // Init on user request
            stack_mem = false;
        }

        /**
         * @brief Constructs a vector with a specified capacity and growth factor.
         *
         * Allocates heap memory and sets the growth factor.
         * Throws std::invalid_argument if capacity is zero or growth factor is invalid.
         *
         * @param capacity Initial capacity.
         * @param growth_factor Growth factor for resizing.
         */
        explicit vector(size_t capacity, double growth_factor)
        {
            if (capacity == 0)
            {
                throw std::invalid_argument("Capacity must be greater than zero");
            }

            if (growth_factor <= 1)
            {
                throw std::invalid_argument("Growth factor must be greater than 1");
            }

            if (growth_factor > 10)
            {
                throw std::invalid_argument("Growth factor must not exceed 10");
            }

            this->growth_factor = growth_factor;
            this->capacity = capacity;

            // Do not initialize heap memory if the capacity is less than 100
            if (capacity < 100)
            {
                stack_mem = true; // Use stack memory
                return; // No need to allocate heap memory
            }

            init(); // Init on user request
            stack_mem = false; // Use heap memory
        }

        /**
         * @brief Ensures the vector has at least the specified capacity.
         *
         * Reallocates memory if needed.
         *
         * @param capacity Minimum required capacity.
         */
        void reserve(size_t capacity)
        {
            if (this->capacity >= capacity)
            {
                return; // Already have enough capacity
            }

            // Check if we are on stack memory
            if (stack_mem)
            {
                // If we are using stack memory, we need to switch to heap memory
                if (capacity < 100)
                {
                    return; // No need to allocate heap memory for small capacity
                }

                init(); // Initialize heap memory
                stack_mem = false; // Switch to heap memory
            }

            this->capacity = capacity;
            realloc_();
        }

        /**
         * @brief Adds an element to the end of the vector.
         *
         * Switches to heap allocation if stack limit is exceeded.
         * Resizes if needed.
         *
         * @param element Element to add.
         */
        void push_back(const T &element)
        {
            if (stack_mem)
            {
                stack[len] = element;
                len++;

                if (len >= 100)
                {
                    // Move to heap allocation if we exceed stack size
                    if (capacity < 100)
                    {
                        capacity = 100; // Ensure we have enough capacity
                    }

                    init();
                    stack_mem = false; // Switch to heap memory
                }

                return; // Element added to stack
            }

            // Check if we need to resize
            if (len >= capacity)
            {
                // Resize the vector
                capacity = static_cast<size_t>(capacity * growth_factor);
                realloc_();
            }

            // Add the element to the end
            data[len] = element;
            len++;
        }

        /**
         * @brief Constructs an element in place at the end of the vector.
         *
         * Switches to heap allocation if stack limit is exceeded.
         * Resizes if needed.
         *
         * @tparam Args Argument types for element construction.
         * @param args Arguments to forward to the element constructor.
         */
        template <typename... Args>
        void emplace_back(Args &&...args)
        {
            if (stack_mem)
            {
                new (&stack[len]) T(std::forward<Args>(args)...);
                len++;

                if (len >= 100)
                {
                    // Move to heap allocation if we exceed stack size
                    if (capacity < 100)
                    {
                        capacity = 100; // Ensure we have enough capacity
                    }

                    init();
                    stack_mem = false; // Switch to heap memory
                }

                return; // Element added to stack
            }

            // Check if we need to resize
            if (len >= capacity)
            {
                // Resize the vector
                capacity = static_cast<size_t>(capacity * growth_factor);
                realloc_();
            }

            // Construct the element in place at the end of the vector
            new (&data[len]) T(std::forward<Args>(args)...);
            len++;
        }

        /**
         * @brief Returns the last element in the vector, if any.
         *
         * @return std::optional<T> Last element or std::nullopt if empty.
         */
        std::optional<T> back()
        {
            // Make sure the length is greater than 0
            if (len == 0)
            {
                return std::nullopt; // No elements in the vector
            }

            return data[len - 1]; // Return the last element
        }

        /**
         * @brief Removes the last element from the vector.
         *
         * Throws std::out_of_range if the vector is empty.
         */
        void pop_back()
        {
            // Make sure the length is greater than 0
            if (len == 0)
            {
                throw std::out_of_range("Cannot pop from an empty vector");
            }

            if (stack_mem)
            {
                // Call the destructor for the last element
                stack[len - 1].~T();
            }
            else
            {
                // Call the destructor for the last element
                data[len - 1].~T();
            }

            len--; // Decrease the length
        }

        /**
         * @brief Manually sets the length of the vector.
         *
         * Warning: Use with caution. This does not construct or destruct elements.
         * May lead to undefined behavior if `n` is greater than the current capacity or
         * if elements are not properly initialized.
         *
         * @param n New length to set.
         */
        void calibrate(const size_t n)
        {
            // Warning: use with caution
            this->len = n;
        }

        /**
         * @brief Returns the current number of elements in the vector.
         *
         * @return size_t Number of elements.
         */
        [[nodiscard]] size_t size() const
        {
            return len; // Return the current length of the vector
        }

        /**
         * @brief Returns the current capacity of the vector.
         *
         * @return size_t Capacity.
         */
        [[nodiscard]] size_t get_capacity() const
        {
            return capacity; // Return the current capacity of the vector
        }

        /**
         * @brief Checks if the vector is empty.
         *
         * @return true if empty, false otherwise.
         */
        [[nodiscard]] bool empty() const
        {
            return len == 0; // Return true if the vector is empty
        }

        /**
         * @brief Returns a pointer to the underlying data.
         *
         * @return T* Pointer to data.
         */
        [[nodiscard]] T *data_ptr()
        {
            return stack_mem ? stack : data; // Return the pointer to the data
        }

        /**
         * @brief Returns the element at the specified index, if within bounds.
         *
         * @param index Index to access.
         * @return std::optional<T> Element or std::nullopt if out of bounds.
         */
        std::optional<T> at(const size_t index) const
        {
            // Check if the index is within bounds
            if (index >= len)
            {
                return std::nullopt; // Index out of bounds
            }

            if (stack_mem)
            {
                return stack[index]; // Return the element at the specified index
            }

            return data[index]; // Return the element at the specified index
        }

        /**
         * @brief Returns a reference to the element at the specified index.
         *
         * Throws std::out_of_range if index is out of bounds.
         *
         * @param index Index to access.
         * @return T& Reference to element.
         */
        T &operator[](const size_t index)
        {
            if (index >= len)
            {
                throw std::out_of_range("Index out of bounds");
            }

            if (stack_mem)
            {
                return stack[index]; // Return a reference to the element at the specified index
            }

            return data[index]; // Return a reference to the element at the specified index
        }

        /**
         * @brief Copy constructor.
         *
         * Copies elements from another vector.
         *
         * @param other Vector to copy from.
         */
        vector(const vector& other)
        {
            // Avoid copying
            clear();
            reserve(other.size());

            for (size_t i = 0; i < other.len; ++i)
            {
                // Use placement new to construct the element in place
                emplace_back(other[i]);
            }
        }

        /**
         * @brief Move constructor.
         *
         * Moves resources from another vector.
         *
         * @param other Vector to move from.
         */
        vector(vector&& other) noexcept
            : len(other.len),
            capacity(other.capacity),
            growth_factor(other.growth_factor),
            data(other.data),
            stack_mem(other.stack_mem)
        {
            other.data = nullptr;
            other.capacity = 0;

            if (other.stack_mem)
            {
                // Move stack memory to this instance
                for (size_t i = 0; i < other.len; i++)
                {
                    stack[i] = std::move(other.stack[i]); // Move elements from other stack to this stack
                }
            }

            other.len = 0;
            other.stack_mem = true; // Reset the other vector to stack memory
        }

        /**
         * @brief Move assignment operator.
         *
         * Moves resources from another vector.
         *
         * @param other Vector to move from.
         * @return vector& Reference to this vector.
         */
        vector& operator=(vector&& other) noexcept
        {
            if (this != &other)
            {
                clear();
                if (!stack_mem)
                {
                    delete[] data;
                }

                len = other.len;
                capacity = other.capacity;
                growth_factor = other.growth_factor;
                data = other.data;
                stack_mem = other.stack_mem;

                if (stack_mem)
                {
                    // Copy stack memory from other
                    for (size_t i = 0; i < other.len; i++)
                    {
                        stack[i] = std::move(other.stack[i]); // Move elements from other stack to this stack
                    }
                }

                other.data = nullptr;
                other.len = 0;
                other.capacity = 0;
            }
            return *this;
        }

        /**
         * @brief Copy assignment operator.
         *
         * Copies elements from another vector.
         *
         * @param other Vector to copy from.
         * @return vector& Reference to this vector.
         */
        vector& operator=(const vector& other)
        {
            if (this != &other)
            {
                clear();
                reserve(other.len);
                for (size_t i = 0; i < other.len; ++i)
                {
                    emplace_back(other[i]);
                }
            }
            return *this;
        }

        /**
         * @brief Returns a pointer to the beginning of the vector.
         *
         * @return T* Pointer to beginning.
         */
        T *begin()
        {
            if (stack_mem)
            {
                return stack; // Return pointer to the beginning of the stack-allocated array
            }

            return data;
        }

        /**
         * @brief Returns a pointer to the end of the vector.
         *
         * @return T* Pointer to end.
         */
        T *end()
        {
            if (stack_mem)
            {
                return stack + len; // Return pointer to the end of the stack-allocated array
            }

            return data + len; // Return pointer to the end of the vector
        }

        const T* begin() const
        {
            if (stack_mem)
            {
                return stack; // Return pointer to the beginning of the stack-allocated array
            }

            return data;
        }

        const T* end() const
        {
            if (stack_mem)
            {
                return stack + len; // Return pointer to the end of the stack-allocated array
            }

            return data + len;
        }

        /**
         * @brief Clears the vector, destroying all elements.
         */
        void clear()
        {
            if (stack_mem)
            {
                for (size_t i = 0; i < len; ++i)
                {
                    stack[i].~T(); // Call the destructor for each element
                }
            }
            else
            {
                for (size_t i = 0; i < len; ++i)
                {
                    data[i].~T(); // Call the destructor for each element
                }
            }

            len = 0;
        }

        /**
         * @brief Destructor. Cleans up resources.
         */
        ~vector()
        {
            clear();

            if (!stack_mem)
            {
                delete[] data; // Free the allocated memory
            }
        }
    };
}