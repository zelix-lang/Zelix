use code::types::ParamTypeImpl;
use shared::code::{file_code::{FileCode, FileCodeImpl}, function::{Function, FunctionImpl}, param::ParamImpl};

use super::{body::transpile_body, type_transpiler::transpile_type};

fn transpile_arguments(function: &Function, transpiled_code: &mut String) {
    for (name, argument) in function.get_arguments() {
        transpile_type(argument.get_data_type(), transpiled_code);
        transpiled_code.push_str(name);
        transpiled_code.push_str(", ");
    }

    // Remove the last comma and space
    // No arguments = no last comma and space
    // still have to check for the comma and space
    if transpiled_code.ends_with(", ") {
        transpiled_code.truncate(transpiled_code.len() - 2);
    }
}

pub fn transpile_functions(file_code: &FileCode, transpiled_code: &mut String) {
    let imports = file_code.get_imports();
    let functions = file_code.get_functions();

    for (_, functions) in functions {
        let public_functions:Vec<&String> = functions.keys()
            .filter(|key| functions.get(*key).unwrap().is_public())
            .collect();

        let private_functions : Vec<&String> = functions.keys()
            .filter(|key| !public_functions.contains(key))
            .collect(); 

        // Transpile the public functions and add the private functions
        // as lambda functions inside the public functions' body
        // this way, the private functions are only visible inside the public functions
        // and we avoid name collisions

        for function_name in public_functions {
            let function = functions.get(function_name).unwrap();
            let is_main = function_name == "main";

            if is_main {
                transpiled_code.push_str("int ");
            } else {
                transpile_type(function.get_return_type().get_raw_tokens(), transpiled_code);
            }

            transpiled_code.push_str(function_name);
            transpiled_code.push_str("(");
            
            if is_main {
                transpiled_code.push_str("int arg_count, char* args[]");
            } else {
                transpile_arguments(function, transpiled_code);
            }

            transpiled_code.push_str(") {\n");

            // Here, we have to transpile the private functions too
            for private_function_name in &private_functions {
                let private_function = functions.get(*private_function_name).unwrap();

                transpiled_code.push_str("auto ");
                transpiled_code.push_str(private_function_name);

                // Make the lambda isolated from the outer scope
                // by adding "[]" instead of "[&]"
                transpiled_code.push_str(" = [](");
                transpile_arguments(private_function, transpiled_code);
                transpiled_code.push_str(") {\n");
                transpile_body(
                    private_function.get_body(),
                    transpiled_code,
                    imports
                );
                transpiled_code.push_str("};\n");
            }

            transpile_body(
                function.get_body(),
                transpiled_code,
                imports
            );

            if is_main {
                transpiled_code.push_str("return 0;\n");
            }

            transpiled_code.push_str("}\n\n");
        }
    }
}