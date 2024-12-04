use std::collections::HashMap;
use std::fs::File;
use std::io::Write;
use std::path::PathBuf;

use code::token::Token;
use extractor::extract_parts;
use shared::code::file_code::FileCodeImpl;
use shared::code::import::Import;
use checker::analyze;

use super::util::function_transpiler::transpile_functions;
use super::util::import_transpiler::transpile_imports;

/// Transpiles Surf code into C++ code
/// This can be later compiled with G++, Clang++ or GCC

// Returns a vector of imports that are needed for later compilation
pub fn transpile(
    tokens: Vec<Token>,
    out_dir: PathBuf,
    source: PathBuf,
    bindings: HashMap<String, String>
) -> Vec<Import> {

    let mut transpiled_code = String::new();
    let file_code = extract_parts(&tokens, source);

    // Borrow to avoid moving the value or cloning it
    analyze(&file_code);

    transpile_imports(&file_code, &mut transpiled_code);

    // Transpile the functions
    transpile_functions(&file_code, &mut transpiled_code);

    // Save the transpiled code to a file in the output directory
    let mut file = File::create(out_dir.join("out.cpp")).unwrap();
    file.write_all(transpiled_code.as_bytes()).unwrap();

    file_code.get_seen_imports().clone()

}