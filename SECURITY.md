<div align="center">
    <img src="https://assets.zelixlang.dev/logo.png?update=true" height="60" width="60">
    <h1>The Zelix Programming Language</h1>
    Zelix a modern and blazing-fast programming language.
</div>

---

## 👋 Welcome

Welcome to the official Zelix repository! In this file, you'll find a guide on how to
report security vulnerabilities in the Zelix programming language.

---

## ⛔ What is forbidden

Do NOT open an issue on GitHub for a security vulnerability.
The issue will be immediately closed and deleted.
This is done to prevent further exploitation of the vulnerability.

---

## 📝 The basics

Before reporting a security vulnerability, make sure to properly
test it in various environments and configurations (Windows, Linux, macOS, etc.);
(VMs are recommended for this purpose). This will help us understand the
scope of the vulnerability and its impact on the language.

**Before you proceed, please make sure that the vulnerability
is NOT tied to:**

- Having ASLR or any other security feature disabled.
- Stack smashing due to the previous point.
- A bug that has no ability to propagate with malicious intent in a targeted environment. **[3]**
- A bug that requires physical access to the machine.
- A bug that requires the attacker to be a privileged user (Unless there is a privilege escalation vulnerability).

> In case of **[3]**: You might want to open a public issue instead if this bug is NOT malicious, has no security
implications, is not a vulnerability, and modifies the intended behavior of the language (e.g. crashes, incorrect output, etc.).

**What will NOT be fixed:**

- Security vulnerabilities that are not exploitable.
- Security vulnerabilities that come from the OS itself and not from the language's code.
- Any security vulnerability that is the result of the points mentioned in the **NOT tied to** section.

> In case of a vulnerability that is caused by a third-party library, please report it to the library's maintainers.
> Zelix will update this library as soon as they release a fix, please do not report this kind of vulnerabilities to us.

---

## 🚨 Reporting a Vulnerability

If you have found a vulnerability that meets the [above requirements](#-what-is-forbidden),
please send a detailed report to: fluent-security.purifier127@passfwd.com

**What your report must include:**

- Detailed steps to reproduce the vulnerability.
- The environment in which the vulnerability was tested.
- The impact of the vulnerability on a scale of 1 to 10.
- Any additional information that might be helpful.
- A PoC (Proof of Concept) if possible.
- Your name and email address if you want us to credit you (optional).
- All the previous points with clarity, do not add riddles or puzzles. Doing so will get your message ignored.

Please, do not email this address for any other reason than reporting a security vulnerability.
Doing so will result in your email being blocked.

**Additional information and FAQ:**

1. **What is the disclosure policy?**
   - The vulnerability will be disclosed to the public after a fix is released.
   - During the fix, you will see no commits regarding the vulnerability, as this will expose the vulnerability to the public.
   - If the vulnerability is of high severity, a new minor version will be released as soon as the fix is ready.
   - If the vulnerability is of moderate or low severity, the fix will be included in the next minor version.
   - You will be credited in the release notes if you provided your name and email address.

2. **Can I choose to disclose my name but not my email address?**
   - Yes, you can choose to disclose your name only. To do so, you must specify this clearly in the message that you send.

3. **Can I choose to disclose my email address but not my name?**
   - No, you must disclose your name if you want to disclose your email address. This is done to ensure that the vulnerability reporter is a real person.

4. **Can I choose to disclose neither my name nor my email address?**
   - Yes, you can choose to disclose neither your name nor your email address. In this case, you will not be credited in the release notes.

5. **How can I measure the vulnerability's impact on a scale of 1 to 10?**
   - The impact is measured by the potential damage that the vulnerability can cause.
   - An issue that can cause incorrect terminal input or CLI crashes is considered a **BUG**, not a vulnerability.
   - A vulnerability that gives an attacker any control over (RCE, privilege escalation, reading files, changing files, modifying system behavior, etc.) a foreign machine is considered a **10**.
   - A vulnerability that allows DoS (albeit highly unlikely with Zelix) or network attacks is considered a **5**.
   - A vulnerability that triggers an infinite loop and resource exhaustion is considered a **5**. **[3]**
   - A vulnerability that allows information disclosure in a negligent manner is considered a **1**. (E.g. Compiler versions in an error message).
   - A vulnerability that is caused by a third-party library is **initially** considered a **4**, if the vulnerability is severe, we might develop a custom solution or switch to another library.

    > In the case of **[3]**: This kind of vulnerability will only be considered if it can be triggered remotely, open a public issue otherwise.

6. **What kind of vulnerabilities are non-exploitable?**
   - A vulnerability that requires the attacker to have physical access to the machine.
   - A vulnerability that requires the attacker to be a privileged user.
   - A vulnerability that requires the attacker to disable security features.
   - A vulnerability that only occurs if the attacker modifies Zelix's source code (i.e., self-inflicted security risk).

7. **How long can it take for my message to get a response?**
   - It can take up to 48 hours for your message to get a response.
   - If you do not get a response within 48 hours, please send your message again if your report meets all the requirements listed above.

8. **How long can it take for any vulnerability to be disclosed?**
    - We will get back to you when you report a vulnerability. We will mention an estimate based on the severity and the time it will take to fix it.