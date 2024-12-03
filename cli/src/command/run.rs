use std::path::PathBuf;

use logger::{Logger, LoggerImpl};

use super::compile::compile_command;

pub fn run_command(path: Option<PathBuf>) {

    let output_file = compile_command(path);

    Logger::log(&[
        "  <black_bright>-></black_bright> <green_bright>Running</green_bright>",
        ""
    ]);

    let on_windows = cfg!(target_os = "windows");
    let shell_program = if on_windows {
        "cmd"
    } else {
        "sh"
    };

    let separator = if on_windows {
        "/C"
    } else {
        "-c"
    };

    let mut command = std::process::Command::new(shell_program);

    command
        .arg(separator)
        .arg(output_file.to_str().unwrap());

    command
        .stdin(std::process::Stdio::inherit())
        .stdout(std::process::Stdio::inherit())
        .stderr(std::process::Stdio::inherit());

    // Execute the command
    let status = command.spawn().expect("Failed to execute command")
        .wait().expect("Failed to wait on child");

    if !status.success() {
        eprintln!("Command executed with failing error code");
    }
}