mod function_analyzer;
mod scope_analyzer;
mod main_function_checker;

use function_analyzer::analyze_functions;
use scope_analyzer::analyze_scope;

use crate::shared::file_code::{FileCode, FileCodeImpl};

pub fn analyze(code: &FileCode) {

    analyze_functions(code.clone().get_functions(), code.get_source().clone());
    analyze_scope(code.clone());

}