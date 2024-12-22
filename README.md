<div align="center">
    <img src="assets/logo.png" height="60" width="60">
    <h1>Zyro Language</h1>
    Zyro is a modern and blazing-fast programming language.
</div>

---

## 👋 Welcome

Welcome to the official Zyro repository! Here, you will find the source code of the Zyro language. Have fun!

> **NOTE**
> This project is still in development and may not be stable.
> **Zyro is still not a language for production use.**

---

## 📦 Features

- Blazing-fast execution times
- Simple syntax
- Flexible standard libraries
- Easily package your app into an executable
- And more!

---

## 🎆 Installation

To install Zyro, you can download the official Zyro installer from the **Releases** pages.
Once the installer shows that the installation was successful, you may need to restart your terminal for changes to take effect.

- **On Windows**: Just close and re-open the terminal (CMD, PowerShell, etc.)
- **On Linux/macOS**: Execute `source ~/.bashrc` or `source ~/.zshrc` depending on your shell

---

## 🚀 Getting Started

To create a new Zyro project, you can use the `zyro init` command. This will create a new Zyro project in the current directory.

> **NOTE:** Make sure the current directory is empty before running the `zyro init` command, otherwise, add a name after the command to create a new directory with the project name, e.g. `zyro init my_project`.

```shell
zyro init
```

You will be prompted to fill relevant information about your project like project name, author, etc.
Once you've filled in the information, the project will be created and you can start coding!

**Example structure:**

```
my_project/
    src/
        main.zyro
    Zyro.yml
    .gitignore
```

---

## 📚 Documentation

The official Zyro documentation can be found on the [official website](https://rodri-r-z.github.io/Zyro/docs).

---

## 📦 Contributions

Contributions are welcome! If you'd like to contribute to the Zyro language, please read the [CONTRIBUTING.md](CONTRIBUTING.md) file for more information.

---

## 🎲 Building from Source

To build Zyro from source into an executable, you need to execute either one of the **build scripts**:

- **On Windows**: `build.bat`
- **On Linux/macOS**: `build.sh`

After the build process is complete, you will find the executable in the `bin/` directory.

---

## 📝 License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for more information.

```
Copyright (C) 2024 Rodrigo R. & All Zyro Contributors
This program comes with ABSOLUTELY NO WARRANTY; for details type `zyro license`.
This is free software, and you are welcome to redistribute it under certain conditions;
type `zyro license --full` for details.
```