use std::path::PathBuf;

use crate::result::try_unwrap;

pub fn retrieve_path(path: PathBuf) -> PathBuf {
    let cwd = try_unwrap(
        std::env::current_dir(),
        "Failed to get current directory"
    );

    if path.is_absolute() {
        return path;
    }

    return cwd.join(path);
}

pub fn discard_cwd(path: String) -> String {
    let cwd = try_unwrap(
        std::env::current_dir(),
        "Failed to get current working directory",
    );

    let mut cwd_string = cwd.to_str().unwrap().to_string();

    // Add an "/" to the end of the cwd string if it doesn't have one
    if !cwd_string.ends_with("/") {
        cwd_string.push_str("/");
    }

    
    path.replace(cwd_string.as_str(), "")
}