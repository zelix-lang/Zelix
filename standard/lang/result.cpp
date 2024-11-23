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

#include "err.h"
#include "result.h"
#include "panic.cpp"
#include <optional>

using namespace std;

template <typename T>
Result<T>::Result(T value, std::optional<Err> error) : value(value), error(error) {}

template <typename T>
bool Result<T>::has_error() const {
    return error.has_value();
}

template <typename T>
T* Result<T>::unwrap() {
    if (has_error()) {
        panic(&(error.value()).get_message());
    }

    return value;
}

template <typename T>
T* Result<T>::unwrap_or(T* default_value) {
    if (has_error()) {
        return default_value;
    }

    return value;
}

template <typename T>
Err Result<T>::get_error() const {
    return error.value();
}