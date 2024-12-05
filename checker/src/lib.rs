mod function;
mod header;
mod body;

use function::function_checker::analyze_functions;
use function::main_function_checker::check_main_function;
use body::scope_checker::analyze_scope;

use shared::code::file_code::{FileCode, FileCodeImpl};

pub fn analyze(code: &FileCode) {

    analyze_functions(code.get_functions(), code.get_imports());
    check_main_function(code.get_functions());
    analyze_scope(code);

}