mod function_checker;
mod main_function_checker;
mod scope_checker;
mod header_checker;
mod variable_checker;

use function_checker::analyze_functions;
use main_function_checker::check_main_function;
use scope_checker::analyze_scope;

use shared::code::file_code::{FileCode, FileCodeImpl};

pub fn analyze(code: &FileCode) {

    analyze_functions(code.get_functions());
    check_main_function(code.get_functions());
    analyze_scope(code);

}