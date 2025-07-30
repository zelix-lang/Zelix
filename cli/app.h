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

#include <ankerl/unordered_dense.h>
#include <fluent/ansi/ansi.h>

#include "args.h"
#include "container/external_string.h"
#include "container/owned_string.h"
#include "except/exception.h"
#include "value.h"
#define ANSI_UNDERLINE "\033[4m"
#define ANSI_DIM "\e[2m"
#define ANSI_DIM_END "\e[22m"

namespace fluent::cli
{
    class app
    {
        const char *name_ = nullptr;
        const char *desc_ = nullptr;

        ankerl::unordered_dense::map<
            container::external_string,
            value,
            container::external_string_hash
        > commands;

        ankerl::unordered_dense::map<
            container::external_string,
            value,
            container::external_string_hash
        > flags;

        // Aliases (cmd name -> alias)
        ankerl::unordered_dense::map<
            container::external_string,
            container::external_string,
            container::external_string_hash
        > cmd_aliases;

        ankerl::unordered_dense::map<
            container::external_string,
            container::external_string,
            container::external_string_hash
        > flag_aliases;

        // Aliases (alias -> cmd name)
        ankerl::unordered_dense::map<
            container::external_string,
            container::external_string,
            container::external_string_hash
        > cmd_aliases_reverse;

        ankerl::unordered_dense::map<
            container::external_string,
            container::external_string,
            container::external_string_hash
        > flag_aliases_reverse;

        const int argc;
        const char **argv;

        void write_val_info(
            container::string &msg,
            const container::external_string &name,
            const value &val,
            const bool flag,
            const int alias_padding = 0
        )
        {
            const auto &desc = val.get_description();
            const auto &alias = flag ? flag_aliases[name] : cmd_aliases[name];

            msg.push(ANSI_RESET, 4);
            msg.push(ANSI_BRIGHT_BLACK, 5);
            msg.push(" ~ ", 3);
            msg.push(desc.ptr(), desc.size());
            msg.push("\n   ", 4);

            for (size_t i = 0; i < name.size() + 6 + alias_padding + alias.size(); ++i)
            {
                msg.push(' ');
            }

            switch (val.get_type())
            {
                case value::STRING:
                {
                    msg.push("[type=str, default=", 19);
                    const auto str_val = val.get<container::external_string>();
                    msg.push(str_val.ptr(), str_val.size());
                    break;
                }

                case value::BOOL:
                {
                    msg.push("[type=bool, default=", 20);
                    if (val.get<bool>())
                    {
                        msg.push("true", 4);
                    }
                    else
                    {
                        msg.push("false", 5);
                    }

                    break;
                }

                case value::FLOAT:
                {
                    msg.push("[type=float, default=", 21);
                    const std::string float_str = std::to_string(val.get<float>());
                    msg.push(float_str.c_str(), float_str.size());
                    break;
                }

                case value::INTEGER:
                {
                    msg.push("[type=int, default=", 19);
                    const std::string int_str = std::to_string(val.get<int>());
                    msg.push(int_str.c_str(), int_str.size());
                    break;
                }
            }

            msg.push("]" ANSI_RESET "\n", 6);
        }
    public:
        explicit app(
            const char *name,
            const char *description,
            const int argc,
            const char **argv
        )
            : name_(name), desc_(description), argc(argc), argv(argv)
        {
            if (name == nullptr || description == nullptr)
            {
                throw except::exception("Name and description cannot be null");
            }
        }

        [[nodiscard]] const char *name() const
        {
            return name_;
        }

        [[nodiscard]] const char *description() const
        {
            return desc_;
        }

        template <typename T>
        void command(
            const container::external_string &name,
            const container::external_string &alias,
            const container::external_string &description,
            const T &def
        )
        {
            if (commands.contains(name) || cmd_aliases_reverse.contains(alias))
            {
                throw except::exception("Command already exists");
            }

            cmd_aliases[name] = alias;
            cmd_aliases_reverse[alias] = name;
            commands[name] = value(def, description);
        }

        template <typename T>
        void command(
            const char *name,
            const char *alias,
            const char *description,
            const T &value
        )
        {
            command(
                container::external_string(name, strlen(name)),
                container::external_string(alias, strlen(alias)),
                container::external_string(description, strlen(description)),
                value
            );
        }

        template <typename T>
        void flag(
            const container::external_string &name,
            const container::external_string &alias,
            const container::external_string &description,
            const T &def
        )
        {
            if (flags.contains(name) || flag_aliases_reverse.contains(alias))
            {
                throw except::exception("Flag already exists");
            }

            flag_aliases[name] = alias;
            flag_aliases_reverse[alias] = name;
            flags[name] = value(def, description);
        }

        template <typename T>
        void flag(
            const char *name,
            const char *alias,
            const char *description,
            const T &value
        )
        {
            flag(
                container::external_string(name, strlen(name)),
                container::external_string(alias, strlen(alias)),
                container::external_string(description, strlen(description)),
                value
            );
        }

        std::optional<args> parse()
        {
            args parsed_args(
                commands,
                flags,
                cmd_aliases,
                flag_aliases,
                cmd_aliases_reverse,
                flag_aliases_reverse
            );

            if (
                !parsed_args.parse(argc, argv)
            )
            {
                return std::nullopt;
            }

            return std::make_optional<args>(std::move(parsed_args));
        }

        [[nodiscard]] container::string help()
        {
            container::string msg;
            msg.push(ANSI_BOLD_BRIGHT_BLUE, 7);
            msg.push(name_);
            msg.push(ANSI_RESET, 4);
            msg.push("\n", 1);
            msg.push(ANSI_BRIGHT_BLACK, 5);
            msg.push(desc_);
            msg.push(ANSI_RESET, 4);
            msg.push("\n\n", 2);

            const size_t bin_len = strlen(argv[0]) - 1;
            if (global_error.error_type != error::UNKNOWN)
            {
                msg.push(ANSI_BOLD_BRIGHT_RED, 7);
                msg.push("Error: ", 7);
                msg.push(ANSI_RESET, 4);
                msg.push(ANSI_BRIGHT_RED, 5);

                switch (global_error.error_type)
                {
                    case error::EXPECTED_VALUE:
                        msg.push("Expected a value", 16);
                        break;

                    case error::NOT_EXPECTED_VALUE:
                        msg.push("Unexpected value", 16);
                        break;

                    case error::TYPE_MISMATCH:
                        msg.push("Type mismatch", 13);
                        break;

                    case error::UNKNOWN_COMMAND:
                        msg.push("Unknown command", 15);
                        break;

                    case error::UNKNOWN_FLAG:
                        msg.push("Unknown flag", 12);
                        break;

                    default:
                        break;
                }

                msg.push("\n  ", 3);
                msg.push(ANSI_RESET, 4);
                msg.push(ANSI_BRIGHT_BLACK, 5);
                msg.push("➤ ");
                msg.push(ANSI_RESET, 4);

                const bool error_first = global_error.argv_pos == 0;
                if (error_first)
                {
                    msg.push(ANSI_BRIGHT_RED, 5);
                    msg.push(ANSI_UNDERLINE, 4);
                }
                else
                {
                    msg.push(ANSI_BOLD_BRIGHT_BLUE, 7);
                }

                msg.push(argv[0], bin_len + 1);

                if (!error_first || argc > 1)
                {
                    msg.push(ANSI_RESET, 4);
                    msg.push(ANSI_BRIGHT_BLACK, 5);
                    msg.push(" ... ");
                    msg.push(ANSI_RESET, 4);
                }

                if (!error_first)
                {
                    msg.push(ANSI_BRIGHT_RED, 5);
                    msg.push(ANSI_UNDERLINE, 4);
                    msg.push(argv[global_error.argv_pos]);
                    msg.push(ANSI_RESET "\n         ", 14);

                    for (size_t i = 0; i < bin_len + 1; ++i)
                    {
                        msg.push(' ');
                    }
                }
                else
                {
                    msg.push(ANSI_RESET "\n    ", 9);
                }

                msg.push(ANSI_GREEN, 5);
                msg.push("⤷");
                msg.push(" help: ", 7);

                switch (global_error.error_type)
                {
                    case error::EXPECTED_VALUE:
                        msg.push("add a value after this", 22);
                        break;

                    case error::NOT_EXPECTED_VALUE:
                        msg.push("remove this", 11);
                        break;

                    case error::TYPE_MISMATCH:
                        msg.push("change the value to match the expected type", 43);
                        break;

                    case error::UNKNOWN_COMMAND:
                        msg.push("use --help to see a list of commands", 36);
                        break;

                    case error::UNKNOWN_FLAG:
                        msg.push("use --help to see a list of flags", 33);
                        break;

                    default:
                        break;
                }

                msg.push(ANSI_RESET, 4);
                msg.push(ANSI_RESET, 4);
                msg.push("\n\n", 2);
            }

            msg.push(
                ANSI_BRIGHT_YELLOW
                "Usage:\n"
                ANSI_RESET
                ANSI_BOLD_BRIGHT_BLUE,
                23
            );

            msg.push(argv[0], bin_len + 1);

            msg.push(
                " "
                ANSI_RESET
                ANSI_UNDERLINE
                ANSI_BRIGHT_PURPLE
                "[--flags]"
                ANSI_RESET
                " "
                ANSI_UNDERLINE
                ANSI_BRIGHT_GREEN
                "<command>"
                ANSI_RESET
                " "
                ANSI_YELLOW
                ANSI_UNDERLINE
                "[<args>]"
                ANSI_RESET
                " "
                ANSI_UNDERLINE
                ANSI_BRIGHT_PURPLE
                "[--flags]"
                ANSI_RESET
                "\n",
                96
            );

            msg.push(ANSI_DIM, 4);
            msg.push(ANSI_BRIGHT_BLACK, 5);

            for (size_t i = 0; i < bin_len; ++i)
            {
                msg.push('-');
            }

            msg.push('>');
            msg.push(ANSI_DIM_END, 5);
            msg.push(ANSI_RESET, 4);

            msg.push(
                " "
                ANSI_DIM
                ANSI_BRIGHT_PURPLE
                "optional"
                ANSI_DIM_END
                ANSI_RESET
                ANSI_DIM
                ANSI_BRIGHT_GREEN
                "  required"
                ANSI_DIM_END
                ANSI_RESET
                ANSI_DIM
                ANSI_YELLOW
                "  optional"
                ANSI_DIM_END
                ANSI_RESET
                " "
                ANSI_DIM
                ANSI_BRIGHT_PURPLE
                "optional"
                ANSI_DIM_END
                ANSI_RESET
                "\n\n",
                112
            );

            if (!commands.empty())
            {
                msg.push(
                    ANSI_YELLOW
                    "Available commands:\n"
                    ANSI_RESET,
                    29
                );
            }

            for (const auto &[name, val] : commands)
            {
                msg.push(ANSI_BRIGHT_BLACK, 5);
                msg.push("  ", 2);
                msg.push("➤ ");
                msg.push(ANSI_RESET, 4);
                msg.push(ANSI_CYAN, 5);
                msg.push(name.ptr(), name.size());

                // Honor aliases
                const auto &alias = cmd_aliases[name];
                msg.push(", ", 2);
                msg.push(alias.ptr(), alias.size());

                // Write the value's info
                write_val_info(msg, name, val, false);
            }

            if (!flags.empty())
            {
                msg.push(
                    ANSI_YELLOW
                    "Available flags:\n"
                    ANSI_RESET,
                    26
                );
            }

            for (auto &[name, val] : flags)
            {
                msg.push("  ", 2);
                msg.push(ANSI_BRIGHT_BLUE, 5);
                msg.push("--", 2);
                msg.push(name.ptr(), name.size());

                // Honor aliases
                const auto &alias = flag_aliases[name];
                msg.push(", -", 3);
                msg.push(alias.ptr(), alias.size());

                write_val_info(msg, name, val, true, 1);
            }

            return msg;
        }
    };
}