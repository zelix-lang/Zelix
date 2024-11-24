use std::{env::var, fs::exists, path::PathBuf, process::exit};

use shared::logger::{Logger, LoggerImpl};

pub fn locate_standard(name: String) -> PathBuf {

    let standard_path_result = var("SURF_STANDARD_PATH");
    
    if standard_path_result.is_err() {
        Logger::err(
            "Failed to locate standard libraries!",
            &[
                "Make sure the SURF_STANDARD_PATH environment variable is set",
                "Try reinstalling Surf CLI"
            ],
            &[
                "No trace available"
            ]
        );

        exit(1);
    }

    let standard_path = PathBuf::from(standard_path_result.unwrap());
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