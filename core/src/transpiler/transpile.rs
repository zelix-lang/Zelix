use std::fs::File;
use std::io::Write;
use std::path::PathBuf;

use shared::token::token::{Token, TokenImpl};
use shared::token::token_type::TokenType;

use crate::extractor::extract_parts;
use crate::shared::file_code::FileCodeImpl;
use crate::shared::function::FunctionImpl;
use crate::shared::import::{Import, Importable};
use crate::syntax::analyze;

/// Transpiles Surf code into C++ code
/// This can be later compiled with G++, Clang++ or GCC

// Returns a vector of imports that are needed for later compilation
pub fn transpile(tokens: Vec<Token>, out_dir: PathBuf, source: PathBuf) -> Vec<Import> {

    let mut transpiled_code = String::new();
    let file_code = extract_parts(&tokens, source);

    // Borrow to avoid moving the value or cloning it
    analyze(&file_code);

    // First add the imports
    let imports = file_code.get_imports();

    for import in imports {
        transpiled_code.push_str("#include \"");

        transpiled_code.push_str(
            import.get_from().to_str().unwrap()
        );

        transpiled_code.push_str("\"\n");
    }

    // Add imports that are needed
    transpiled_code.push_str("#include <string>\n");

    // Use the namespace std
    transpiled_code.push_str("using namespace std;\n");

    // Transpile the functions
    for (name, function) in file_code.get_functions() {

        let is_main = name == "main";

        if is_main {
            transpiled_code.push_str("int ");
        } else {
            if function.get_return_type() == &TokenType::Nothing {
                transpiled_code.push_str("void ");
            } else if function.get_return_type() == &TokenType::Num {
                transpiled_code.push_str("double ");
            } else if function.get_return_type() == &TokenType::String {
                transpiled_code.push_str("string ");
            } else if function.get_return_type() == &TokenType::Bool {
                transpiled_code.push_str("bool ");
            }
        }

        transpiled_code.push_str(name);
        transpiled_code.push_str("() {\n");

        for token in function.get_body() {
            let is_string = token.get_token_type() == TokenType::StringLiteral;

            if is_string {
                transpiled_code.push_str("\"");
            }

            transpiled_code.push_str(&token.get_value());

            if is_string {
                transpiled_code.push_str("\"");
            }

            transpiled_code.push_str(" ");
        }

        if is_main {
            transpiled_code.push_str("\nreturn 0;");
        }

        transpiled_code.push_str("\n}");

    }

    // Save the transpiled code to a file in the output directory
    let mut file = File::create(out_dir.join("out.cpp")).unwrap();
    file.write_all(transpiled_code.as_bytes()).unwrap();

    imports.clone()

}