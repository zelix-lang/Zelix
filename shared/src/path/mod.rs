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

    let cwd_string = cwd.to_str().unwrap();
    
    path.replace(cwd_string, "")
}