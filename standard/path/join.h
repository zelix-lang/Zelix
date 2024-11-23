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

#ifndef JOIN_PATH_H
#define JOIN_PATH_H

#include <string>

std::string get_path_separator();
std::string join_path(const std::string* path1, const std::string* path2);

#endif