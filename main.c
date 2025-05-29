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

// ============= FLUENT LIB C =============
#include <fluent/cli/cli.h> // fluent_libc
#include <fluent/cli/help/generator.h> // fluent_libc

// ============= MACROS =============
#ifndef PROGRAM_NAME
#   define PROGRAM_NAME "fluent"
#endif

#ifndef PROGRAM_DESCRIPTION
#   define PROGRAM_DESCRIPTION "The Fluent Programming Language"
#endif

int main(const int argc, const char **const argv)
{
    cli_app_t app;
    if (!cli_new_app(&app)) // Initialize the CLI application
    {
        // Handle failure
        puts("Error: Failed to initialize the CLI application.");
        return 1; // Exit if app initialization fails
    }

    // Fill the CLI application with flags
    cli_value_t help_flag = cli_new_value("Displays this help message", CLI_TYPE_STATIC, "h", FALSE);
    cli_insert_flag(&app, "help", &help_flag);

    // Parse the command line arguments
    argv_t args = parse_argv(argc, argv, &app);

    // Handle failure
    if (args.success == FALSE || hashmap_get(args.statics, "help") != NULL)
    {
        // Generate a help message
        char *help_message = cli_generate_help(&app, PROGRAM_NAME, PROGRAM_DESCRIPTION, 18);

        // Print the help message
        printf("%s", help_message);

        // Free the help message memory
        free(help_message);

        // Free the CLI application resources
        cli_destroy_app(&app, FALSE);
        destroy_argv(&args);

        // Return failure
        return 1;
    }

    // Free the CLI application resources
    cli_destroy_app(&app, FALSE);
    destroy_argv(&args);
    return 0;
}