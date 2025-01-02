<div align="center">
    <img src="assets/logo.png" height="60" width="60">
    <h1>Fluent Language</h1>
    Fluent is a modern and blazing-fast programming language.
</div>

---

## 👋 Welcome

Welcome to the official Fluent repository! In this file, you'll find a guide on how to contribute to the Fluent language.

---

## 📝 The basics

Before contributing to the Fluent language, make sure you have read the [Code of Conduct](CODE_OF_CONDUCT.md).
On top of that, certain restrictions apply to the code contributed to the Fluent language, such as:

- You may not contribute code that is not yours, or wasn't given permission to use
- You may not contribute code that is harmful to the community
- **You may not attempt to pull obfuscated code into the Fluent language**
- **You may not attempt to pull malicious code into the Fluent language**
- **You may not attempt to pull code that is not related to the Fluent language**
- All contributed code should have an explanation of what it does and why it was added, code that is self-explanatory is an exception.

**Please, send any PRs to the `dev` branch, changes are only merged into the `main` branch after they have been reviewed and tested thoroughly.**

---

## 📝 Our standards

Before you jump into editing and/or pulling code into the Fluent language,
make sure your code follows our [Style Guide](STYLE_GUIDE.md),
which provides a set of rules and guidelines to follow when writing code for the Fluent language.

Also, please **Benchmark** your code before submitting it, we want
to ensure that the code you submit is as fast as possible.

This repository holds a built-in benchmarking tool that you can use to benchmark your code.
This tool will run the static analyzer on the file located inside the `example` directory and output the results.

To use it, you can use the following command in your terminal:

```bash
go run main.go -b
```

---

## 🚀 Getting Started

You may submit your code through a pull request. You can submit code independently or create a fork of the repository and submit your code through that fork.

To submit code through a fork, follow these steps:

1. Fork the repository
2. Clone the forked repository to your local machine
3. Create a new branch for your changes
4. Make your changes
5. Test your changes and ensure they work as expected
6. Commit your changes
7. Push your changes to your fork
8. Create a pull request on the `dev` branch