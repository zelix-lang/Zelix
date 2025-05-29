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

#include <fluent/cli/cli.h> // fluent_libc

int main()
{
    cli_app_t app;
    if (!cli_new_app(&app)) // Initialize the CLI application
    {
        // Handle failure
        puts("Error: Failed to initialize the CLI application.");
        return 1; // Exit if app initialization fails
    }

    // Free the CLI application resources
    cli_destroy_app(&app, FALSE);
    return 0;
}