use shared::code::{file_code::{FileCode, FileCodeImpl}, import::Importable};

pub fn transpile_imports(file_code: &FileCode, transpiled_code: &mut String) {
    // First add the imports
    let imports = file_code.get_seen_imports();

    for import in imports {
        transpiled_code.push_str("#include \"");

        transpiled_code.push_str(
            import.get_from().to_str().unwrap()
        );

        transpiled_code.push_str("\"\n");
    }

    // Add imports that are needed
    transpiled_code.push_str("#include <string>\n");
}