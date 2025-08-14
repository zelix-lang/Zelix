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
// Created by rodri on 8/13/25.
//

#include "analyzer.h"

#include "error.h"
using namespace zelix;

template <typename T>
bool check_duplicates(
    analyzer::result &res,
    ankerl::unordered_dense::map<
        container::external_string,
        analyzer::tool::ref,
        container::external_string_hash
    > &symbols,
    container::external_string &name,
    size_t &idx,
    T *const &duplicate
)
{
    if (symbols.contains(name))
    {
        if constexpr (std::is_same_v<T, code::function> || std::is_same_v<T, code::mod>)
        {
            // If the symbol already exists, add an error to the result
            res.errors.emplace_back(
                duplicate->line,
                duplicate->column,
                idx ,
                analyzer::error::REDEFINITION
            );
        }
        else
        {
            static_assert(false, "Unsupported type for duplicate check");
        }

        return true;
    }

    return false;
}

template <bool IsRoot>
void inject_context(
    ankerl::unordered_dense::map<
        container::external_string,
        analyzer::tool::ref,
        container::external_string_hash
    > &symbols,
    code::file_code *file,
    size_t &i,
    analyzer::result &res
)
{
    // Insert all modules and functions into the symbol table
    for (auto &[fst, snd] : file->modules)
    {
        // Honor visibility
        if constexpr (!IsRoot)
        {
            if (!snd->pub)
            {
                continue;
            }
        }

        check_duplicates(
            res,
            symbols,
            fst,
            i,
            snd
        );

        symbols.insert({fst, analyzer::tool::ref(snd)});
    }

    for (auto &[fst, snd] : file->functions)
    {
        // Honor visibility
        if constexpr (!IsRoot)
        {
            if (!snd->pub)
            {
                continue;
            }
        }

        check_duplicates(
            res,
            symbols,
            fst,
            i,
            snd
        );

        symbols.insert({fst, analyzer::tool::ref(snd)});
    }
}

analyzer::result analyzer::analyze(container::vector<code::file_code *> &files)
{
    result res;

    // Iterate over the files in reverse to process the files
    // without dependencies first
    const auto file_count = files.size();
    for (size_t i = file_count; i-- > 0;)
    {
        ankerl::unordered_dense::map<
            container::external_string,
            tool::ref,
            container::external_string_hash
        > symbols; // The symbol table for the current file

        // Get the current file
        auto *file = files[i];

        // Inject the current file's context into the symbol table
        inject_context<true>(symbols, file, i, res);

        // Inject all imports into the symbol table
        for (const auto &idx : file->imports)
        {
            // Get the imported file
            auto *imported_file = files[idx];
            inject_context<false>(symbols, imported_file, i, res);
        }
    }

    return res;
}