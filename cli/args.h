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
#include <charconv>
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
            value,
            container::external_string_hash
        > &commands;

        ankerl::unordered_dense::map<
            container::external_string,
            value,
            container::external_string_hash
        > &flags;

        // Aliases (cmd name -> alias)
        ankerl::unordered_dense::map<
            container::external_string,
            container::external_string,
            container::external_string_hash
        > &cmd_aliases;

        ankerl::unordered_dense::map<
            container::external_string,
            container::external_string,
            container::external_string_hash
        > &flag_aliases;

        // Aliases (alias -> cmd name)
        ankerl::unordered_dense::map<
            container::external_string,
            container::external_string,
            container::external_string_hash
        > &cmd_aliases_reverse;

        ankerl::unordered_dense::map<
            container::external_string,
            container::external_string,
            container::external_string_hash
        > &flag_aliases_reverse;

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

        ankerl::unordered_dense::map<
            container::external_string,
            container::external_string,
            container::external_string_hash
        > str_flags; ///< String flags map

        ankerl::unordered_dense::map<
            container::external_string,
            int,
            container::external_string_hash
        > int_flags; ///< Integer flags map

        ankerl::unordered_dense::map<
            container::external_string,
            float,
            container::external_string_hash
        > float_flags; ///< Float flags map

        ankerl::unordered_dense::map<
            container::external_string,
            bool,
            container::external_string_hash
        > bool_flags; ///< Boolean flags map

        container::external_string cmd;

        bool parse_flag(
            container::external_string &flag,
            bool &waiting_value,
            value::type &expected,
            const int i
        )
        {
            // Handle aliases
            if (flag_aliases_reverse.contains(flag))
            {
                const auto &alias = flag_aliases_reverse.at(flag);
                flag = alias;
            }

            // Get the flag from the map
            if (flags.contains(flag))
            {
                const auto &flag_val = flags.at(flag);
                waiting_value = flag_val.get_type() != value::BOOL;
                expected = flag_val.get_type();

                if (!waiting_value)
                {
                    bool_flags[flag] = true;
                }
            }
            else
            {
                global_error.error_type = error::UNKNOWN_FLAG;
                global_error.argv_pos = i;
                return false;
            }

            return true;
        }

        template <typename T, typename Flag>
        bool parse_value(
            const container::external_string &value,
            container::external_string &name
        )
        {
            if constexpr (std::is_same_v<T, container::external_string>)
            {
                if constexpr (std::is_same_v<Flag, bool>)
                {
                    str_flags[name] = value;
                }
                else
                {
                    str_args[name] = value;
                }

                return true;
            }
            else if constexpr (std::is_same_v<T, bool>)
            {
                if constexpr (std::is_same_v<Flag, bool>)
                {
                    if (memcmp(value.ptr(), "true", value.size()) == 0)
                    {
                        bool_flags[name] = true;
                        return true;
                    }

                    if (memcmp(value.ptr(), "false", value.size()) == 0)
                    {
                        bool_flags[name] = false;
                        return true;
                    }
                }
                else
                {
                    if (memcmp(value.ptr(), "true", value.size()) == 0)
                    {
                        bool_args[name] = true;
                        return true;
                    }

                    if (memcmp(value.ptr(), "false", value.size()) == 0)
                    {
                        bool_args[name] = false;
                        return true;
                    }
                }

                return false; // Invalid boolean value
            }
            else if constexpr (std::is_same_v<T, int>)
            {
                int result = 0;
                const auto value_ptr = value.ptr();

                // Iterate over the ptr
                for (size_t i = 0; i < value.size(); ++i)
                {
                    const char c = value_ptr[i];
                    if (c < '0' || c > '9')
                    {
                        return false; // Invalid integer value
                    }

                    result = result * 10 + (c - '0');
                }

                if constexpr (std::is_same_v<Flag, bool>)
                {
                    int_flags[name] = result;
                }
                else
                {
                    int_args[name] = result;
                }

                return true;
            }
            else if constexpr (std::is_same_v<T, float>)
            {
                float result = 0.0f;
                const auto value_ptr = value.ptr();

                if (
                    auto [ptr, ec] = std::from_chars(value_ptr, value_ptr + value.size(), result);
                    ec != std::errc()
                ) {
                    return false; // Invalid float value
                }

                if constexpr (std::is_same_v<Flag, bool>)
                {
                    float_flags[name] = result;
                }
                else
                {
                    float_args[name] = result;
                }

                return true;
            }
            else
            {
                static_assert(
                    false,
                    "Unsupported type for value parsing"
                );

                return false; // Should never reach here
            }
        }

        template <typename T, typename Flag>
        T val(const container::external_string &name)
        {
            if constexpr (
                std::is_same_v<T, container::external_string>
                || std::is_same_v<T, const char *>
            )
            {
                if constexpr (std::is_same_v<Flag, bool>)
                {
                    // Get from the flags
                    if (!str_flags.contains(name))
                    {
                        // Return default value
                        const auto &def = flags.at(name);
                        return def.get<T>();
                    }

                    return str_flags.at(name);
                }
                else
                {
                    // Get from the commands
                    if (!str_args.contains(name))
                    {
                        // Return default value
                        const auto &def = commands.at(name);
                        return def.get<T>();
                    }

                    return str_flags.at(name);
                }
            }
            else if constexpr (std::is_same_v<T, int>)
            {
                if constexpr (std::is_same_v<Flag, bool>)
                {
                    // Get from the flags
                    if (!int_flags.contains(name))
                    {
                        // Return default value
                        const auto &def = flags.at(name);
                        return def.get<T>();
                    }

                    return int_flags.at(name);
                }
                else
                {
                    // Get from the commands
                    if (!int_args.contains(name))
                    {
                        // Return default value
                        const auto &def = commands.at(name);
                        return def.get<T>();
                    }

                    return int_flags.at(name);
                }
            }
            else if constexpr (std::is_same_v<T, float>)
            {
                if constexpr (std::is_same_v<Flag, bool>)
                {
                    // Get from the flags
                    if (!float_flags.contains(name))
                    {
                        // Return default value
                        const auto &def = flags.at(name);
                        return def.get<T>();
                    }

                    return float_flags.at(name);
                }
                else
                {
                    // Get from the commands
                    if (!float_args.contains(name))
                    {
                        // Return default value
                        const auto &def = commands.at(name);
                        return def.get<T>();
                    }

                    return float_args.at(name);
                }
            }
            else if constexpr (std::is_same_v<T, bool>)
            {
                if constexpr (std::is_same_v<Flag, float>)
                {
                    // Get from the flags
                    if (!bool_flags.contains(name))
                    {
                        // Return default value
                        const auto &def = flags.at(name);
                        return def.get<T>();
                    }

                    return bool_flags.at(name);
                }
                else
                {
                    // Get from the commands
                    if (!bool_args.contains(name))
                    {
                        // Return default value
                        const auto &def = commands.at(name);
                        return def.get<T>();
                    }

                    return bool_args.at(name);
                }
            }
            else
            {
                static_assert(
                    false,
                    "Invalid type"
                );

                // Unreachable code
                return (T){0};
            }
        }

    public:
        explicit args(
            ankerl::unordered_dense::map<
                container::external_string,
                value,
                container::external_string_hash
            > &commands,

            ankerl::unordered_dense::map<
                container::external_string,
                value,
                container::external_string_hash
            > &flags,

            // Aliases (cmd name -> alias)
            ankerl::unordered_dense::map<
                container::external_string,
                container::external_string,
                container::external_string_hash
            > &cmd_aliases,


            ankerl::unordered_dense::map<
                container::external_string,
                container::external_string,
                container::external_string_hash
            > &flag_aliases,

            // Aliases (alias -> cmd name)
            ankerl::unordered_dense::map<
                container::external_string,
                container::external_string,
                container::external_string_hash
            > &cmd_aliases_reverse,

            ankerl::unordered_dense::map<
                container::external_string,
                container::external_string,
                container::external_string_hash
            > &flag_aliases_reverse
        ) :
            commands(commands),
            flags(flags),
            cmd_aliases(cmd_aliases),
            flag_aliases(flag_aliases),
            cmd_aliases_reverse(cmd_aliases_reverse),
            flag_aliases_reverse(flag_aliases_reverse)
        {}

        bool parse(
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

            // Whether we are waiting for a value
            bool waiting_value = false;
            bool value_command = false; ///< Whether the expected value is for a command
            bool has_command = false;
            value::type expected;
            container::external_string flag;
            for (int i = 1; i < argc; ++i)
            {
                const auto arg = argv[i];
                if (waiting_value)
                {
                    auto val = container::external_string(arg);
                    bool parsing_success = false;

                    if (value_command)
                    {
                        switch (expected)
                        {
                            case value::BOOL:
                            {
                                parsing_success = parse_value<bool, int>(
                                    val,
                                    flag
                                );

                                break;
                            }

                            case value::FLOAT:
                            {
                                parsing_success = parse_value<float, int>(
                                    val,
                                    flag
                                );

                                break;
                            }

                            case value::INTEGER:
                            {
                                parsing_success = parse_value<int, int>(
                                    val,
                                    flag
                                );

                                break;
                            }

                            case value::STRING:
                            {
                                parsing_success = parse_value<container::external_string, int>(
                                    val,
                                    flag
                                );

                                break;
                            }
                        }
                    }
                    else
                    {
                        switch (expected)
                        {
                            case value::BOOL:
                            {
                                parsing_success = parse_value<bool, bool>(
                                    val,
                                    flag
                                );

                                break;
                            }

                            case value::FLOAT:
                            {
                                parsing_success = parse_value<float, bool>(
                                    val,
                                    flag
                                );

                                break;
                            }

                            case value::INTEGER:
                            {
                                parsing_success = parse_value<int, bool>(
                                    val,
                                    flag
                                );

                                break;
                            }

                            case value::STRING:
                            {
                                parsing_success = parse_value<container::external_string, bool>(
                                    val,
                                    flag
                                );

                                break;
                            }
                        }
                    }

                    if (!parsing_success)
                    {
                        global_error.error_type = error::TYPE_MISMATCH;
                        global_error.argv_pos = i;
                        return false;
                    }

                    waiting_value = false;
                    value_command = false; // Reset the command value expectation
                    continue;
                }

                // Check if we have a flag
                if (arg[0] == '-')
                {
                    // Check if the flag is valid
                    if (arg[1] == '\0')
                    {
                        global_error.error_type = error::UNKNOWN_FLAG;
                        global_error.argv_pos = i;
                        return false;
                    }

                    // Check if we have a long flag
                    if (arg[1] == '-')
                    {
                        // Make sure we have a name
                        if (arg[2] == '\0')
                        {
                            global_error.error_type = error::EXPECTED_VALUE;
                            global_error.argv_pos = i;
                        }

                        flag = container::external_string(arg + 2);
                    }
                    else
                    {
                        flag = container::external_string(arg + 1);
                    }

                    // Check if we have a value
                    const auto flag_ptr = flag.ptr();
                    if (
                        const auto equals = strchr(flag_ptr, '=');
                        equals != nullptr
                    )
                    {
                        // Find the position of the equals sign
                        const size_t equals_pos = equals - flag.ptr();
                        const auto value = flag_ptr + equals_pos + 1; // Skip the equals sign

                        flag.set_size(equals_pos); // Exclude the equals sign

                        if (!parse_flag(flag, waiting_value, expected, i))
                        {
                            return false;
                        }

                        // Check if we are expecting a value
                        if (!waiting_value)
                        {
                            global_error.error_type = error::NOT_EXPECTED_VALUE;
                            global_error.argv_pos = i;
                            return false;
                        }

                        // Validate the value
                        if (value[0] == '\0')
                        {
                            global_error.error_type = error::EXPECTED_VALUE;
                            global_error.argv_pos = i;
                            return false;
                        }

                        auto val = container::external_string(value);
                        waiting_value = false; // We are no longer waiting for a value
                        bool parsing_success = false;

                        // Parse the value
                        switch (expected)
                        {
                            case value::BOOL:
                            {
                                parsing_success = parse_value<bool, bool>(
                                    val,
                                    flag
                                );

                                break;
                            }

                            case value::FLOAT:
                            {
                                parsing_success = parse_value<float, bool>(
                                    val,
                                    flag
                                );

                                break;
                            }

                            case value::INTEGER:
                            {
                                parsing_success = parse_value<int, bool>(
                                    val,
                                    flag
                                );

                                break;
                            }

                            case value::STRING:
                            {
                                parsing_success = parse_value<container::external_string, bool>(
                                    val,
                                    flag
                                );

                                break;
                            }
                        }

                        if (!parsing_success)
                        {
                            global_error.error_type = error::TYPE_MISMATCH;
                            global_error.argv_pos = i;
                            return false;
                        }

                        continue;
                    }

                    // Parse the flag
                    if (!parse_flag(flag, waiting_value, expected, i))
                    {
                        return false;
                    }

                    continue;
                }

                // Parse commands
                if (!has_command)
                {
                    cmd = container::external_string(arg);
                    has_command = true;

                    // Handle aliases
                    if (cmd_aliases_reverse.contains(cmd))
                    {
                        const auto &alias = cmd_aliases_reverse.at(cmd);
                        cmd = alias;
                    }

                    // Check if the command is valid
                    if (commands.contains(cmd))
                    {
                        const auto &cmd_val = commands.at(cmd);
                        waiting_value = cmd_val.get_type() != value::BOOL;
                        expected = cmd_val.get_type();

                        if (!waiting_value)
                        {
                            str_args[cmd] = cmd_val.get<container::external_string>();
                        }
                        else
                        {
                            value_command = true; // We are expecting a value for the command
                        }
                    }
                    else
                    {
                        global_error.error_type = error::UNKNOWN_COMMAND;
                        global_error.argv_pos = i;
                        return false;
                    }

                    continue;
                }

                // Invalid argument
                global_error.error_type = error::NOT_EXPECTED_VALUE;
                global_error.argv_pos = i;
                return false;
            }

            // Make sure we are not waiting for a value
            if (!has_command)
            {
                global_error.error_type = error::EXPECTED_VALUE;
                global_error.argv_pos = argc - 1;
                return false;
            }

            // Default values if we are still expecting a value
            if (waiting_value)
            {
                // Retrieve the value
                const value &val = value_command
                    ? commands.at(cmd)
                    : flags.at(flag);


                switch (val.get_type())
                {
                    case value::STRING:
                    {
                        if (value_command)
                        {
                            str_args[cmd] = val.get<container::external_string>();
                        }
                        else
                        {
                            str_flags[flag] = val.get<container::external_string>();
                        }
                    }

                    case value::BOOL:
                    {
                        if (value_command)
                        {
                            bool_args[cmd] = val.get<bool>();
                        }
                        else
                        {
                            bool_flags[flag] = val.get<bool>();
                        }
                    }

                    case value::FLOAT:
                    {
                        if (value_command)
                        {
                            float_args[cmd] = val.get<float>();
                        }
                        else
                        {
                            float_flags[flag] = val.get<float>();
                        }
                    }

                    case value::INTEGER:
                    {
                        if (value_command)
                        {
                            int_args[cmd] = val.get<int>();
                        }
                        else
                        {
                            int_flags[flag] = val.get<int>();
                        }
                    }
                }
            }

            return true;
        }

        template <typename T>
        T flag(const container::external_string &name)
        {
            return val<T, bool>(name);
        }

        template <typename T>
        T flag(const char *name)
        {
            return val<T, bool>(container::external_string(name));
        }

        template <typename T>
        T command(const container::external_string &name)
        {
            return val<T, int>(name);
        }

        template <typename T>
        T command(const char *name)
        {
            return val<T, int>(container::external_string(name));
        }

        container::external_string &get_cmd()
        {
            return this->cmd;
        }

        [[nodiscard]] static bool is_err()

        {
            return global_error.error_type != error::UNKNOWN;
        }
    };
}