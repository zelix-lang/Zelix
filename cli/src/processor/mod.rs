use std::{collections::HashMap, process::exit};

use logger::{Logger, LoggerImpl};

use crate::structs::surf_config_file::SurfConfigFile;

fn build_bind_trace(name:&String) -> String {
    format!(
        "At Surf.yml: bind -> {}",
        name
    )
}

/// Parses and processes the bindings
/// of the configuration file.
pub fn process_bindings(config_file: &SurfConfigFile) -> HashMap<String, String> {
    let mut result = HashMap::new();

    for binding in &config_file.bind {
        let name = binding.name.clone();
        let mut value = binding.value.clone();

        if value.starts_with("$self.bind") {
            Logger::err(
                "Cannot bind another binding",
                &[
                    "Sometimes you don't need more than 1 bindings",
                    "that point to the same value, as all bindings",
                    "are available globally."
                ],
                &[
                    build_bind_trace(&name).as_str()
                ]
            );

            exit(1);
        } else if value.starts_with("$self.") {
            value = match value.as_str() {
                "$self.name" => {
                    config_file.name.clone()
                }

                "$self.description" => {
                    config_file.description.clone()
                }

                "$self.version" => {
                    config_file.version.clone()
                }

                "$self.author" => {
                    config_file.author.clone()
                }

                "$self.license" => {
                    config_file.license.clone()
                }

                "$self.git" => {
                    config_file.git.clone()
                }

                "$self.main_file" => {
                    config_file.main_file.clone()
                }
                
                _ => {
                    Logger::err(
                        "Invalid binding",
                        &[
                            "You can only bind values that are part",
                            "of the project's metadata",
                            "If you want to bind custom values that aren't",
                            "part of the metadata, use strings that don't start",
                            "with '$self.'"
                        ],
                        &[
                            build_bind_trace(&name).as_str()
                        ]
                    );
                    exit(1);
                }
            }
        }

        result.insert(
            name,
            value
        );
    }

    result
}