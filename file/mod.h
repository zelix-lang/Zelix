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
// Created by rodri on 8/15/25.
//

#pragma once
#include "ankerl/unordered_dense.h"
#include "declaration.h"
#include "derivable.h"
#include "function.h"
#include "global/trace/trace.h"
#include "zelix/container/external_string.h"

namespace zelix::code
{
    struct mod : globals::trace, derivable
    {
        ankerl::unordered_dense::map<
            container::external_string,
            declaration,
            container::external_string_hash
        > declarations;

        ankerl::unordered_dense::map<
            container::external_string,
            function,
            container::external_string_hash
        > functions;
    };
}