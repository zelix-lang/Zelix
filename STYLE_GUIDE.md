<div align="center">
    <img src="https://assets.zelixlang.dev/logo.png?update=true" height="60" width="60">
    <h1>The Zelix Programming Language</h1>
    Zelix is a modern and blazing-fast programming language.
</div>

---

## 👋 Welcome

Welcome to the official Zelix repository! In this file, you'll find a guide on how to style your code  
to contribute to the Zelix language.

---

## 📝 The Basics

> ⚠️ Code that doesn't follow these guidelines **will not** be accepted.

### **General Rules**

- Use **4 spaces** for indentation — no tabs (check the `.clang-format` file for this).
- No **trailing whitespace** or empty lines at the end of files.
- Avoid **long spaghetti lines** — keep logic clean and concise.
- No **commented-out code**. If code is removed, **explain why** in a comment.
- Remove **unused variables** and **imports**.
- Avoid **unnecessary type casting**.
- Code must be **free of syntax errors**.
- Use **Allman brace style**: every `{` and `}` goes on **its own line**, no exceptions.

---

## 1. Function Invocations

If parameters are short, keep them on one line. If long, split clearly.

```c++
// Good
myFunction(param1, param2, param3);

// Also Good
myFunction(
    myOtherFunction(param1).accessSomething(),
    myOtherFunction(param2).accessSomething(),
    myOtherFunction(param3).accessSomething()
);

// Bad
myFunction(
    param1,
    param2,
    param3
);
```

---

## 2️. Braces (Allman Style)

Always use the **Allman brace convention**:

```c++
// Good
if (condition)
{
    doSomething();
}

// Bad
if (condition) {
    doSomething();
}
```

---

## 3️. Comments

Use **single-line comments** instead of block-style comments.

```c++
// Good
// This is a comment

// Bad
/*
This is a comment
*/
```

---

## 4️⃣ Naming Conventions

- `snake_case` → variables and functions
- `PascalCase` → types and interfaces
- `UPPER_SNAKE_CASE` → constants
- `snake_case` → file and directory names
- lowercase → package names

---

## 5️. Copyright

Every file must begin with the following header:

```c++
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