/*
    These files are part of the Surf's standard library.
    They're bundled with the main program by the compiler
    which is then converted to machine code.

    -----
    License notice:

    This code is released under the GNU GPL v3 license.
    The code is provided as is without any warranty
    Copyright (c) 2024 Rodrigo R. & all Surf contributors
*/

#include "read_line.hpp"
#include <iostream>
#include <string>

std::string read_line() {
    std::string line;

    getline(std::cin, line);
    return line;
}