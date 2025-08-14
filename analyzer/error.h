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

#pragma once

namespace zelix::analyzer::error
{
    inline constexpr auto NAME_SHOULD_BE_SNAKE_CASE     = 0;  // E0001
    inline constexpr auto PARAM_TYPE_NOTHING            = 1;  // E0002
    inline constexpr auto DATA_OUTLIVES_STACK           = 2;  // E0003
    inline constexpr auto MUST_RETURN_A_VALUE           = 3;  // E0004
    inline constexpr auto UNDEFINED_REFERENCE           = 4;  // E0005
    inline constexpr auto UNUSED_VARIABLE               = 5;  // E0006
    inline constexpr auto REDEFINITION                  = 6;  // E0007
    inline constexpr auto INVALID_DEREFERENCE           = 7;  // E0008
    inline constexpr auto TYPE_MISMATCH                 = 8;  // E0009
    inline constexpr auto PARAMETER_COUNT_MISMATCH      = 9;  // E0010
    inline constexpr auto CANNOT_INFER_TYPE             = 10; // E0011
    inline constexpr auto SHOULD_NOT_RETURN             = 11; // E0012
    inline constexpr auto CANNOT_TAKE_ADDRESS           = 12; // E0013
    inline constexpr auto INVALID_PROP_ACCESS           = 13; // E0014
    inline constexpr auto ILLEGAL_PROP_ACCESS           = 14; // E0015
    inline constexpr auto CONSTANT_REASSIGNMENT         = 15; // E0016
    inline constexpr auto DOES_NOT_HAVE_CONSTRUCTOR     = 16; // E0017
    inline constexpr auto SHOULD_NOT_HAVE_GENERICS      = 17; // E0018
    inline constexpr auto VALUE_NOT_ASSIGNED            = 18; // E0019
    inline constexpr auto CIRCULAR_MODULE_DEPENDENCY    = 19; // E0020
    inline constexpr auto SELF_REFERENCE                = 20; // E0021
    inline constexpr auto INVALID_LOOP_INSTRUCTION      = 21; // E0022
    inline constexpr auto INVALID_POINTER               = 22; // E0023
    inline constexpr auto NO_MAIN_FUNCTION              = 23; // E0024
    inline constexpr auto MAIN_FUNCTION_HAS_PARAMETERS  = 24; // E0025
    inline constexpr auto MAIN_FUNCTION_HAS_RETURN      = 25; // E0026
    inline constexpr auto MAIN_FUNCTION_HAS_GENERICS    = 26; // E0027
    inline constexpr auto MAIN_FUNCTION_IS_PUBLIC       = 27; // E0028
}