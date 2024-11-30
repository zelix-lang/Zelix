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

#include <iostream>
#include "panic.h"

void panic(const char* message) {
    std::cerr << "Panic: " << message << std::endl;
    exit(EXIT_FAILURE);
}