<div align="center">
    <img src="assets/logo.png" height="60" width="60">
    <h1>Fluent Language</h1>
    Fluent is a modern and blazing-fast programming language.
</div>

---

## 👋 Welcome

Welcome to the official Fluent repository! Here, you will find the source code of the Fluent language. Have fun!

> **NOTE**
> This project is still in development and may not be stable.
> **Fluent is still not a language for production use.**

---

## 📦 Features

- Blazing-fast execution times
- Simple syntax
- Flexible standard library
- Easily package your app into an executable
- **[Free - as in Freedom.](LICENSE)**
- And more!

---

## 🎆 Installation

To install Fluent, you can download the official Fluent installer from the **Releases** pages.
Once the installer shows that the installation was successful, you may need to restart your terminal for changes to take effect.

- **On Windows**: Just close and re-open the terminal (CMD, PowerShell, etc.)
- **On Linux/macOS**: Execute `source ~/.bashrc` or `source ~/.zshrc` depending on your shell

---

## 🚀 Getting Started

To create a new Fluent project, you can use the `fluent init` command. This will create a new Fluent project in the current directory.

> **NOTE:** Make sure the current directory is empty before running the `fluent init` command, otherwise, add a name after the command to create a new directory with the project name, e.g. `fluent init my_project`.

```shell
fluent init
```

You will be prompted to fill relevant information about your project like project name, author, etc.
Once you've filled in the information, the project will be created and you can start coding!

**Example structure:**

```
my_project/
    src/
        main.fluent
    Fluent.yml
    .gitignore
```

To run your code use:

```
fluent run
```

---

## 📚 Documentation

The official Fluent documentation can be found on the [official website](https://fluent-lang.github.io/docs).

---

## 📦 Contributions

Contributions are welcome! If you'd like to contribute to the Fluent language, please read the [CONTRIBUTING.md](CONTRIBUTING.md) file for more information.

---

## 🎲 Building from Source

To build Fluent from source into an executable, you need to execute either one of the **build scripts**:

- **On Windows**: `build.bat`
- **On Linux/macOS**: `build.sh`

After the build process is complete, you will find the executable in the `bin/` directory.

---

## 🔒 Security

Please refer to [SECURITY.md](SECURITY.md) for more information on how to report security vulnerabilities.

---

## 📝 License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for more information.

```
Copyright (C) 2025 Rodrigo R. & All Fluent Contributors
This program comes with ABSOLUTELY NO WARRANTY; for details type `fluent license`.
This is free software, and you are welcome to redistribute it under certain conditions;
type `fluent license --full` for details.
```