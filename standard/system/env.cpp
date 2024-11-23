#include <string>
#include "../lang/result.h"

Result<std::string> get_env(const std::string* key) {
    char* value = getenv(key->c_str());

    if (value == NULL) {
        return Result(
            std::string(""), 
            optional<Err>(Err("Environment variable not found"))
        );
    }

    return Result(std::string(value), optional<Err>());
}

Result<bool> set_env(const std::string* key, const std::string* value) {
    if (setenv(key->c_str(), value->c_str(), 1) != 0) {
        return Result(
            false, 
            optional<Err>(Err("Failed to set environment variable"))
        );
    }

    return Result(true, optional<Err>());
}