use std::process::exit;

use c_parser::header::Header;
use code::types::ParamTypeImpl;
use extractor::return_extractor::extract_function_returns;
use logger::{Logger, LoggerImpl};
use shared::code::function::{Function, FunctionImpl};

pub fn check_function_return(
    function: &Function,
    imports: &Vec<Header>
) {
    let return_type = function.get_return_type();
    let return_sentences = extract_function_returns(function.get_body());

    // Check for nothing-type return
    if 
        return_type.get_params().is_empty()
        && return_type.get_name() == "nothing"
    {
        // The return sentences must be inexistent or have no tokens
        // If they don't have tokens the representation is "return;", which is valid
        if return_sentences.iter().any(|sentence| !sentence.is_empty()) {
            Logger::err(
                "Function with nothing return type must not return anything",
                &[
                    "Functions with nothing return type must not return anything",
                ],
                &[
                    function.get_trace().as_str()
                ]
            );

            exit(1);
        }
    }
}