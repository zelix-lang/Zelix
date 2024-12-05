use std::{collections::{HashMap, HashSet}, process::exit};
use c_parser::header::Header;
use logger::{Logger, LoggerImpl};

use shared::code::{function::{Function, FunctionImpl}, value_name::value_name::{CPP_KEYWORDS, VALUE_NAME_REGEX}};

use super::lifetime_checker::check_lifetime;

pub fn analyze_functions(
    // Pass by reference to avoid moving the value or cloning it 
    functions: &HashMap<String, HashMap<String, Function>>,
    imports: &Vec<Header>
) {
    // Save the functions we've seen so far so we can detect multiple definitions
    let mut seen_functions: HashSet<String> = HashSet::new();

    // Extractor checks for duplicated definitions but only in the same file
    // so we need to check for duplicated definitions across all files
    for (_, file_functions) in functions.iter() {
        for (name, function) in file_functions.iter() {
            if 
                !VALUE_NAME_REGEX.is_match(name.as_str()).unwrap_or(false)
                || CPP_KEYWORDS.contains(&name.as_str())
            {
                Logger::err(
                    format!("Invalid function name: {}", name).as_str(),
                    &[
                        "Function names must start with a letter or an underscore",
                    ],
                    &[
                        function.get_trace().as_str()
                    ]
                );
    
                exit(1);
            }
    
            if seen_functions.contains(name) {
                Logger::err(
                    format!("Function {} already defined", name).as_str(),
                    &[
                        "Functions can only be defined once",
                    ],
                    &[
                        function.get_trace().as_str()
                    ]
                );
    
                exit(1);
            }

            if function.is_public() {
                seen_functions.insert(name.clone());
            }

            check_lifetime(function, imports);
        }
    }

}