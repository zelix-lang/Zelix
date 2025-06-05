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

// ============= INCLUDES =============
#include <math.h>

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

static inline char *build_error_message(
    const char *const code,
    const char *const file,
    const char *const highlight_color,
    const size_t line,
    const size_t column
)
{
    // Create a new string builder for the error message
    string_builder_t builder; // Don't initialize it just yet
    size_t line_count = 1; // Initialize line count
    const size_t start_counting_at = line - 1; // Start counting lines from the line before the error
    const size_t end_line = line + 1; // The line where we should stop counting
    const size_t end_counting_at = line + 2; // End counting lines at the line after the error
    bool allowed_to_write = FALSE; // Flag to control writing to the builder

    // Iterate over the code and file to build the error message
    for (size_t i = 0; code[i] != '\0'; i++)
    {
        // Get the current character
        const char c = code[i];

        // Check if we are allowed to write to the builder
        if (allowed_to_write)
        {
            // Write the character to the builder
            write_char_string_builder(&builder, c);
        }

        // Check if we have a newline character
        if (c == '\n')
        {
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
            if (end_counting_at == line_count)
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

                // Convert the line number to a string
                char *line_number = itoa(line_count);
                write_string_builder(&builder, line_number);
                // line_number is now copied, it is safe to free it
                free(line_number); // Free the line number string

                // Write format
                write_string_builder(&builder, " | ");
            }
        }
    }

    // Return the built error message
    return collect_string_builder_no_copy(&builder);
}

#endif //FLUENT_MESSAGE_GENERATOR_H
