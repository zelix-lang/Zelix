use std::{collections::VecDeque, fs::read_to_string, path::PathBuf};

use shared::{code::{import::{Import, Importable}, import_extractor::extract_import_matches, import_path_extractor::extract_imports_paths}, logger::{Logger, LoggerImpl}};

// Check if Vec<&Import> contains a specific import
// based on the import path (PathBuf)
fn list_contains_import(list: &Vec<Import>, import: &PathBuf) -> Option<Import> {
    for i in list {
        if i.get_from() == *import {
            return Some(i.clone());
        }
    }

    None
}

pub fn check_imports(matches: &Vec<String>, file: &String) {
    // We'll use a queue to keep track of the imports we need to check
    let mut chain : VecDeque<String> = VecDeque::new();
    // Initialize the queue with the imports we found
    chain.extend(matches.iter().cloned());

    // Store the imports we've seen
    let mut seen_imports : Vec<Import> = Vec::new();

    while !chain.is_empty() {
        // Unwrap here is safe because we check if the chain is empty
        let import = chain.pop_front().unwrap();

        let import_path = PathBuf::from(import.clone());

        // Check if we've seen this import before
        // If we have, we have a circular dependency
        let potential_circular = list_contains_import(&seen_imports, &import_path);
        if potential_circular.is_some() {
            // This may not be triggered due to duplicate imports
            // because we exclude them during the extraction
            // so we're ensured this will never happen
            Logger::err(
                "Circular dependency detected",
                &[
                    "Use a different module to separate the circular dependency"
                ],
                &[
                    potential_circular.unwrap().get_trace().as_str()
                ]
            );
        }

        // Let's check if the file exists
        if !import_path.exists() {
            Logger::err(
                "Import not found",
                &[
                    "The import you're trying to use doesn't exist"
                ],
                &[
                    "Make sure the file exists and the path is correct"
                ]
            );
        }

        // Add to the seen imports
        seen_imports.push(
            Import::new(
                import_path.clone(),
                file.clone()
            )
        );

        // Read the import's content and extract the imports
        // We'll add them to the queue

        // Using unwrap here is safe because we check if the file exists
        let import_content = read_to_string(import_path.clone()).unwrap();

        // Add the new imports to the queue
        let new_imports = extract_imports_paths(
            &extract_import_matches(
                import_content.clone()
            )
        );

        chain.extend(new_imports.iter().cloned());
    }
}