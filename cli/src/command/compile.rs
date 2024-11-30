use shared::code::import::{Import, Importable};
use core::transpiler::transpile::transpile;
use std::{env::current_dir, fs::{remove_dir_all, remove_file}, os::unix::process::ExitStatusExt, path::PathBuf, process::{exit, ExitStatus}};

use shared::{logger::{Logger, LoggerImpl}, message::print_header, path::retrieve_path, result::try_unwrap, token::token::TokenImpl};

use crate::command::lexe_base::lexe_base;

pub fn compile_command(path: Option<PathBuf>) -> PathBuf {
    let preferred_compiler = try_unwrap(
        std::env::var("SURF_PREFERRED_COMPILER"),
        "Failed to get the compiler"
    );

    if 
        preferred_compiler != "clang++"
        && preferred_compiler != "g++"
    {
        Logger::err(
            "Invalid compiler",
            &[
                "Set the SURF_PREFERRED_COMPILER environment variable to clang++ or g++",
                "depending on your preferred compiler"
            ],
            &[
                "The compiler you set is invalid",
                "Please use either clang++ or g++"
            ],
        );

        exit(1);
    }

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


    // The path of the transpiled file will always be
    // join(out_dir, "out.cpp")
    // we get the imports so we know what to include

    // Determine if we're on Windows or Unix so we can add the correct extension
    let on_windows = cfg!(target_os = "windows");

    let transpiled_file_path = out_dir.join("out.cpp");
    let mut command = format!(
        "{} -o ",
        preferred_compiler
    );
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

    let mut files_to_link : Vec<PathBuf> = vec![];

    // Add all the imports
    // We import only .h and .hpp files
    // Some libraries don't have a corresponding .cpp file

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

    for import in imports {
        let import_path = import.get_from().clone();

        let extension_optional = import_path.extension();
        if extension_optional.is_none() {
            continue;
        }

        let prev_extension = extension_optional.unwrap().to_str().unwrap();
        // Some libraries may include only .hpp files
        // And don't provide a .cpp file
        if prev_extension != "h" && prev_extension != "hpp" {
            continue;
        }

        // Check if there's a .cpp file to build
        // otherwise, skip as the necessary code is
        // already in the header file
        let new_path = import_path.with_extension("cpp");

        if !new_path.exists() {
            // Nothing to do, the code is already in the header file
            continue;
        }

        // Convert the library to an object file
        let out_path = out_dir.join(
            new_path.with_extension("o")
                .file_name()
                .unwrap()
        );

        Logger::log(&[
            format!(
                "  <black_bright>-></black_bright> <magenta_bright>Building {}</magenta_bright>",
                new_path
                    .with_extension("")
                    .file_name().unwrap().to_str().unwrap()
            ).as_str()
        ]);

        let command = format!(
            "{} -c {} -o {}",
            preferred_compiler,
            new_path.to_str().unwrap(),
            out_path.to_str().unwrap()
        );

        let output = try_unwrap(
            std::process::Command::new(shell_program)
                .arg(separator)
                .arg(command)
                .output(),
            "Failed to execute the command",
        );

        if output.status != ExitStatus::from_raw(0) {
            Logger::err(
                "Failed to compile a library/import",
                &[
                    "Maybe your syntax is incorrect?",
                    "If this is a persistent issue, please report it on the GitHub repository",
                ],
                &[
                    "The compiler threw an error while compiling the code"
                ]
            );

            println!("{}", String::from_utf8_lossy(&output.stderr));
            
            exit(1);
        }

        files_to_link.push(out_path);
    }

    // Build the .cpp file
    {
        let path_with_new_extension = transpiled_file_path.with_extension("o");
        let transpiled_code_path = path_with_new_extension
            .to_str().unwrap();

        Logger::log(&[
            format!(
                "  <black_bright>-></black_bright> <magenta_bright>Compiling {}</magenta_bright>",
                path_with_new_extension.file_name().unwrap().to_str().unwrap()
            ).as_str()
        ]);

        let command = format!(
            "{} -o {} -c {}",
            preferred_compiler,
            transpiled_code_path,
            transpiled_file_path.to_str().unwrap()
        );

        let output = try_unwrap(
            std::process::Command::new(shell_program)
                .arg(separator)
                .arg(command)
                .output(),
            "Failed to execute the command",
        );

        if output.status != ExitStatus::from_raw(0) {
            Logger::err(
                "Failed to compile the main file",
                &[
                    "Maybe your syntax is incorrect?",
                    "If this is a persistent issue, please report it on the GitHub repository",
                ],
                &[
                    "The compiler threw error while compiling the code"
                ]
            );

            println!("{}", String::from_utf8_lossy(&output.stderr));
            
            exit(1);
        }

        files_to_link.push(PathBuf::from(transpiled_code_path));
    }

    for file in files_to_link {
        command.push_str(file.to_str().unwrap());
        command.push_str(" ");
    }

    // Execute the command
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
                "The compiler threw error while compiling the code"
            ]
        );

        println!("{}", String::from_utf8_lossy(&output.stderr));
        
        exit(1);
    } else {
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