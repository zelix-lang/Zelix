use shared::code::{file_code::{FileCode, FileCodeImpl}, function::FunctionImpl, param::ParamImpl};

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

        if is_main {
            // We don't have to worry about additional parameters here
            // the static analyzer makes sure that the main function
            // doesn't have any parameters
            transpiled_code.push_str("(int arg_count, char* args[]");
        } else {
            transpiled_code.push_str("(");
        }

        for param in function.get_arguments() {
            transpile_type(param.get_data_type(), transpiled_code);

            if param.is_reference() {
                // Add a pointer to the data type
                transpiled_code.push_str("*");
            }
            
            transpiled_code.push_str(param.get_name().as_str());
        }

        transpiled_code.push_str(") {\n");

        transpile_body(function.get_body(), transpiled_code);

        if is_main {
            transpiled_code.push_str("\nreturn 0;");
        }

        transpiled_code.push_str("\n}");

    }
}