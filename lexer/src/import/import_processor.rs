use std::{collections::HashSet, fs::read_to_string, path::PathBuf, process::exit};

use shared::{
    code::{
        import_extractor::{extract_import_matches, IMPORT_REGEX},
        import_path_extractor::extract_imports_paths,
    },
    logger::{Logger, LoggerImpl},
    result::try_unwrap,
};

// Builds the chain trace for a single element of the chain
fn build_single_chain(n: usize, import: &PathBuf, chain: &mut Vec<String>) {
    // It's safe to unwrap here because we're iterating through the set
    // therefore the index will always be valid
    let mut string_to_be_added = String::new();

    string_to_be_added.push_str(
        " ".repeat(n).as_str()
    );

    string_to_be_added.push_str(
        "-> "
    );

    string_to_be_added.push_str(
        import.to_str().unwrap()
    );

    chain.push(string_to_be_added);
}

// Build the chain of imports to print in the error message
fn build_chain(imports: &HashSet<PathBuf>, root: &PathBuf) -> Vec<String> {
    let mut chain: Vec<String> = Vec::new();

    build_single_chain(0, root, &mut chain);

    for n in 0..imports.len() {
        let import = imports.iter().nth(n).unwrap();
        build_single_chain(n + 1, import, &mut chain);
    }

    chain
}

pub fn process_imports(raw_code: &String, file: &PathBuf) -> String {
    let mut result = raw_code.clone();

    // Queue for files to process
    let mut work_queue: Vec<(String, PathBuf)> = vec![(raw_code.clone(), file.clone())];

    // Track seen imports to avoid circular dependencies
    let mut seen_imports: HashSet<PathBuf> = HashSet::new();

    // Iterate through the queue until all files are processed
    while let Some((current_code, current_file)) = work_queue.pop() {
        // Track the current directory context
        let current_context = current_file.parent().unwrap_or_else(|| {
            Logger::err(
                "Invalid file path",
                &["The file path is invalid"],
                &["Make sure the file path is correct"],
            );
            exit(1);
        }).to_path_buf();

        while try_unwrap(
            IMPORT_REGEX.is_match(&current_code),
            "Failed to check if the code contains imports",
        ) {
            let matches: Vec<String> = extract_import_matches(current_code.clone());
            let matches_paths: Vec<PathBuf> =
                extract_imports_paths(&matches, &current_context);

            for n in 0..matches.len() {
                let import_match = matches[n].clone();
                let import = matches_paths[n].clone();

                if seen_imports.contains(&import) {
                    Logger::err(
                        "Circular dependency detected",
                        &["Make sure the imports are correct"],
                        build_chain(&seen_imports, file)
                        .iter()
                        .map(|a| a.as_str()).collect::<Vec<&str>>()
                        .as_slice(),
                    );
                    exit(1);
                }

                seen_imports.insert(import.clone());

                let import_code = match read_to_string(&import) {
                    Ok(code) => code,
                    Err(_) => {
                        Logger::err(
                            "Failed to read import",
                            &["Make sure the import path is correct"],
                            &[
                                "Failed to read the import file",
                                format!("Import path: {}", import.to_str().unwrap())
                                    .as_str(),
                                format!("At {}", file.to_str().unwrap()).as_str(),
                            ],
                        );
                        exit(1);
                    }
                };

                // Replace the import match with a placeholder for now
                result = result.replacen(&import_match, "", 1);

                // Enqueue the imported file for processing
                work_queue.push((import_code, import));
            }
        }
    }

    result
}
