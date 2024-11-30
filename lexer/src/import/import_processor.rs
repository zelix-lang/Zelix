use shared::logger::{Logger, LoggerImpl};

use super::import_extractor::extract_import_matches;

// Returns the modified code with the imports removed
// This is useful so the lexer can process the imports in one go
// instead of calling the lexer each time a new import is found
// That would give an unnecessary complexity of O(n^2)
pub fn process_imports(raw_code: &String, file: &String) -> String {
    let mut result = raw_code.clone();
    let matches: Vec<&str> = extract_import_matches(&raw_code);
    let mut processed_imports: Vec<String> = Vec::new();

    for raw_import in matches {
        let import_info = raw_import.replacen("@import", "", 1);
        // Remove the quotation marks and semicolon
        // Removes the first and last 2 characters
        let import_path = import_info[1..import_info.len()-2].to_string();

        // Not checking for empty imports because the regex pattern
        // matches full imports only, empty imports are processed
        // as tokens, this will cause the program to refuse to compile
        // because you can't have anything other than imports
        // outside of a function

        // Check for duplicate imports
        if processed_imports.contains(&import_path) {
            // Not a critical mistake, just skip the import and leave a warning
            Logger::warn(
                "Duplicate import found",
                &[
                    format!("The import {} was already processed", import_path)
                        .as_str(),
                    format!("At: {}", file)
                    .as_str()
                ]
            );
            continue;
        }

        // Imports of the standard library are later processed
        // by the lexer, so we don't need to process them here
        if !import_path.starts_with("@Surf:standard") {
            
        }

        // Save the import
        processed_imports.push(import_path);
    }

    result
}