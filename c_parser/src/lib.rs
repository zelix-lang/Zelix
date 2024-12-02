// This workspace is used to parse C and C++ header files using the clang crate.
// It's not inteded to execute or analyze Surf code, but to provide a way to parse
// C and C++ header files which are needed for static analysis and code generation.

pub mod wrapper;
pub mod class;
pub mod function;
pub mod header;
use clang::{Clang, Index, TranslationUnit};

pub fn create_c_instance() -> Clang {
    Clang::new().unwrap()
}

// Specify the lifetime of the Clang instance and the Index returned
pub fn create_index(c: &Clang) -> Index {
    Index::new(&c, false, false)
}

pub fn parse_header_file<'a>(path: &String, index: &'a Index) -> TranslationUnit<'a> {
    index.parser(path)
            .arguments(&["-x", "c++"])
            .parse()
            // Can't use try_unwrap without creating a circular dependency
            // It's highly unlikely that this will fail, so just use expect
            .expect("Failed to parse header file")
}
