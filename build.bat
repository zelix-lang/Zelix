@echo off
rem         ==== The Zelix Programming Language ====
rem    ---------------------------------------------------------
rem      - This file is part of the Zelix Programming Language
rem        codebase. Zelix is a fast, statically-typed and
rem        memory-safe programming language that aims to
rem        match native speeds while staying highly performant.
rem    ---------------------------------------------------------
rem      - Zelix is categorized as free software; you can
rem        redistribute it and/or modify it under the terms of
rem        the GNU General Public License as published by the
rem        Free Software Foundation, either version 3 of the
rem        License, or (at your option) any later version.
rem    ---------------------------------------------------------
rem      - Zelix is distributed in the hope that it will
rem        be useful, but WITHOUT ANY WARRANTY; without even
rem        the implied warranty of MERCHANTABILITY or FITNESS
rem        FOR A PARTICULAR PURPOSE. See the GNU General Public
rem        License for more details.
rem    ---------------------------------------------------------
rem      - You should have received a copy of the GNU General
rem        Public License along with Zelix. If not, see
rem        <https://www.gnu.org/licenses/>.

if not exist "cmake-build-debug" (
    mkdir cmake-build-debug
    cmake -B cmake-build-debug -S . -G "Ninja"
)

cmake --build cmake-build-debug --target Fluent -j 14