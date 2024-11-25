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

#ifndef FILES_H
#define FILES_H

#include <iostream>
#include <string>
#include <vector>
#include "../lang/result.hpp"
#include "../lang/err.h"

Result<bool> write_file(const std::string* path, const std::string* content);
Result<std::string> read_file(const std::string* path);
Result<bool> delete_file(const std::string* path);
Result<std::vector<std::string>> delete_dir(const std::string* path);
Result<std::vector<std::string>> walk_dir(const std::string* path);

#endif