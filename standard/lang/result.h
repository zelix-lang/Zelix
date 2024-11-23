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

#include "err.h"
#include <optional>

using namespace std;

template <typename T>
class Result {
    private:
        T value;
        optional<Err> error;
    public:
        // Constructors
        Result(T value, optional<Err> error);

        // Methods
        bool has_error() const;
        T* unwrap();
        T* unwrap_or(T* default_value);
        Err get_error() const;
};

#endif
