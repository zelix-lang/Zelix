/*
        ==== The Zelix Programming Language ====
---------------------------------------------------------
  - This file is part of the Zelix Programming Language
    codebase. Zelix is a fast, statically-typed and
    memory-safe programming language that aims to
    match native speeds while staying highly performant.
---------------------------------------------------------
  - Zelix is categorized as free software; you can
    redistribute it and/or modify it under the terms of
    the GNU General Public License as published by the
    Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.
---------------------------------------------------------
  - Zelix is distributed in the hope that it will
    be useful, but WITHOUT ANY WARRANTY; without even
    the implied warranty of MERCHANTABILITY or FITNESS
    FOR A PARTICULAR PURPOSE. See the GNU General Public
    License for more details.
---------------------------------------------------------
  - You should have received a copy of the GNU General
    Public License along with Zelix. If not, see
    <https://www.gnu.org/licenses/>.
*/

//
// Created by rodrigo on 8/1/25.
//

#pragma once

#include <cstdio>
#include <cstdlib>
#include "zelix/container/owned_string.h"
#include "zelix/except/exception.h"

namespace zelix::util
{
    /**
     * Reads the contents of a file into a container::owned_string.
     * @param path The path to the file to read.
     * Reads the contents of a file into a container::string.
     * @param path The path to the file to read.
     * @return A container::string containing the file contents.
     */
    inline container::string read_file(const char* path)
    {
        // Open the file in read mode
        FILE *file = fopen(path, "rb");

        if (file == nullptr)
        {
            throw except::exception("Failed to open file");
        }

        // Seek to the end to determine file size
        if (fseek(file, 0, SEEK_END) != 0)
        {
            fclose(file);
            throw except::exception("Failed to seek file");
        }

        const long file_size = ftell(file);
        if (file_size < 0)
        {
            fclose(file);
            throw except::exception("Failed to determine file size");
        }

        rewind(file);

        // Allocate buffer and read file
        auto* buffer = static_cast<char*>(malloc(file_size + 1));
        if (!buffer)
        {
            fclose(file);
            throw except::exception("Failed to allocate buffer");
        }

        const size_t read_size = std::fread(buffer, 1, file_size, file);
        buffer[read_size] = '\0';

        // Close the file
        fclose(file);

        // Construct the string and free buffer
        container::string content;
        content.push(buffer);
        free(buffer);
        return content;
    }
}