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

//
// Created by rodri on 8/12/25.
//

#include "print_chain.h"

#include "fluent/ansi/ansi.h"
#include "util/nested_spaces.h"
using namespace zelix;

void print_chain_el(
    container::string &el,
    const container::string &target,
    const size_t nesting = 0
)
{
    if (nesting > 0)
    {
        util::print_nested_spaces(nesting);
        printf(ANSI_BRIGHT_BLACK "└─ " ANSI_RESET);
    }

    if (el == target)
    {
        printf(ANSI_RED "%s" ANSI_RESET "\n", el.c_str());
    }
    else
    {
        printf(ANSI_BRIGHT_BLACK "%s" ANSI_RESET "\n", el.c_str());
    }
}

void code::converter::helper::print_import_chain(
    container::vector<file_code *> &files,
    const container::string &target,
    container::string &root_path
)
{
    if (files.empty())
    {
        return;
    }

    // Print the root path
    print_chain_el(root_path, target);

    // Print the import chain
    for (size_t i = 0; i < files.size(); i++)
    {
        print_chain_el(files[i]->path, target, i);
    }
}