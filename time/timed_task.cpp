/*
        ==== The Fluent Programming Language ====
---------------------------------------------------------
  - This file is part of the Fluent Programming Language
    codebase. Fluent is a fast, statically-typed and
    memory-safe programming language that aims to
    match native speeds while staying highly performant.
---------------------------------------------------------
  - Fluent is categorized as free software; you can
    redistribute it and/or modify it under the terms of
    the GNU General Public License as published by the
    Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.
---------------------------------------------------------
  - Fluent is distributed in the hope that it will
    be useful, but WITHOUT ANY WARRANTY; without even
    the implied warranty of MERCHANTABILITY or FITNESS
    FOR A PARTICULAR PURPOSE. See the GNU General Public
    License for more details.
---------------------------------------------------------
  - You should have received a copy of the GNU General
    Public License along with Fluent. If not, see
    <https://www.gnu.org/licenses/>.
*/

#include "timed_task.h"

#include <chrono>

struct timed_task
{
    const char *name = nullptr;
    std::chrono::system_clock::time_point start_time = std::chrono::system_clock::now();
    size_t took = 0;
    int steps = 0;
    int max_steps = 0;
};

// Save the current task
timed_task task;

void fluent::time::complete(const bool recompute_time)
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
    printf("[%d/%d] %s - %lldµs\n", task.steps, task.max_steps, task.name, task.took);
    task.name = nullptr; // Clear the task name
}

void fluent::time::advance()
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

    printf("[%d/%d] %s\r", task.steps, task.max_steps, task.name);
}

void fluent::time::post(const char *name, const int max_steps)
{
    complete();
    task.name = name; // Set the task name
    task.max_steps = max_steps; // Set the maximum steps
    task.start_time = std::chrono::system_clock::now(); // Reset the start time
    task.took = 0; // Reset the time taken
    task.steps = 0; // Reset the steps
    printf("[%d/%d] %s\r", task.steps, task.max_steps, name);
}
