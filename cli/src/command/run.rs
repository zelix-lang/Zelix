use std::path::PathBuf;

use extractor::extract_parts;

use super::lexe_base::lexe_base;

pub fn run_command(path: Option<PathBuf>) {

    let (tokens, bindings, final_path) = lexe_base(path.clone());
    // Wrap the tokens into an AST
    let ast = extract_parts(&tokens, final_path);

    // Run the interpreter
    
}