use core::{shared::import::{Import, Importable}, transpiler::transpile::transpile};
use std::{env::current_dir, fs::{remove_dir_all, remove_file}, os::unix::process::ExitStatusExt, path::PathBuf, process::{exit, ExitStatus}};

use shared::{logger::{Logger, LoggerImpl}, message::print_header, path::retrieve_path, result::try_unwrap, token::token::TokenImpl};

use crate::command::lexe_base::lexe_base;

pub fn compile_command(path: Option<PathBuf>) -> PathBuf {
    print_header();
    println!();

    let tokens = lexe_base(path.clone());
    let cwd = try_unwrap(
        current_dir(),
        "Failed to get current working directory",
    );

    let final_path = retrieve_path(
        path.unwrap_or(cwd.clone())
    );

    let out_dir = final_path.join("out");

    if out_dir.exists() {
        if out_dir.is_dir() {
            try_unwrap(
                remove_dir_all(out_dir.clone()),
                "Failed to remove out directory",
            );
        } else {
            try_unwrap(
                remove_file(out_dir.clone()),
                "Failed to remove out file",
            );
        }
    }

    try_unwrap(
        std::fs::create_dir(out_dir.clone()),
        "Failed to create out directory",
    );

    Logger::log(&[
        "  <black_bright>-></black_bright> <magenta_bright>Building</magenta_bright>"
    ]);

    let source : PathBuf;
    if tokens.len() > 0 {
        source = PathBuf::from(tokens[0].get_file());
    } else {
        source = final_path.clone();
    }

    let imports: Vec<Import> = transpile(
        tokens,
        out_dir.clone(),
        source
    );

    Logger::log(&[
        "  <black_bright>-></black_bright> <magenta_bright>Compiling</magenta_bright>"
    ]);

    // The path of the transpiled file will always be
    // join(out_dir, "out.cpp")
    // we get the imports so we know what to include

    // Determine if we're on Windows or Unix so we can add the correct extension
    let on_windows = cfg!(target_os = "windows");

    let transpiled_file_path = out_dir.join("out.cpp");
    let mut command = String::from("clang++ -o ");
    let mut final_out_path = out_dir.join("out");

    // Add the output file
    command.push_str(
        final_out_path.to_str().unwrap()
    );

    if on_windows {
        final_out_path.set_extension("exe");
        command.push_str(".exe");
    }

    command.push_str(" ");

    // Add the main file to the command
    command.push_str(transpiled_file_path.to_str().unwrap());

    // Add all the imports
    // We import only .h files, but the standard library is built
    // in such a way that for each .h file, there's a corresponding
    // .cpp file, so we can just replace the .h extension with .cpp

    for import in imports {
        let mut import_path = import.get_from().clone();
        import_path.set_extension("cpp");

        command.push_str(" ");
        command.push_str(import_path.to_str().unwrap());
    }

    // Execute the command
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

    let output = try_unwrap(
        std::process::Command::new(shell_program)
            .arg(separator)
            .arg(command)
            .output(),
        "Failed to execute the command",
    );

    if output.status != ExitStatus::from_raw(0) {
        Logger::err(
            "Failed to compile the code",
            &[
                "Maybe your syntax is incorrect?",
                "If this is a persistent issue, please report it on the GitHub repository",
            ],
            &[
                "Clang++ gave an error while compiling the code"
            ]
        );

        println!("{}", String::from_utf8_lossy(&output.stderr));
        
        exit(1);
    } else {
        // Remove the .cpp file
        try_unwrap(
            remove_file(transpiled_file_path.clone()),
            "Failed to remove the transpiled file",
        );

        Logger::log(&[
            "  <black_bright>-></black_bright> <green_bright>Finished</green_bright>",
            format!(
                "    <black_bright>-> {}</black_bright>",
                final_out_path.to_str().unwrap()
            ).as_str(),
        ]);
    }

    final_out_path

}