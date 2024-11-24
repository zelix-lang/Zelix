mod function_analyzer;

use function_analyzer::analyze_functions;

use crate::shared::file_code::{FileCode, FileCodeImpl};

pub fn analyze(code: &FileCode) {

    analyze_functions(code.clone().get_functions());

}