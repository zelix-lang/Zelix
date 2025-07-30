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
#include "ankerl/unordered_dense.h"
#include "container/external_string.h"
#include "value.h"

namespace fluent::cli
{
    class error
    {
    public:
        enum type
        {
            UNKNOWN,
            EXPECTED_VALUE,
            NOT_EXPECTED_VALUE,
            UNKNOWN_COMMAND,
            UNKNOWN_FLAG,
            TYPE_MISMATCH,
        };

        type error_type = UNKNOWN; ///< Type of the error
        size_t argv_pos = 0; ///< Position in the argv array where the error occurred
    };

    inline error global_error; ///< Global error object to store the last error

    class args
    {
        ankerl::unordered_dense::map<
            container::external_string,
            container::external_string,
            container::external_string_hash
        > str_args; ///< String arguments map

        ankerl::unordered_dense::map<
            container::external_string,
            int,
            container::external_string_hash
        > int_args; ///< Integer arguments map

        ankerl::unordered_dense::map<
            container::external_string,
            float,
            container::external_string_hash
        > float_args; ///< Float arguments map

        ankerl::unordered_dense::map<
            container::external_string,
            bool,
            container::external_string_hash
        > bool_args; ///< Boolean arguments map

    public:
        bool parse(
            const ankerl::unordered_dense::map<
                container::external_string,
                value,
                container::external_string_hash
            > &commands,

            const ankerl::unordered_dense::map<
                container::external_string,
                value,
                container::external_string_hash
            > &flags,

            // Aliases (cmd name -> alias)
            const ankerl::unordered_dense::map<
                container::external_string,
                container::external_string,
                container::external_string_hash
            > &cmd_aliases,

            const ankerl::unordered_dense::map<
                container::external_string,
                container::external_string,
                container::external_string_hash
            > &flag_aliases,

            // Aliases (alias -> cmd name)
            const ankerl::unordered_dense::map<
                container::external_string,
                container::external_string,
                container::external_string_hash
            > &cmd_aliases_reverse,

            const ankerl::unordered_dense::map<
                container::external_string,
                container::external_string,
                container::external_string_hash
            > &flag_aliases_reverse,

            const int argc,
            const char **argv
        )
        {
            if (argc < 2 || argv == nullptr)
            {
                global_error.argv_pos = 0;
                global_error.error_type = error::EXPECTED_VALUE;
                return false;
            }

            global_error.argv_pos = argc - 1;
            global_error.error_type = error::UNKNOWN_COMMAND;
            return false;
        }
    };
}