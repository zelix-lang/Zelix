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
#include <fluent/ansi/ansi.h> // fluent_libc

// ============= INCLUDES =============
#include "token/token_map.h"
#include "command/check.h"
#include "state/state.h"

// ============= MACROS =============
#ifndef PROGRAM_NAME
#   define PROGRAM_NAME "fluent"
#endif

#ifndef PROGRAM_DESCRIPTION
#   define PROGRAM_DESCRIPTION "The Fluent Programming Language"
#endif

int main(const int argc, const char **const argv)
{
    atexit(timer_done); // Ensure all timers are stopped on exit

    cli_app_t app;
    if (!cli_new_app(&app)) // Initialize the CLI application
    {
        // Handle failure
        puts("Error: Failed to initialize the CLI application.");
        return 1; // Exit if app initialization fails
    }

    // Fill the CLI application with flags
    cli_value_t help_flag = cli_new_value("Displays this help message", CLI_TYPE_STATIC, "h", FALSE);
    cli_value_t c_flags = cli_new_value("Specifies compiler flags", CLI_TYPE_STRING, "cf", FALSE);
    cli_value_t opt_level = cli_new_value("Sets the optimization level", CLI_TYPE_INTEGER, "O", FALSE);
    cli_insert_flag(&app, "help", &help_flag);
    cli_insert_flag(&app, "c_flags", &c_flags);
    cli_insert_flag(&app, "optimization", &opt_level);

    // Fill the CLI application with commands
    cli_value_t compile_cmd = cli_new_value("Builds Fluent source code and outputs an executable", CLI_TYPE_STRING, "b", FALSE);
    cli_value_t run_cmd = cli_new_value("Runs the Fluent source code", CLI_TYPE_STRING, "r", FALSE);
    cli_value_t check_cmd = cli_new_value("Performs static analyzer checks", CLI_TYPE_STRING, "c", FALSE);
    cli_insert_command(&app, "build", &compile_cmd);
    cli_insert_command(&app, "run", &run_cmd);
    cli_insert_command(&app, "check", &check_cmd);

    // Parse the command line arguments
    argv_t args = parse_argv(argc, argv, &app);

    // Handle failure
    if (args.success == FALSE || hashmap_cli_i_get(args.statics, "help") != NULL)
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

    // Print the header message
    printf("%sThe Fluent Programming Language%s\n", ANSI_BOLD_BRIGHT_BLUE, ANSI_RESET);
    printf("%sA blazingly fast programming language%s\n\n", ANSI_BRIGHT_BLACK, ANSI_RESET);

    // Initialize the token map for further processing
    get_token_map();
    get_punctuation_map();
    get_chainable_map();

    // Handle the commands and flags
    if (args.cmd_ptr == &compile_cmd)
    {
        puts("Build command executed.");
    } else if (args.cmd_ptr == &run_cmd)
    {
        puts("Run command executed.");
    }
    else if (args.cmd_ptr == &check_cmd)
    {
        check_command(args.command.value);
    }

    // Free the CLI application resources
    cli_destroy_app(&app, FALSE);
    destroy_argv(&args);
    return 0;
}