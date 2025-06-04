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

#ifndef FLUENT_STATE_H
#define FLUENT_STATE_H

// ============= FLUENT LIB C =============
#include <fluent/clock/clock.h> // fluent_libc
#include <fluent/std_bool/std_bool.h> // fluent_libc

/**
 * \brief Represents a timer state in the Fluent system.
 *
 * This structure holds information about a timer, including its start time
 * and an associated message.
 */
typedef struct
{
    hr_clock_t start_time;  /**< Start time of the timer */
    const char *message;    /**< The message of the timer */
} state_timer_t;

// ============= MACROS =============
#ifndef FLUENT_TIMER_SUCCESS_CHAR
#   define FLUENT_TIMER_SUCCESS_CHAR '✔' // Character to indicate successful timer completion
#endif

#ifndef FLUENT_TIMER_FAILURE_CHAR
#   define FLUENT_TIMER_FAILURE_CHAR '✘' // Character to indicate failed timer completion
#endif

#ifndef FLUENT_TIMER_CLOCK_CHAR
#   define FLUENT_TIMER_CLOCK_CHAR '⏱' // Character to indicate a timer is running
#endif

// ============= GLOBAL VARIABLES =============
state_timer_t current;
bool timer_running = FALSE; // Flag to indicate if a timer is currently running

static inline void timer_done()
{
    // Skip if no timer is running
    if (!timer_running)
    {
        return; // No timer to complete
    }

    // Set the timer running flag to FALSE
    timer_running = FALSE;

    // Calculate the elapsed time
    const long long elapsed = hr_clock_distance_from_now(&current.start_time, CLOCK_MICROSECONDS);

    // Print the timer message with elapsed time
    printf(
        "%s(%s%s%c%s%s) %s - %lldμs%s\n",
        ANSI_BRIGHT_BLACK,
        ANSI_RESET,
        ANSI_BOLD_BRIGHT_GREEN,
        FLUENT_TIMER_SUCCESS_CHAR,
        ANSI_RESET,
        ANSI_BRIGHT_BLACK,
        current.message,
        elapsed,
        ANSI_RESET
    );
}

static inline void new_timer(const char *message) {
    // Check if we have a timer running
    if (timer_running)
    {
        timer_done(); // If a timer is already running, complete it
    }

    // Set the timer running flag to TRUE
    timer_running = TRUE;

    // Print the start message
    printf(
        "%s(%s%s%c%s%s) %s%s\r",
        ANSI_BRIGHT_BLACK,
        ANSI_RESET,
        ANSI_BOLD_BRIGHT_GREEN,
        FLUENT_TIMER_CLOCK_CHAR,
        ANSI_RESET,
        ANSI_BRIGHT_BLACK,
        current.message,
        ANSI_RESET
    );

    // Create a new timer state
    state_timer_t timer;
    hr_clock_tick(&timer.start_time); // Initialize the start time of the timer
    timer.message = message;

    current = timer; // Set the current timer state to the new timer
}

#endif //FLUENT_STATE_H
