use shared::code::{file_code::{FileCode, FileCodeImpl}, function::FunctionImpl};

use super::{body_transpiler::transpile_body, type_transpiler::transpile_type};

pub fn transpile_functions(file_code: &FileCode, transpiled_code: &mut String) {
    for (name, function) in file_code.get_functions() {

        let is_main = name == "main";

        if is_main {
            transpiled_code.push_str("int ");
        } else {
            transpile_type(function.get_return_type(), transpiled_code);
        }

        transpiled_code.push_str(name);
        transpiled_code.push_str("() {\n");

        transpile_body(function.get_body(), transpiled_code);

        if is_main {
            transpiled_code.push_str("\nreturn 0;");
        }

        transpiled_code.push_str("\n}");

    }
}