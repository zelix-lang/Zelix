// This workspace is used to parse C and C++ header files using the clang crate.
// It's not inteded to execute or analyze Surf code, but to provide a way to parse
// C and C++ header files which are needed for static analysis and code generation.

use clang::{Clang, Index, TranslationUnit};

pub fn create_c_instance() -> Clang {
    Clang::new().unwrap()
}

// Specify the lifetime of the Clang instance and the Index returned
pub fn create_index<'a>(c: &'a Clang) -> Index<'a> {
    Index::new(c, false, false)
}

pub fn parse_header_file<'a>(path: &String, index: &'a Index<'a>) -> TranslationUnit<'a> {
    index.parser(path)
            .arguments(&["-x", "c++"])
            .parse()
            // Can't use try_unwrap without creating a circular dependency
            // It's highly unlikely that this will fail, so just use expect
            .expect("Failed to parse header file")
}
