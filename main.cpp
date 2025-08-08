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

#include "command/compile.h"
#include "zelix/cli/app.h"
using namespace zelix;

#define APP_NAME "The Zelix Programming Language"
#define APP_DESC "A blazingly fast programming language"

int main(const int argc, const char **argv)
{
    const auto *compile = "compile";
    const auto *run = "run";

    cli::app app(
        APP_NAME,
        APP_DESC,
        argc,
        argv
    );

    app.command<const char*>(
        compile,
        "c",
        "compiles a Zelix project",
        "."
    );

    app.command<const char*>(
        run,
        "r",
        "runs a Zelix project",
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
    else
    {
        // Print the header
        printf(ANSI_BOLD_BRIGHT_BLUE APP_NAME ANSI_RESET "\n");
        printf(ANSI_BRIGHT_BLACK APP_DESC ANSI_RESET "\n\n");
    }

    // Get the command
    const auto &cmd = args.get_cmd();

    if (
        const auto cmd_ptr = cmd.ptr();
        cmd_ptr == compile
    )
    {
        command::compile(args);
    } else if (cmd_ptr == run)
    {

    }
    return 0;
}