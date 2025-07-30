@echo off
rem The Fluent Programming Language
rem -----------------------------------------------------
rem This code is released under the GNU GPL v3 license.
rem For more information, please visit:
rem https://www.gnu.org/licenses/gpl-3.0.html
rem -----------------------------------------------------
rem Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
rem This program comes with ABSOLUTELY NO WARRANTY.
rem For details type `fluent l`. This is free software,
rem and you are welcome to redistribute it under certain
rem conditions; type `fluent l -f` for details.

if not exist "cmake-build-debug" (
    mkdir cmake-build-debug
    cmake -B cmake-build-debug -S . -G "Ninja"
)

cmake --build cmake-build-debug --target Fluent -j 14