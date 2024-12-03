use std::{fs::exists, path::PathBuf, process::exit};

use logger::{Logger, LoggerImpl};
use shared::env::STANDARD_LIBRARY_LOCATION;

pub fn locate_standard(name: String) -> PathBuf {

    let standard_path = STANDARD_LIBRARY_LOCATION.clone();
    let final_path = standard_path.join(name.clone());
    let exists_result = exists(final_path.clone());

    if exists_result.is_err() || !exists_result.unwrap() {
        Logger::err(
            "Failed to locate standard libraries!",
            &[
                "Make sure the SURF_STANDARD_PATH environment variable is set",
                "Try reinstalling Surf CLI"
            ],
            &[
                format!(
                    "File {} does not exist",
                    final_path.to_str().unwrap()
                ).as_str()
            ]
        );

        exit(1);
    }

    final_path

}