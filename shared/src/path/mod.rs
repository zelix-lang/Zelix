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