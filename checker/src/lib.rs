mod function_checker;
mod scope_checker;
mod main_function_checker;

use function_checker::analyze_functions;
use scope_checker::analyze_scope;

use shared::code::file_code::{FileCode, FileCodeImpl};

pub fn analyze(code: &FileCode) {

    analyze_functions(code.clone().get_functions(), code.get_source().clone());
    analyze_scope(code.clone());

}