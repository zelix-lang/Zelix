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
// Created by rodrigo on 6/4/25.
//

#ifndef FLUENT_MESSAGE_GENERATOR_H
#define FLUENT_MESSAGE_GENERATOR_H

// ============= FLUENT LIB C =============
#include <fluent/string_builder/string_builder.h> // fluent_libc
#include <fluent/std_bool/std_bool.h> // fluent_libc
#include <fluent/ansi/ansi.h> // fluent_libc
#include <fluent/itoa/itoa.h> // fluent_libc
#include <fluent/str_has_prefix/str_has_prefix.h> // fluent_libc

// ============= INCLUDES =============
#include <math.h>

// ============= GLOBAL VARS =============
static size_t cwd_len = 0; // Length of the current working directory

static bool have_same_digits(size_t a, size_t b)
{
    // Fast path: Quick-n-dirty O(log n) for small numbers
    if (a <= 1000)
    {
        // Check if b is greater than 1999
        if (b > 1999 && a == 1000 || b >= 1000 && a < 1000)
        {
            // Drop immediately
            // This is a quick check for numbers that are too far apart
            return FALSE;
        }

        size_t count_a = a == 0 ? 1 : 0; // Edge case for zero
        size_t count_b = b == 0 ? 1 : 0; // Edge case for zero

        // Iterate over the digits of a
        while (a > 0)
        {
            // Increment the counter for a
            count_a++;
            a /= 10; // Remove the last digit
        }

        // Iterate over the digits of b
        while (b > 0)
        {
            // Increment the counter for b
            count_b++;
            b /= 10; // Remove the last digit
        }

        // Check for equality
        return count_a == count_b;
    }

    // Slow path: O(log n) for larger numbers, usually faster for larger numbers
    // Check if the two numbers have the same number of digits
    return a == b || (a == 0 ? 1 : (int)log10(a) + 1) == (b == 0 ? 1 : (int)log10(b) + 1);
}

static void write_line(
    string_builder_t *builder,
    const size_t line_count
)
{
    // Convert the line number to a string
    char *line_number = itoa_convert(line_count);
    write_string_builder(builder, line_number);
    // line_number is now copied, it is safe to free it
    free(line_number); // Free the line number string

    // Write format
    write_string_builder(builder, " | ");
}

static void write_pinpoint(
    string_builder_t *builder,
    const size_t line,
    const size_t column,
    const char *const highlight_color,
    const size_t space_count,
    const size_t char_count,
    const size_t col_start
)
{
    write_string_builder(builder, "     ");
    write_string_builder(builder, highlight_color);
    write_string_builder(builder, "> ");

    // Write the line
    write_line(builder, line);

    // Write spaces before the error
    for (size_t i = 0; i < space_count; i++)
    {
        write_char_string_builder(builder, ' ');
    }

    // Get the real column
    const size_t real_column = column - 1;
    const size_t real_col_start = col_start - 1;

    // Write the caret to highlight the error
    for (size_t i = 0; i < char_count; i++)
    {
        // Write the caret character at the column
        if (i >= real_col_start && i <= real_column)
        {
            write_char_string_builder(builder, '^');
        }
        else
        {
            write_char_string_builder(builder, '-');
        }
    }

    // Write an ANSI reset code
    write_string_builder(builder, ANSI_RESET);
    write_string_builder(builder, "\n");
}

static inline char *build_error_message(
    const char *const code,
    const char *const real_path,
    const char *const highlight_color,
    const size_t line,
    const size_t column,
    const size_t col_start
)
{
    // See if we have to request the CWD
    if (cwd_len == 0)
    {
        cwd_len = strlen(get_cwd());
    }

    // Create a new string builder for the error message
    string_builder_t builder; // Don't initialize it just yet
    size_t line_count = 1; // Initialize line count
    const size_t start_counting_at = line - 1; // Start counting lines from the line before the error
    const size_t end_line = line + 1; // The line where we should stop counting
    const size_t end_counting_at = line + 2; // End counting lines at the line after the error
    bool allowed_to_write = FALSE; // Flag to control writing to the builder
    size_t space_count = 0; // Count the spaces before the first character in the error message
    size_t char_count = 0; // Count the characters in the error message

    // Iterate over the code and file to build the error message
    for (size_t i = 0; code[i] != '\0'; i++)
    {
        // Get the current character
        const char c = code[i];

        // Check if we are allowed to write to the builder
        if (allowed_to_write)
        {
            if (line_count == line)
            {
                // Check if we have a space character
                if (c == ' ')
                {
                    // Check if we have already written a character
                    if (char_count > 0)
                    {
                        // Increment the character count
                        char_count++;
                    }
                    else
                    {
                        // Increment the space count
                        space_count++;
                    }
                }
                else
                {
                    // Increment the character count
                    char_count++;
                }
            }

            // Write the character to the builder
            write_char_string_builder(&builder, c);
        }

        // Check if we have a newline character
        if (c == '\n')
        {
            // Check if we have to write the pinpoint
            if (line_count == line)
            {
                // Write the pinpoint for the error
                write_pinpoint(
                    &builder,
                    line,
                    column,
                    highlight_color,
                    space_count,
                    char_count,
                    col_start
                );
            }

            // Increment the line count
            line_count++;

            // Init the string builder if we have reached the counting range
            if (line_count == start_counting_at)
            {
                // Initialize the string builder
                init_string_builder(&builder, 256, 1.5); // Start with a capacity of 256 characters
                allowed_to_write = TRUE; // Allow writing to the builder
            }

            // Bail out if we have reached the end of the counting range
            else if (end_counting_at == line_count)
            {
                // Write an ANSI reset code
                write_string_builder(&builder, ANSI_RESET);

                break; // Stop counting lines
            }

            // Write the line number if we are allowed to write
            if (allowed_to_write)
            {
                // Write an ANSI reset code
                write_string_builder(&builder, ANSI_RESET);

                // Write the appropriate spacing
                if (line_count == line)
                {
                    write_string_builder(&builder, "     ");
                    write_string_builder(&builder, highlight_color);
                    write_string_builder(&builder, "> ");
                }
                else
                {
                    const char *spaces = "       "; // 7 spaces for the line number

                    // Check if we have to add 1 space less
                    if (
                        line_count == end_line &&
                        line % 10 == 9 &&
                        !have_same_digits(line_count, line)
                    )
                    {
                        spaces = "      "; // 6 spaces for the line number
                    }

                    // Check if we have to add 1 space more
                    if (
                        line_count == start_counting_at &&
                        line_count % 10 == 9 &&
                        !have_same_digits(line_count, line)
                    )
                    {
                        spaces = "        "; // 8 spaces for the line number
                    }

                    write_string_builder(&builder, ANSI_BRIGHT_BLACK);
                    write_string_builder(&builder, spaces);
                }

                // Write the line
                write_line(&builder, line_count);
            }
        }
    }

    // Handle cases where (EOF == end_line)
    if (line_count == end_line && allowed_to_write)
    {
        // Write an ANSI reset code
        write_string_builder(&builder, ANSI_RESET);
        write_string_builder(&builder, "\n"); // Write a newline
    }

    // Check if ended writing at the line (EOF)
    if (line_count == line && allowed_to_write)
    {
        write_string_builder(&builder, ANSI_RESET); // Reset the ANSI color
        write_string_builder(&builder, "\n"); // Write a newline

        // Write the pinpoint for the error
        write_pinpoint(
            &builder,
            line,
            column,
            highlight_color,
            space_count,
            char_count,
            col_start
        );
    }

    // Write the file location
    write_string_builder(&builder, ANSI_BOLD_BRIGHT_PURPLE);
    write_string_builder(&builder, "           => ");

    // Check if the path starts with the current working directory
    if (str_has_prefix(real_path, get_cwd()))
    {
        // Calculate the path without the CWD
        const char *write_path = real_path + cwd_len;

        // Make sure the path doesn't start with the path separator
        if (*write_path == PATH_SEPARATOR)
        {
            write_path++; // Skip the path separator
        }

        // Write the file path without the current working directory
        write_string_builder(&builder, write_path);
    }
    else
    {
        // Write the full file path
        write_string_builder(&builder, real_path);
    }

    // Write the file location line
    write_string_builder(&builder, ":");
    char *line_str = itoa_convert(line);
    write_string_builder(&builder, line_str);
    free(line_str); // Free the line number string

    char *column_str = itoa_convert(column);
    write_string_builder(&builder, ":");
    write_string_builder(&builder, column_str);
    free(column_str); // Free the column number string

    // Write an ANSI reset code
    write_string_builder(&builder, ANSI_RESET);
    write_string_builder(&builder, "\n");

    // Return the built error message
    return collect_string_builder_no_copy(&builder);
}

#endif //FLUENT_MESSAGE_GENERATOR_H
