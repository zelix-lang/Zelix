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

#include "result.hpp"

#include <iostream>

template<typename T>
std::ostream& operator<<(std::ostream& os, const Result<T>& res) {
    os << "Result(";

    if (res.has_error()) {
        os << "None, Err(" << res.get_error() << "))";
    } else {
        os << "Some(" << res.unwrap() << "), None)";
    }
    
    return os;
}