use std::fs::File;
use std::io::Write;
use std::path::PathBuf;

use shared::token::token::Token;

use crate::extractor::extract_parts;
use crate::shared::file_code::{FileCode, FileCodeImpl};
use crate::syntax::analyze;

pub fn transpile(tokens: Vec<Token>, out_dir: PathBuf) {

    let mut transpiled_code = "";
    let file_code = extract_parts(&tokens);

    // Borrow to avoid moving the value or cloning it
    analyze(&file_code);

    println!("{:?}", file_code.get_functions());

}