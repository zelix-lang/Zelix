<div align="center">
    <img src="assets/logo.png" height="60" width="60">
    <h1>Fluent Language</h1>
    Fluent is a modern and blazing-fast programming language.
</div>

---

## 👋 Welcome

Welcome to the official Fluent repository! In this file, you'll find a guide on how to style your code
to contribute to the Fluent language.

---

## 📝 The basics

Please note that any code that does not follow the guidelines in this document will not be accepted.

**Guidelines:**

- 4 spaces for indentation, no tabs (You should have this configured as per the .editorconfig file)
- No trailing whitespace or empty lines at the end of files
- No outrageously long lines of code, try to keep logic as concise as possible
  - For example:
    - ```go
      // Bad
      if err != nil { 
        return err
      }
        ```
      
    - ```go
      // Bad
	  if err != nil || my_file_is_okay > 0 && ((myOtherFile == "okay" || myOtherFile == "not okay") || fetchSomeResource() == "13.5") {
        return err
	  }
        ```
      
- No commented-out code, if you remove code, you must provide a comment explaining why it is no longer needed
- No unused variables or imports
- No unnecessary type casting
- No syntax errors

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

- Use `camelCase` for variables and functions
- Use `PascalCase` for types and interfaces
- Use `UPPER_SNAKE_CASE` for constants
- Use `camelCase` for file and directory names
- Use lowercase for package names

**5. Error handling**

Always handle errors, do not ignore them.

Example:

```go
// Good
if err != nil {
    return err
}

// Bad
if err != nil {
    // Do nothing
}
```

Unless there is underlying logic that requires it, do not use `panic()`.
Instead use the `logger` package to log information, warns and errors.

**6. Copyright**

All files must contain the following header at the top of the file:

```go
/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 (Your Name) & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/
```