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

#include "timed_task.h"
#include <fluent/ansi/ansi.h>
#include <chrono>

struct timed_task
{
    const char *name = nullptr;
    int name_len = 0;
    std::chrono::system_clock::time_point start_time = std::chrono::system_clock::now();
    size_t took = 0;
    int steps = 0;
    int max_steps = 0;
    size_t nested = 0;
};

// Save the current task
timed_task task;

template <bool Failed, bool Complete>
void print_task(const size_t nested, const char *reason)
{
    if (nested > 0)
    {
        for (size_t i = 0; i < nested; i++)
        {
            printf("  ");
        }

        if constexpr (Failed)
        {
            printf(ANSI_BRIGHT_RED "└─" ANSI_RESET);
        }
        else if constexpr (Complete)
        {
            printf(ANSI_BRIGHT_GREEN "└─" ANSI_RESET);
        }
        else
        {
            printf(ANSI_BRIGHT_BLACK "└─" ANSI_RESET);
        }
    }

    if constexpr (Failed)
    {
        printf(
            ANSI_BRIGHT_BLACK
            "\033[2m["
            ANSI_RESET
            ANSI_BRIGHT_RED
            "\033[2m%d"
            ANSI_RESET
            ANSI_BRIGHT_BLACK
            "\033[2m/%d]"
            ANSI_RESET
            ANSI_BRIGHT_RED
            " %.*s [x]"
            ANSI_RESET
            "\n    \033[38;5;214m\033[2m(!) what "
            ANSI_RESET
            ANSI_BRIGHT_BLACK
            "~ "
            ANSI_RESET
            "\033[38;5;214m%s\n",
            task.steps,
            task.max_steps,
            task.name_len,
            task.name,
            reason
        );
    }

    else if constexpr (Complete)
    {
        printf(
            ANSI_BRIGHT_GREEN
            "\033[2m[%d/%d]\033[22m"
            ANSI_RESET
            ANSI_BRIGHT_GREEN
            " %.*s"
            ANSI_RESET,
            task.steps,
            task.max_steps,
            task.name_len,
            task.name
        );
    }

    else
    {
        printf(
            ANSI_BRIGHT_BLACK
            "[%d/"
            ANSI_RESET
            ANSI_BRIGHT_BLUE
            "%d"
            ANSI_RESET
            ANSI_BRIGHT_BLACK
            "]"
            ANSI_RESET
            ANSI_BRIGHT_BLUE
            " %.*s"
            ANSI_RESET
            "\r",
            task.steps,
            task.max_steps,
            task.name_len,
            task.name
        );
    }

    fflush(stdout);
}

void zelix::time::complete(const bool recompute_time)
{
    if (task.name == nullptr)
    {
        return; // No task to complete
    }

    // Recompute the time taken if requested
    if (recompute_time)
    {
        // Add the time taken since the last start
        const auto now = std::chrono::system_clock::now();
        task.took += std::chrono::duration_cast<std::chrono::microseconds>(now - task.start_time).count();
    }

    task.steps = task.max_steps;
    print_task<false, true>(task.nested, nullptr); // Print the completed task

    if (task.took < 1000)
    {
        printf(
            ANSI_BRIGHT_BLACK
            " ~ %zuµs"
            ANSI_RESET
            "\n",
            task.took
        );
    }
    else if (task.took < 1000000)
    {
        printf(
            ANSI_BRIGHT_BLACK
            " ~ %.2fms"
            ANSI_RESET
            "\n",
            static_cast<double>(task.took) / 1000.0
        );
    }
    else
    {
        printf(
            ANSI_BRIGHT_BLACK
            " ~ %.2fs"
            ANSI_RESET
            "\n",
            static_cast<double>(task.took) / 1000000.0
        );
    }

    task.name = nullptr; // Clear the task name
}

void zelix::time::advance()
{
    if (task.name == nullptr)
    {
        return; // No task to advance
    }

    // Get the current time
    const auto now = std::chrono::system_clock::now();

    // Update the task with the time taken since the last start
    // (to explicitly exclude the time taken for print syscalls)
    task.took += std::chrono::duration_cast<std::chrono::microseconds>(now - task.start_time).count();
    task.steps++;
    task.start_time = now;

    // Complete the task if the maximum steps have been reached
    if (task.steps >= task.max_steps)
    {
        complete(false);
        return;
    }

    print_task<false, false>(task.nested, nullptr);
}

void zelix::time::post(
    const char *name,
    const int len,
    const int max_steps,
    const size_t nested
)
{
    complete();
    task.name = name; // Set the task name
    task.name_len = len; // Set the task name length
    task.max_steps = max_steps; // Set the maximum steps
    task.start_time = std::chrono::system_clock::now(); // Reset the start time
    task.took = 0; // Reset the time taken
    task.steps = 0; // Reset the steps
    task.nested = nested; // Set the nesting level

    print_task<false, false>(nested, nullptr);
}

void zelix::time::fail(const char *reason)
{
    if (task.name == nullptr)
    {
        return; // No task to advance
    }

    print_task<true, true>(task.nested, reason);
}