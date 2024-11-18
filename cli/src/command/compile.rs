use core::compiler::compile::compile;
use std::{env::current_dir, fs::{remove_dir_all, remove_file}, path::PathBuf};

use shared::{path::retrieve_path, result::try_unwrap};

use crate::command::lexe_base::lexe_base;

pub fn compile_command(path: Option<PathBuf>) {
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

    compile(tokens, out_dir);
}