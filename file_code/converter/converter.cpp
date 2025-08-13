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

#include "converter.h"

#include "function.h"
#include "import.h"
#include "mod.h"
#include "util/dirname.h"
using namespace zelix;

container::vector<code::file_code *> code::convert(
    memory::lazy_allocator<file_code> &allocator,
    memory::lazy_allocator<parser::ast> &ast_allocator,
    memory::lazy_allocator<lexer::token> &token_allocator,
    memory::lazy_allocator<function> &fun_allocator,
    memory::lazy_allocator<mod> &mod_allocator,
    parser::ast *const &root,
    container::string &root_path
)
{
    container::vector<file_code *> files;
    container::vector<converter::queue_el> queue;
    ankerl::unordered_dense::set<
        container::string,
        container::string_hash
    > chain; // The import chain in a set for fast lookup
    queue.emplace_back(root, util::dirname(root_path.c_str()));

    // Process the queue until it's empty
    while (!queue.empty())
    {
        // Get the last node in the queue
        // We have to copy here since pop_back() calls the destructor
        auto [node, dir] = queue.back();
        queue.pop_back();

        // Allocate a new file_code object
        auto *file = allocator.alloc();

        // Walk the tree
        for (
            const auto &children = node->children;
            const auto &child: children
        )
        {
            switch (child->rule)
            {
                case parser::ast::IMPORT:
                {
                    converter::imp(
                        chain,
                        child,
                        queue,
                        ast_allocator,
                        token_allocator,
                        dir
                    );
                    break;
                }

                case parser::ast::FUNCTION:
                {
                    converter::function(
                        file->functions,
                        child,
                        fun_allocator
                    );

                    break;
                }

                case parser::ast::MOD:
                {
                    converter::mod(
                        file->modules,
                        child,
                        fun_allocator,
                        mod_allocator
                    );

                    break;
                }

                default:
                {
                    // Impossible case
                    break;
                }
            }
        }
    }

    return files; ///< Return the vector of file_code objects
}
