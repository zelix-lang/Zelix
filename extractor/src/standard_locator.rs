use std::{fs::exists, path::PathBuf, process::exit};

use clang::Index;
use logger::{Logger, LoggerImpl};
use shared::{code::{file_code::{FileCode, FileCodeImpl}, import::{Import, Importable}}, env::STANDARD_LIBRARY_LOCATION};
use util::result::try_unwrap;

pub fn locate_and_import_package(name: &str, result: &mut FileCode, trace: String, index: &Index) {
    // Let the standard locator locate and validate the package
    let package_path = locate_standard(name.to_string());

    if !package_path.is_dir() {
        Logger::err(
            "Failed to locate package",
            &[
                "Make sure the package is a directory"
            ],
            &[
                format!(
                    "The path {} is not a directory",
                    package_path.to_str().unwrap()
                ).as_str()
            ]
        );

        exit(1);
    };

    // Scan the package for files
    let files = try_unwrap(
        package_path.read_dir(),
        "Failed to read package directory"
    );

    for file in files {
        let file = try_unwrap(
            file,
            "Failed to read package directory"
        );

        let file_path = file.path();
        let file_name = file_path.file_name().unwrap().to_str().unwrap().to_string();

        // Only import .hpp and .h files, .cpp files should be ignored
        if file_name.ends_with(".cpp") {
            continue;
        }

        result.add_import(
            Import::new(
                locate_standard(
                    format!(
                        "{}/{}",
                        name,
                        file_name
                    )
                ),
                trace.clone()
            ),
            index
        );
    }
}

pub fn locate_standard(name: String) -> PathBuf {

    let standard_path = STANDARD_LIBRARY_LOCATION.clone();
    let final_path = standard_path.join(name.clone());
    let exists_result = try_unwrap(
        exists(final_path.clone()),
        "Failed to locate standard libraries!"
    );

    if !exists_result {
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