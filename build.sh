#    The Fluent Programming Language
#    -----------------------------------------------------
#    This code is released under the GNU GPL v3 license.
#    For more information, please visit:
#    https://www.gnu.org/licenses/gpl-3.0.html
#    -----------------------------------------------------
#    Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
#    This program comes with ABSOLUTELY NO WARRANTY.
#    For details type `fluent l`. This is free software,
#    and you are welcome to redistribute it under certain
#    conditions; type `fluent l -f` for details.

if [ ! -d "cmake-build-debug" ]; then
    mkdir cmake-build-debug
    cmake -B cmake-build-debug -S . -G "Ninja"
fi

cmake --build cmake-build-debug --target Fluent -j 14