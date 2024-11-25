// Standard library imports
use std::{collections::HashMap, fs::read_to_string, path::PathBuf, process::exit};

// External module imports
use lexer::{Lexer, LexerImpl};
use shared::{logger::{Logger, LoggerImpl}, path::discard_cwd, result::try_unwrap};

// Internal module imports
use crate::{extractor::extract_parts, shared::{file_code::{FileCode, FileCodeImpl}, import::{Import, Importable}}};

// Builds a human-readable trace of the import chain, showing dependencies in order
fn build_trace(import_chain : Vec<PathBuf>) -> Vec<String> {
    let mut chain : Vec<String> = Vec::new();

    // Loop through the chain of imports and format each level
    for n in 0..import_chain.len() {
        let mut el = String::from(" ".repeat(n).as_str());

        el.push_str("-> ");
        el.push_str(
            discard_cwd(
                import_chain[n].to_str().unwrap().to_string()
            ).as_str()
        );

        chain.push(el);
    }

    chain
}

// Checks whether the extension of the import file is valid
fn check_extension(import: &&Import) -> bool {
    let path = import.get_from();
    let extension_optional = path.extension();

    // If there is no extension, log an error and exit
    if extension_optional.is_none() {
        Logger::err(
            "Invalid import file extension!",
            &[
                "The import file must have an extension"
            ],
            &[path.to_str().unwrap()]
        );

        exit(1);
    }

    // Returns true if the extension is not `.h`
    extension_optional.unwrap() != "h" && extension_optional.unwrap() != "hpp"
}

// Analyzes the imports in a source file to check for issues like circular dependencies
pub fn analyze_imports(source: FileCode) {
    let imports = source.get_imports();
    // Collects all initial imports that pass the extension check
    let mut import_chain = imports
        .iter().filter(check_extension)
        .map(|i| i.get_from().clone())
        .collect::<Vec<PathBuf>>();

    // Initializes the lexer and tracking data structures
    let mut lexer: Lexer = Lexer::new();
    let mut checked_files: Vec<PathBuf> = Vec::new();
    let mut all_import_chains: HashMap<PathBuf, Vec<PathBuf>> = HashMap::new();

    // Processes the import chain until all dependencies are resolved
    while !import_chain.is_empty() {
        let current_import = import_chain.pop().unwrap();

        // Check if the current file is already in a tracked chain to detect circular dependencies
        if let Some(chain) = all_import_chains.get(&current_import) {
            if chain.contains(&current_import) {
                Logger::err(
                    "Circular dependency detected!",
                    &[
                        "You have a circular dependency in your imports!",
                        "This may also happen because you have imported a file twice or more"
                    ],
                    build_trace(chain.clone())
                        .iter()
                        .map(|s| s.as_str())
                        .collect::<Vec<&str>>().as_slice()
                );
                exit(1);
            }
        }

        // If the file does not exist, log an error and exit
        if !current_import.exists() {
            // Try to figure out where the import is coming from
            let import_trace = imports
                .iter()
                .filter(|i| i.get_from() == current_import)
                .map(|i| i.get_trace())
                .collect::<Vec<String>>();

            Logger::err(
                "Import file not found!",
                &["The import file doesn't exist"],
                import_trace
                    .iter()
                    .map(|s| s.as_str())
                    .collect::<Vec<&str>>().as_slice()
            );
            exit(1);
        }

        // Skip already processed files to avoid redundant work
        if checked_files.contains(&current_import) {
            continue;
        }

        // Read the contents of the current file, logging an error if reading fails
        let contents = try_unwrap(
            read_to_string(current_import.clone()),
            "Failed to read the import file"
        );

        // Tokenize the file's contents and extract its import dependencies
        let tokens = lexer.tokenize(
            &contents,
            &current_import.to_str().unwrap().to_string()
        );

        let parts = extract_parts(&tokens, current_import.clone());
        let chained_imports = parts.get_imports();
        let mut total_impots : Vec<Import> = vec![];

        total_impots.extend(chained_imports.clone());
        total_impots.extend(imports.clone());

        // Track the current import chain and prevent circular dependencies
        let current_chain = all_import_chains.entry(current_import.clone()).or_insert_with(Vec::new);
        current_chain.push(current_import.clone());

        // Add dependencies from the current file to the import chain for further processing
        let current_chain_clone = current_chain.clone();
        for import in chained_imports.iter().filter(check_extension) {
            let import_path = import.get_from().clone();
            import_chain.push(import_path.clone());
            let mut new_chain = current_chain_clone.clone();
            new_chain.push(import_path.clone());
            all_import_chains.insert(import_path, new_chain);
        }

        // Mark the current file as checked to avoid reprocessing
        checked_files.push(current_import);
    }
}