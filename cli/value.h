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

namespace fluent::cli
{
    class value
    {
    public:
        enum type
        {
            STRING,
            INTEGER,
            FLOAT,
            BOOL
        };

    private:
        type value_type = STRING; ///< Type of the value
        container::external_string description; ///< Description of the value

        // Default values
        container::external_string default_str;
        container::external_string alias_str;
        std::optional<container::external_string> original_name;
        bool alias = false;
        int default_int = 0;
        float default_float = 0.0f;
        bool default_bool = false;

    public:
        explicit value() :
            description("", 1), default_str("", 1),
            alias_str("", 1), original_name(std::nullopt)
        {}

        /**
         * @brief Constructs a value with a given type and description.
         * @param default_value The default value of the value.
         * @param description The description of the value.
         * @param alias_str The alias string for the value.
         * @param original_name The original name of the value, if it is an alias.
         * @param alias Whether this value is an alias for another value.
         */
        template <typename T>
        explicit value(
            const T default_value,
            const container::external_string &description,
            const container::external_string &alias_str,
            const std::optional<container::external_string> &original_name,
            const bool alias
        ) :
            description(description), default_str("", 1),
            alias_str(alias_str),
            original_name(original_name), alias(alias)
        {
            if (description.size() == 0)
            {
                throw except::exception("Description cannot be empty");
            }

            if constexpr (std::is_same_v<T, container::external_string>)
            {
                value_type = STRING;
                default_str = default_value;
            }
            else if constexpr (std::is_same_v<T, int>)
            {
                value_type = INTEGER;
                default_int = default_value;
            }
            else if constexpr (std::is_same_v<T, float>)
            {
                value_type = FLOAT;
                default_float = default_value;
            }
            else if constexpr (std::is_same_v<T, bool>)
            {
                value_type = BOOL;
                default_bool = default_value;
            }
            else if constexpr (std::is_same_v<T, const char *>)
            {
                value_type = STRING;
                default_str = container::external_string(default_value);
            }
            else
            {
                throw except::exception("Unsupported type for default value");
            }
        }

        [[nodiscard]] type get_type() const
        {
            return value_type;
        }

        [[nodiscard]] const container::external_string &get_description()
        const
        {
            return description;
        }

        [[nodiscard]] bool is_alias() const
        {
            return alias;
        }

        [[nodiscard]] const std::optional<container::external_string> &get_original_name()
        const {
            return original_name;
        }

        [[nodiscard]] container::external_string get_alias()
        const {
            return alias_str;
        }

        template <typename T>
        T get()
        const {
            if constexpr (std::is_same_v<T, container::external_string>)
            {
                return default_str;
            }
            else if constexpr (std::is_same_v<T, int>)
            {
                return default_int;
            }
            else if constexpr (std::is_same_v<T, float>)
            {
                return default_float;
            }
            else if constexpr (std::is_same_v<T, bool>)
            {
                return default_bool;
            }
            else
            {
                throw except::exception("Unsupported type for default value");
            }
        }
    };
}