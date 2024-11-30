use std::fs::read_to_string;

use checker::import_checker::check_imports;
use shared::{code::{import_extractor::{extract_import_matches, IMPORT_REGEX}, import_path_extractor::extract_imports_paths}, result::try_unwrap};

// Returns the modified code with the imports removed
// This is useful so the lexer can process the imports in one go
// instead of calling the lexer each time a new import is found
// That would give an unnecessary complexity of O(n^2)
pub fn process_imports(raw_code: &String, file: &String) -> String {
    let mut result = raw_code.clone();

    while try_unwrap(
        IMPORT_REGEX.is_match(&result),
        "Failed to check if the code contains imports"
    ) {
        let matches: Vec<String> = extract_import_matches(result.clone());
        let matches_paths : Vec<String> = extract_imports_paths(&matches);

        // Let the import checker do its job
        check_imports(&matches_paths, file);

        for import_path in matches_paths {

            // Not checking for empty imports because the regex pattern
            // matches full imports only, empty imports are processed
            // as tokens, this will cause the program to refuse to compile
            // because you can't have anything other than imports
            // outside of a function

            // Read the file contents
            // The import checker also checks if the file exists
            let import_content = read_to_string(&import_path).unwrap();
            result = result.replacen(import_path.as_str(), &import_content, 1);
        }

    }

    result
}