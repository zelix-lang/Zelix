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

#ifndef RESULT_H
#define RESULT_H

#include "err.hpp"
#include "panic.h"
#include <optional>
#include <utility>

template <typename T>
class Result {
private:
    T value;
    std::optional<Err> error;

public:
    // Constructors
    Result(const T& value, std::optional<Err> error = std::nullopt)
        : value(value), error(error) {}

    Result(T&& value, std::optional<Err> error = std::nullopt)
        : value(std::move(value)), error(error) {}

    // Methods
    bool has_error() const { return error.has_value(); }

    T& unwrap() {
        if (has_error()) {
            panic(error.value().get_message().c_str());
        }
        return value;
    }

    T unwrap_or(const T& default_value) const {
        return has_error() ? default_value : value;
    }

    const Err& get_error() const {
        if (!has_error()) {
            panic("Attempted to get error from a successful Result");
        }
        return error.value();
    }
};

std::ostream& operator<<(std::ostream& os, const Err& err);

#endif
