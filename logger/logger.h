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

#ifndef FLUENT_LOGGER_H
#define FLUENT_LOGGER_H
#include <stdio.h>
#include <fluent/ansi/ansi.h>

// Global info log prefix
static char *info_prefix = "[INFO] ";
// Global error log prefix
static char *error_log_prefix = "[ERROR] ";
// Global warning log prefix
static char *warn_log_prefix = "[WARN] ";
// Global help log prefix
static char *help_log_prefix = "[HELP] ";

static inline void log_info(const char *message)
{
    printf(
        "%s%s%s%s\n",
        ANSI_BOLD_BRIGHT_BLUE,
        info_prefix,
        message,
        ANSI_RESET
    );
}

static inline void log_error(const char *message)
{
    fprintf(
        stderr,
        "%s%s%s%s\n",
        ANSI_BOLD_BRIGHT_RED,
        error_log_prefix,
        message,
        ANSI_RESET
    );
}

static inline void log_warning(const char *message)
{
    fprintf(
        stderr,
        "%s%s%s%s\n",
        ANSI_BOLD_BRIGHT_YELLOW,
        warn_log_prefix,
        message,
        ANSI_RESET
    );
}

static inline void log_help(const char *message)
{
    printf(
        "%s%s%s%s\n",
        ANSI_BOLD_BRIGHT_GREEN,
        help_log_prefix,
        message,
        ANSI_RESET
    );
}

#endif //FLUENT_LOGGER_H
