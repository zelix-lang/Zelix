mod function_analyzer;
mod import_analyzer;

use function_analyzer::analyze_functions;
use import_analyzer::analyze_imports;

use crate::shared::file_code::{FileCode, FileCodeImpl};

pub fn analyze(code: &FileCode) {

    analyze_functions(code.clone().get_functions());
    analyze_imports(code.clone());

}