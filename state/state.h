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

/**
 * \enum state_event_t
 * \brief Enumerates the possible states/events in the Fluent system.
 *
 * This enum represents the various processing stages that the Fluent system
 * can be in, such as lexing, parsing, processing, analyzing, and building.
 */
typedef enum
{
    STATE_LEXING = 0, /**< The state when the lexer is processing input */
    STATE_PARSING,    /**< The state when the parser is processing input */
    STATE_PROCESSING, /**< The state when the program's internal code is processing */
    STATE_ANALYZING,  /**< The state when the static analyzer is checking code */
    STATE_BUILDING,   /**< The state when the IR builder is building the code */
} state_event_t;

// ============= MACROS =============
#ifndef FLUENT_TIMER_SUCCESS_STR
#   define FLUENT_TIMER_SUCCESS_STR "DONE"
#endif

#ifndef FLUENT_TIMER_FAILURE_STR
#   define FLUENT_TIMER_FAILURE_STR "FAILED"
#endif

#ifndef FLUENT_TIMER_CLOCK_STR
#   define FLUENT_TIMER_CLOCK_STR "..."
#endif

// ============= GLOBAL VARIABLES =============
state_timer_t current;
bool timer_running = FALSE; // Flag to indicate if a timer is currently running

static void timer_finalized(const char *status, const char *color)
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
        "%s(%s%s%s%s%s) %s - %lldμs%s\n",
        ANSI_BRIGHT_BLACK,
        ANSI_RESET,
        color,
        status,
        ANSI_RESET,
        ANSI_BRIGHT_BLACK,
        current.message,
        elapsed,
        ANSI_RESET
    );
}

static inline void timer_failed()
{
    timer_finalized(FLUENT_TIMER_FAILURE_STR, ANSI_BOLD_BRIGHT_RED);
}

static inline void timer_done()
{
    timer_finalized(FLUENT_TIMER_SUCCESS_STR, ANSI_BOLD_BRIGHT_GREEN);
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
        "%s(%s%s%s%s%s) %s%s\r",
        ANSI_BRIGHT_BLACK,
        ANSI_RESET,
        ANSI_BOLD_BRIGHT_GREEN,
        FLUENT_TIMER_CLOCK_STR,
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
