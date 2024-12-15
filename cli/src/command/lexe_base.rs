use std::{collections::HashMap, fs::read_to_string, path::PathBuf, process::exit};

use code::token::Token;
use lexer::{Lexer, LexerImpl};
use logger::{Logger, LoggerImpl};
use util::{path::retrieve_path, result::try_unwrap};

use std::{
    fs::{exists, metadata},
    env::current_dir
};

use crate::{processor::process_bindings, structs::surf_config_file::SurfConfigFile};

pub fn lexe_base(path: Option<PathBuf>) -> (Vec<Token>, HashMap<String, String>, PathBuf) {
    let final_path =
        retrieve_path(
            path.unwrap_or(
                try_unwrap(
                    current_dir(), 
                    "Failed to get current directory"
                )
            )
        );

    if !try_unwrap(
        exists(final_path.clone()),
        "Failed to check if the current dir exists",
    ) {
        Logger::err(
            &"The path doesn't exist!",
            &[&"Make sure the path is correct"],
            &[&"No trace available"],
        );
        exit(1);
    }

    if !try_unwrap(
        metadata(final_path.clone()),
        "Failed to get metadata from the path",
    )
    .is_dir()
    {
        Logger::err(
            &"The path is not a directory!",
            &[
                &"Make sure the path is correct",
                &"Start a new project with surf init",
            ],
            &[&"No trace available"],
        );
        exit(1);
    }

    let config_file_path = final_path.join("Surf.yml");

    if !config_file_path.exists() {
        Logger::err(
            &"The path is not a Surf project!",
            &[
                &"Make sure the path is correct",
                &"Start a new project with surf init",
            ],
            &[&"No trace available"],
        );
        exit(1);
    }

    let config_file: SurfConfigFile = try_unwrap(
        serde_yml::from_str(
            try_unwrap(
                read_to_string(config_file_path),
                "Failed to read the Surf.yml file",
            )
            .as_str(),
        ),
        "Failed to parse the Surf.yml file",
    );

    let main_file = final_path.join(&config_file.main_file);

    if !main_file.exists() {
        Logger::err(
            &"The main file doesn't exist!",
            &[
                &"Make sure the path is correct",
                &"Start a new project with surf init",
            ],
            &[&"No trace available"],
        );
        exit(1);
    }

    let mut main_file_content = try_unwrap(
        read_to_string(main_file.clone()),
        "Failed to read the main file",
    );

    let bindings = process_bindings(&config_file);
    
    (
        Lexer::new().tokenize(&mut main_file_content, &main_file),
        bindings,
        final_path
    )
}
