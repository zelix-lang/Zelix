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
#include "../lang/result.hpp"
#include "../lang/err.hpp"
#include <iostream>
#include <string>
#include <optional>

Result<std::string> read_line() {
    try {
        std::string line;

        getline(std::cin, line);
        return Result(line, std::optional<Err>());
    } catch(const std::exception& e) {
        return Result(std::string(""), std::optional<Err>(Err(e.what())));
    }
    
}