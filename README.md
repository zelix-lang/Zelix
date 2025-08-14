<div align="center">
    <img src="https://assets.zelixlang.dev/logo.png?update=true" height="60" width="60">
    <h1>The Zelix Programming Language</h1>
    Zelix is a modern and blazing-fast programming language.
    <br>
    <br>

[![GitHub](https://img.shields.io/github/license/zelix-lang/Zelix?style=flat-square)](LICENSE)
![GitHub repo size](https://img.shields.io/github/repo-size/zelix-lang/Zelix)
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues/zelix-lang/Zelix)
![GitHub Repo stars](https://img.shields.io/github/stars/zelix-lang/Zelix?style=flat)

[Website][Website] | [Documentation][Documentation] | [Contributing][Contributing]
</div>

---

[Website]: https://zelixlang.dev
[Documentation]: https://docs.zelixlang.dev
[Contributing]: CONTRIBUTING.md
[stdlib]: https://github.com/zelix-lang/stdlib

Welcome to the main repository of the [Zelix Programming Language][Website].
Here, you will find the source code for the Zelix compiler.
The code for the standard library can be found in the [stdlib] repository.

### ⚡ What is Zelix?

Zelix is an imperative, statically-typed programming language
designed with performance and simplicity in mind.

### ⚡ Why Zelix?

- **⚡ Performance** — Designed for speed, with aggressive low-level optimizations.
- **🛡 Memory safety** — Guaranteed at compile time — no garbage collector, no borrow checker headaches.
- **✍ Simplicity** — A clean, approachable syntax that gets out of your way.
- **🔍 Static typing** — Catch type errors early and get better performance, all at compile time.

### 📦 Installation

Before you begin, please make sure you have the following prerequisites installed:
- A C and C++ compiler (GCC or Clang are recommended).
- CMake (version 3.16 or higher), can be uninstalled after Zelix is installed.
- A Git client (for cloning the repository)
- **LLVM (version 16 or higher).**

Installation is simple. There are three main ways to install Zelix:

1. **Build from source**: Clone the repository and run the build script.
   ```bash
   git clone https://github.com/zelix-lang/Zelix.git
   cd Zelix
   ./build.sh
   mv ./build/zelix /usr/local/bin/zelix
    ```
2. **Download and run the automatic installer**:
   ```bash
   git clone https://github.com/zelix-lang/installer.git
   cd installer
   ./install.sh
   ```
3. **Use a package manager**: Zelix is available on various package managers.
   - **Homebrew**: 
    ```bash
    brew tap zelix-lang/installer-brew
    brew install zelix-lang
    ```
    - **AUR**:
    ```bash
    yay -S zelix-lang
    ```
    - **Debian/Fedora-based distributions and Windows**:
    Unfortunately, Zelix is not available on these package managers yet.
    You can use the **build from source** method instead.

### 👾 Supported Platforms & Architectures
Zelix can run on almost any platform that supports a C compiler.
Currently, Windows support is not available, and it is not planned for the near future.
Zelix is primarily developed and tested on Linux and macOS.

### 🤝 Getting Involved
We welcome contributions! Whether you're fixing bugs, adding features, or improving documentation, your help is
appreciated. Check out our [Contributing Guide][Contributing] for more details.

If you consider Zelix has inspired you to create something amazing,
you can read the [Zelix Internals Book](https://docs.zelixlang.dev/zelix-internals-book)
to learn how Zelix works under the hood and how it was implemented and written from scratch.

### 📚 Documentation
For detailed documentation, including language features, standard library usage, and more,
please visit our [Documentation][Documentation] page.

### 📝 License
Zelix is licensed under the [GNU General Public License v3.0](LICENSE).
This means you can use, modify, and distribute it freely, as long as you keep the
same license for your modifications and distributions.