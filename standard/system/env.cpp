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

#include "env.hpp"

#include <string>
#include "../lang/result.hpp"
#include "../lang/err.hpp"

Result<std::string> get_env(const std::string* key) {
    char* value = getenv(key->c_str());

    if (value == NULL) {
        return Result(
            std::string(""), 
            std::optional<Err>(Err("Environment variable not found"))
        );
    }

    return Result(std::string(value), std::optional<Err>());
}

Result<bool> set_env(const std::string* key, const std::string* value) {
    if (setenv(key->c_str(), value->c_str(), 1) != 0) {
        return Result(
            false, 
            std::optional<Err>(Err("Failed to set environment variable"))
        );
    }

    return Result(true, std::optional<Err>());
}