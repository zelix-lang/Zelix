mod function_checker;
mod main_function_checker;
mod scope_checker;
pub mod variable_checker;

use function_checker::analyze_functions;
use main_function_checker::check_main_function;
use scope_checker::analyze_scope;

use shared::code::file_code::{FileCode, FileCodeImpl};

pub fn analyze(code: &FileCode) {

    // Create a Clang instance and an Index to parse .h and .hpp files
    // in order to determine non-native functions and types
    // which are needed for static analysis
    analyze_functions(code.get_functions(), code.get_source());
    check_main_function(code.get_functions());
    analyze_scope(code);

}