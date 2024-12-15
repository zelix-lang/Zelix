use std::{fs::{self, create_dir, create_dir_all}, path::PathBuf, process::exit};

use globals::{BOOL_REGEX, DESC_REGEX, NAME_REGEX, ROCKET_EMOJI, VERSION_REGEX};
use logger::{Logger, LoggerImpl};
use shared::{message::print_header, stdin::question::question};
use util::{path::retrieve_path, result::try_unwrap};
use std::{env::current_dir, fs::exists};

use crate::{example::example_program::{EXAMPLE_GIT_IGNORE, EXAMPLE_PROGRAM}, structs::surf_config_file::{SurfBinds, SurfConfigFile}};

mod globals {
    use fancy_regex::Regex;
    use lazy_static::lazy_static;

    lazy_static! {
        pub static ref ROCKET_EMOJI : String = emojis::get_by_shortcode("rocket").unwrap().as_str().to_string();
        pub static ref NAME_REGEX : Regex = Regex::new(r"^[a-zA-Z\-_]+((\d+)?)$").unwrap();
        pub static ref DESC_REGEX : Regex = Regex::new(r".+").unwrap();
        pub static ref VERSION_REGEX : Regex = Regex::new(r"^\d+\.\d+\.\d+$").unwrap();
        pub static ref AUTHOR_REGEX : Regex = Regex::new(r"^[a-zA-Z\.\s0-9]+$").unwrap();
        pub static ref BOOL_REGEX : Regex = Regex::new(r"^[yYnN]$").unwrap();
    }
}

pub fn init_command(path: Option<PathBuf>) {
    let cwd = try_unwrap(
        current_dir(),
        "Failed to get the current directory"
    );

    let final_path = retrieve_path(
        path.clone().unwrap_or(cwd.clone())
    );

    if !try_unwrap(exists(final_path.clone()), "Failed to check if the current dir exists") {
        try_unwrap(
            create_dir_all(final_path.clone()),
            "Failed to create the directory at the given path"
        );
    } else {
        let files = try_unwrap(
            std::fs::read_dir(final_path.clone()),
            "Failed to read the directory"
        );

        if files.count() > 0 {
            Logger::err(
                &"The directory is not empty",
                &[&"You can create an empty directory or use a new path"],
                &[
                    format!(
                        "Directory: {}",
                        final_path.to_str().unwrap()
                    ).as_str()
                ]
            );

            exit(1);
        }
    }

    print_header();
    Logger::log(&[
        "<black_bright>You may cancel this process at any time by pressing Ctrl + C.</black_bright>",
        "<black_bright>This will not affect the directory in any way unless you confirm the operation.</black_bright>",
        // Add a newline here to separate the message from the prompt
        ""
    ]);

    let name = question("Project name", "Surf-Project", NAME_REGEX.clone());
    let description = question("Project description", "A new Surf project", DESC_REGEX.clone());
    let version = question("Project version", "0.1.0", VERSION_REGEX.clone());
    let author = question("Project author", "John Doe", DESC_REGEX.clone());
    let license = question("Project license", "MIT", DESC_REGEX.clone());
    let git_repo = question("Git repository (None for none)", "None", DESC_REGEX.clone());
    let init_git_repo = question("Initialize a git repository [Y/n]", "Y", BOOL_REGEX.clone()).to_lowercase();

    // Add a newline here to separate the prompt from the message
    println!();

    if init_git_repo == "y" {
        try_unwrap(
            std::process::Command::new("git")
                .arg("init")
                .current_dir(final_path.clone())
                .spawn(),
            "Failed to initialize the git repository (Is git installed?)"
        );

        try_unwrap(
            fs::write(
                final_path.join(".gitignore"),
                EXAMPLE_GIT_IGNORE.clone()
            ),
            "Failed to write the main file"
        );
    }

    let config_struct = SurfConfigFile {
        name,
        description,
        version,
        author,
        license,
        git: git_repo,
        main_file: "src/main.surf".to_string(),
        repositories: vec![],
        dependencies: vec![],
        bind: vec![
            SurfBinds {
                name: "project_name".to_string(),
                value: "$self.name".to_string()
            },
            SurfBinds {
                name: "project_version".to_string(),
                value: "$self.version".to_string()
            }
        ]
    };

    let config_string = serde_yml::to_string(&config_struct).unwrap();
    Logger::log(&[
        "The following contents are going to be written to Surf.yml",
        "",
        config_string.as_str()
    ]);

    let confirm = question("Is this okay? [Y/n]", "Y", BOOL_REGEX.clone()).to_lowercase();

    if !confirm.eq("y") {
        Logger::log(&["<red_bright>The process was cancelled</red_bright>"]);
        exit(0);
    }

    let src_dir = final_path.join("src");

    try_unwrap(
        create_dir(src_dir.clone()),
        "Failed to create the src directory"
    );
    
    try_unwrap(
        fs::write(
            final_path.join("Surf.yml"),
            format!(
                "{}\n{}\n{}\n{}\n\n{}",
                "# ---------------------------------------------------", 
                "#   This file is generated by the Surf CLI.", 
                "#   Feel free to modify the file to suit your needs.",
                "# ----------------------------------------------------",
                config_string
            )
        ),
        "Failed to write the main file"
    );

    try_unwrap(
        fs::write(
            src_dir.join("main.surf"),
            EXAMPLE_PROGRAM.clone()
        ),
        "Failed to write the main file"
    );

    Logger::log(&[
        "",
        format!(
            "{}{}",
            ROCKET_EMOJI.clone(),
            " Your project is ready!"
        ).as_str(),
        "",
        "  Next steps:",
        ""
    ]);

    let print_cd_step = cwd != final_path;
    let next_step: i32 = if print_cd_step { 2 } else { 1 };

    if print_cd_step {
        Logger::log(&[
            format!(
                "    <magenta_bright>1. cd {}</magenta_bright>",
                path.unwrap().to_str().unwrap()
            ).as_str()
        ]);
    }

    Logger::log(&[
        format!(
            "    <magenta_bright>{}. surf run</magenta_bright>",
            next_step
        ).as_str(),
        "",
        "  Happy coding!"
    ]);

}