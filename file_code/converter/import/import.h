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

#pragma once
#include <ankerl/unordered_dense.h>
#include <fluent/str_has_prefix/str_has_prefix.h>
#include <zelix/container/external_string.h>

#include "../converter.h"
#include "fluent/ansi/ansi.h"
#include "global/constants.h"
#include "global/err/print.h"
#include "global/messages/import.h"
#include "lexer/lexer.h"
#include "memory/allocator.h"
#include "parser/ast.h"
#include "parser/parser.h"
#include "print_chain.h"
#include "time/timed_task.h"
#include "util/absolute.h"
#include "util/dirname.h"
#include "util/read_file.h"

namespace zelix::code::converter
{
    inline void imp(
        ankerl::unordered_dense::set<
            container::string,
            container::string_hash
        > &chain,
        container::vector<file_code *> &files,
        parser::ast *const &node,
        container::vector<queue_el> &queue,
        memory::lazy_allocator<parser::ast> &ast_allocator,
        memory::lazy_allocator<lexer::token> &token_allocator,
        memory::lazy_allocator<file_code> &file_allocator,
        container::string &root_dir,
        container::string &root_path
    )
    {
        time::post(node->value.get(), 3, 1);

        // Read the file
        container::string path;
        bool is_std = false;

        if (
            const auto &requested_path = node->value.get();
            str_has_prefix(requested_path.ptr(), "@std/")
        )
        {
            is_std = true;
            path.push(constants::stdlib.c_str(), constants::stdlib.size());
            path.push(
                requested_path.ptr() + 5, // Skip the "@std/" prefix
                requested_path.size() - 5 // Adjust the size accordingly
            );
            path.push(".zx"); // Add the .zx extension
        }
        else
        {
            if (!util::is_absolute(requested_path.ptr()))
            {
                // Add the root dir
                path.push(root_dir.c_str(), root_dir.size());
                path.push("/");
            }

            // Push the path
            path.push(requested_path.ptr(), requested_path.size());
        }

        // Check if the file is already in the chain
        if (chain.contains(path))
        {
            if (is_std)
            {
                time::complete(); // Complete the timed task for stdlib imports
                return;
            }

            time::fail("Circular import detected");

            // Report the error and print the details
            report::err::print(constants::import::circular_err, constants::import::circular_help);
            printf("\n" ANSI_BRIGHT_BLACK "Import chain:\n" ANSI_RESET);
            helper::print_import_chain(
                files,
                path,
                root_path
            );

            throw except::exception("");
        }

        chain.insert(path);
        time::advance();

        // Read the file
        auto contents = util::read_file(path.c_str());
        time::advance();

        // Lex the file
        auto tokens = lexer::lex(
            container::external_string(contents.c_str(), contents.size()),
            token_allocator
        );
        time::advance();

        // Parse the tokens
        auto ast = parser::parse(tokens, ast_allocator);

        // Allocate a new file_code object
        auto *file = file_allocator.alloc();
        file->content = container::move(contents);

        // Add the AST to the queue
        queue.emplace_back(ast, util::dirname(path.c_str()), container::move(file));
    }
}