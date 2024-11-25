use std::fs::File;
use std::io::Write;
use std::path::PathBuf;

use shared::token::token::Token;

use crate::extractor::extract_parts;
use crate::syntax::analyze;

pub fn transpile(tokens: Vec<Token>, out_dir: PathBuf, source: PathBuf) {

    let mut transpiled_code = "";
    let file_code = extract_parts(&tokens, source);

    // Borrow to avoid moving the value or cloning it
    analyze(&file_code);

    //println!("{:?}", file_code.get_functions());

}