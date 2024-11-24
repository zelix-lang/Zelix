use std::fs::File;
use std::io::Write;
use std::path::PathBuf;

use shared::token::token::Token;

use crate::extractor::extract_parts;
use crate::shared::file_code::FileCode;

pub fn transpile(tokens: Vec<Token>, out_dir: PathBuf) {

    let mut transpiled_code = "";
    let file_code = extract_parts(&tokens);

    println!("{:?}", file_code);

}