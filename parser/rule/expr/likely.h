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
// Created by rodrigo on 8/1/25.
//

#pragma once

namespace zelix::parser::rule::expr
{
#if defined(__x86_64__) || defined(_M_X64) || defined(__aarch64__)
    static constexpr uint64_t CALL_LIKELY = 0x1; // Likely a function call
    static constexpr uint64_t PROP_ACCESS_LIKELY = 0x2; // Likely a property access
    static constexpr uint64_t BOOLEAN_OP_LIKELY = 0x4; // Likely a boolean operation
    static constexpr uint64_t ARITHMETIC_OP_LIKELY = 0x8; // Likely an arithmetic operation
    static constexpr uint64_t ALL_LIKELY = CALL_LIKELY | PROP_ACCESS_LIKELY | BOOLEAN_OP_LIKELY | ARITHMETIC_OP_LIKELY; // All likely operations
#else
    static constexpr uint32_t CALL_LIKELY = 0x1; // Likely a function call
    static constexpr uint32_t PROP_ACCESS_LIKELY = 0x2; // Likely a property access
    static constexpr uint32_t BOOLEAN_OP_LIKELY = 0x4; // Likely a boolean operation
    static constexpr uint32_t ARITHMETIC_OP_LIKELY = 0x8; // Likely an arithmetic operation
    static constexpr uint32_t ALL_LIKELY = CALL_LIKELY | PROP_ACCESS_LIKELY | BOOLEAN_OP_LIKELY | ARITHMETIC_OP_LIKELY; // All likely operations
#endif
}