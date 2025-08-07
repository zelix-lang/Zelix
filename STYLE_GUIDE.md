<div align="center">
    <img src="https://assets.zelixlang.dev/logo.png" height="60" width="60">
    <h1>The Zelix Programmign Language</h1>
    Zelix a modern and blazing-fast programming language.
</div>

---

## 👋 Welcome

Welcome to the official Zelix repository! In this file, you'll find a guide on how to style your code
to contribute to the Zelix language.

---

## 📝 The basics

Please note that any code that does not follow the guidelines in this document will not be accepted.

**Guidelines:**

- 4 spaces for indentation, no tabs (You should have this configured as per the .clang-format file)
- No trailing whitespace or empty lines at the end of files
- No outrageously long lines of code, try to keep logic as concise as possible
  - For example:
    - ```c++
      // Bad
      if (err != nullptr)
      { 
        return err
      }
        ```
      
    - ```c++
      // Bad
	  if (err != nullptr || my_file_is_okay > 0 && ((myOtherFile == "okay" || myOtherFile == "not okay") || fetchSomeResource() == "13.5")) {
        return err
	  }
        ```
      
- No commented-out code, if you remove code, you must provide a comment explaining why it is no longer needed
- No unused variables or imports
- No unnecessary type casting
- No syntax errors
- Every curly brace must have its own line, no exceptions.

**1. Function invocations**

Unless the parameters are too long, function invocations should be on the same line.

Example:

```go
// Good
myFunction(param1, param2, param3)

// Also Good
myFunction(
    myOtherFunction(param1).accessSomething(),
	myOtherFunction(param2).accessSomething(),
    myOtherFunction(param3).accessSomething(),
)

// Bad
myFunction(
    param1,
    param2,
    param3,
)
```

**3. Comments**

Avoid block comments, use single-line comments instead.

Example:

```go
// Good
// This is a comment

// Bad
/*
This is a comment
*/
```

**4. Naming conventions**

- Use `snake_case` for variables and functions
- Use `PascalCase` for types and interfaces
- Use `UPPER_SNAKE_CASE` for constants
- Use `snake_case` for file and directory names
- Use lowercase for package names

**6. Copyright**

All files must contain the following header at the top of the file:

```go
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
```