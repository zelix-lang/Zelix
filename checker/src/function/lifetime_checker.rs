use std::process::exit;

use c_parser::header::{Header, HeaderImpl};
use code::{token::{Token, TokenImpl}, types::ParamTypeImpl};
use extractor::return_extractor::extract_function_returns;
use fancy_regex::Regex;
use lazy_static::lazy_static;
use logger::{Logger, LoggerImpl};
use shared::code::{function::{Function, FunctionImpl}, param::ParamImpl};
use util::result::try_unwrap;

lazy_static! {
    pub static ref HEAP_INVOCATION_REGEX : Regex = try_unwrap(
        Regex::new(r#"^Heap\([\s\S]*?\)\.unwrap\(\)"#),
        "Failed to compile the heap invocation regex"
    );
}

pub fn check_lifetime(function: &Function, imports: &Vec<Header>) {
    if !function.get_return_type().is_reference() {
        return;
    }

    // We do allow returning references for heap memory
    // but before we allow it, we have to figure out if Heap is imported
    // so we check we're referencing to the Heap class (part of the standard library)
    let imports_heap = imports.iter().any(|a| {
        a.get_classes().iter().any(|(name, _)| *name == "Heap".to_string())
    });

    let function_body = function.get_body();
    // Also find the return statements
    let returned_values : Vec<Vec<Token>> = extract_function_returns(function_body);

    let parameters = function.get_arguments();

    for returned_value in returned_values.iter() {
        // Check the length of the returned value
        let len = returned_value.len();

        if len == 1 {
            // The only possible way this doesn't cause a runtime error
            // is it's returning a parameter that's a reference
            let token = &returned_value[0];

            if parameters.iter().any(|param| {
                *param.get_name() == token.get_value()
                    && param.is_reference()
            }) {
                continue;
            }
        } else {
            // We'll have to check if this a heap-allocated value
            // To avoid unnecessary complexity, we can safely use regex here
            // To match any possible heap invocation and see if the match is the entire returned value
            let returned_value_str = returned_value.iter()
                .map(|token| token.get_value()).collect::<String>();

            let matches = HEAP_INVOCATION_REGEX.find_iter(&returned_value_str);
            let matches_vec = matches.collect::<Vec<_>>();

            // Only 1 match allowed
            if 
                matches_vec.len() == 1 
                && imports_heap
                && returned_value_str == try_unwrap(
                    matches_vec[0].clone(),
                    "Failed to get the matched string"
                ).as_str()
            {
                continue;
            }
        }

        // Outlived value!
        Logger::err(
            "Invalid return",
            &[
                "Use Heap<T>.unwrap() to return heap-allocated values",
                "Remove the ampersand from the return type (will return a copy)",
            ],
            &[
                "You are returning a value that outlives this function's scope",
                function.get_trace().as_str()
            ]
        );

        exit(1);
    }
}