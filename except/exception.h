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
#include <exception>

namespace fluent::except
{
    class exception final : public std::exception
    {
        const char *message; ///< Error message

    public:
        /**
         * @brief Constructs an exception with a given message.
         * @param msg The error message to associate with the exception.
         */
        explicit exception(const char *msg)
            : message(msg)
        {}

        [[nodiscard]] const char *what()
        const noexcept override {
            return message; ///< Returns the error message
        }
    };
}
