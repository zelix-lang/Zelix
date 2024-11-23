<div align="center">
    <img src="assets/logo.png" height="60" width="60">
    <h1>Surf Language</h1>
    Surf is a modern and blazing-fast programming language.
</div>

---

## 👋 Welcome

Welcome to the official Surf repository! Here, you will find the source code of the Surf language. Have fun!

---

## 📦 Features

- Blazing-fast execution times
- Simple syntax
- Flexible standard libraries
- Easily package your app into an executable
- And more!

---

## 🎆 Installation

To install Surf, **you'll need [CLang++](https://clang.llvm.org/) and [LLVM](https://llvm.org/) installed**. Installation steps may vary depending on your OS.
Once installed, you can download the official Surf installer from the **Releases** pages.

> **NOTE:**
> If you have CLang installed, you can skip the CLang++ installation step, as CLang++ is included in the CLang package.
> To further verify that CLang++ is installed, you can run `clang++ --version` in your terminal.

Once the installer shows that the installation was successful, you may need to restart your terminal for changes to take effect.

- **On Windows**: Just close and re-open the terminal (CMD, PowerShell, etc.)
- **On Linux/macOS**: Execute `source ~/.bashrc` or `source ~/.zshrc` depending on your shell

---

## 🚀 Getting Started

To create a new Surf project, you can use the `surf init` command. This will create a new Surf project in the current directory.

> **NOTE:** Make sure the current directory is empty before running the `surf init` command, otherwise, add a name after the command to create a new directory with the project name, e.g. `surf init my_project`.

```shell
surf init
```

You will be prompted to fill relevant information about your project like project name, author, etc.
Once you've filled in the information, the project will be created and you can start coding!

**Example structure:**

```
my_project/
    src/
        main.surf
    Surf.yml
    .gitignore
```

---

## 📚 Documentation

The official Surf documentation can be found on the [official website](https://rodri-r-z.github.io/surf/docs).

---

## 📦 Contributions

Contributions are welcome! If you'd like to contribute to the Surf language, please read the [CONTRIBUTING.md](CONTRIBUTING.md) file for more information.

---

## 📝 License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for more information.

```
Copyright (C) 2024 Rodrigo R. & All Surf Contributors
This program comes with ABSOLUTELY NO WARRANTY; for details type `surf license`.
This is free software, and you are welcome to redistribute it under certain conditions;
type `surf license --full` for details.
```