pub fn extract_import_path(raw: &String) -> String {
    let mut import_raw = raw.replacen("@import", "", 1);
    import_raw = import_raw.trim().to_string();

    // Drop the semicolon and the quotes
    let import = import_raw[1..import_raw.len()-2].to_string();

    import
}

pub fn extract_imports_paths(imports: &Vec<String>) -> Vec<String> {
    imports.iter().map(|m| {
        extract_import_path(m)
    }).collect()
}