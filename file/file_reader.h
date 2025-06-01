/*
    The Fluent Programming Language
    -----------------------------------------------------
    This code is released under the GNU GPL v3 license.
    For more information, please visit:
    https://www.gnu.org/licenses/gpl-3.0.html
    -----------------------------------------------------
    Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
    This program comes with ABSOLUTELY NO WARRANTY.
    For details type `fluent l`. This is free software,
    and you are welcome to redistribute it under certain
    conditions; type `fluent l -f` for details.
*/

//
// Created by rodrigo on 5/31/25.
//

#ifndef FLUENT_FILE_READER_H
#define FLUENT_FILE_READER_H
#include <stdio.h>
#include <stdlib.h>
#include <fluent/path/path.h>

/**
 * Reads the entire contents of a file into a dynamically allocated buffer.
 *
 * @param path The path to the file to be read.
 * @return A pointer to a null-terminated string containing the file contents,
 *         or NULL if the file could not be opened or memory allocation failed.
 *         The caller is responsible for freeing the returned buffer.
 */
static inline char *read_file(const char *const path)
{
    // Normalize the path
    char *normalized_path = get_real_path(path);
    if (!normalized_path)
    {
        return NULL; // Failed to normalize the path
    }

    FILE *file = fopen(normalized_path, "r");
    if (!file)
    {
        return NULL; // Failed to open the file
    }

    fseek(file, 0, SEEK_END);
    const long length = ftell(file);
    fseek(file, 0, SEEK_SET);

    char *buffer = (char *)malloc(length + 1);
    if (!buffer)
    {
        // Free the normalized path if memory allocation fails
        free(normalized_path);

        fclose(file);
        return NULL; // Failed to allocate memory
    }

    fread(buffer, 1, length, file);
    buffer[length] = '\0'; // Null-terminate the string

    // Free the normalized path as it's no longer needed
    free(normalized_path);
    fclose(file);
    return buffer;
}

#endif //FLUENT_FILE_READER_H
