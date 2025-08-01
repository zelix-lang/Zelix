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

#include "fluent/cli/app.h"
using namespace fluent;

int main(const int argc, const char **argv)
{
    const auto *compile = "compile";
    const auto *run = "run";

    cli::app app(
        "The Fluent Programming Language",
        "A blazingly fast programming language",
        argc,
        argv
    );

    app.command<const char*>(
        compile,
        "c",
        "compiles a Fluent project",
        "."
    );

    app.command<const char*>(
        run,
        "r",
        "runs a Fluent project",
        "."
    );

    app.flag<int>(
        "optimization",
        "O",
        "specifies the optimization level",
        3
    );

    auto args = app.parse();
    if (cli::args::is_err())
    {
        auto help = app.help();
        printf("%s", help.c_str());
    }

    // Get the command
    const auto &cmd = args.get_cmd();

    if (
        const auto cmd_ptr = cmd.ptr();
        cmd_ptr == compile
    )
    {

    } else if (cmd_ptr == run)
    {

    }
    return 0;
}