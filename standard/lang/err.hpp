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

#ifndef ERR_H
#define ERR_H

#include <string>

class Err {
private:
    std::string message;

public:
    Err(std::string message);
    std::string get_message() const;
};

#endif
