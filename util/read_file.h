/*
        ==== The Zelix Programming Language ====
---------------------------------------------------------
  - This file is part of the Zelix Programming Language
    codebase. Fluent is a fast, statically-typed and
    memory-safe programming language that aims to
    match native speeds while staying highly performant.
---------------------------------------------------------
  - Fluent is categorized as free software; you can
    redistribute it and/or modify it under the terms of
    the GNU General Public License as published by the
    Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.
---------------------------------------------------------
  - Fluent is distributed in the hope that it will
    be useful, but WITHOUT ANY WARRANTY; without even
    the implied warranty of MERCHANTABILITY or FITNESS
    FOR A PARTICULAR PURPOSE. See the GNU General Public
    License for more details.
---------------------------------------------------------
  - You should have received a copy of the GNU General
    Public License along with Fluent. If not, see
    <https://www.gnu.org/licenses/>.
*/

//
// Created by rodrigo on 8/1/25.
//

#pragma once

#include <cstdio>
#include <cstdlib>
#include "zelix/container/owned_string.h"

namespace zelix::util
{
    /**
     * Reads the contents of a file into a container::owned_string.
     * @param path The path to the file to read.
     * @return A container::owned_string containing the file contents.
     */
    inline container::string read_file(const char* path)
    {
        char line[256]; // buffer for each line

        // Open the file in read mode
        FILE *file = fopen(path, "r");

        if (file == nullptr)
        {
            throw except::exception("Failed to open file");
        }

        container::string content; // String to hold the file contents
        // Read each line until EOF
        while (fgets(line, sizeof(line), file))
        {
            content.push(line); // Append the line to the content string
        }

        // Close the file
        fclose(file);

        return content;
    }
}