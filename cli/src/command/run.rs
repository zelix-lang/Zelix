use std::{path::PathBuf, process::exit};

use shared::{logger::{Logger, LoggerImpl}, result::try_unwrap};
use std::{env::current_dir, fs::exists};

pub fn run_command(path: Option<PathBuf>) {
    let final_path = path.unwrap_or(
        try_unwrap(
            current_dir(), 
            "Failed to get current directory"
        )
    );

    if !try_unwrap(
        exists(final_path.clone()),
        "Failed to check if the current dir exists"
    ) {
        Logger::err(
            &"The path doesn't exist!",
            &[&"Make sure the path is correct"],
            &[&"No trace available"]
        );
        exit(1);
    }
}