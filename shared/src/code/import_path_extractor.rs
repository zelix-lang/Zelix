use std::path::PathBuf;

pub fn extract_import_path(raw: &String, start_point: &PathBuf) -> PathBuf {
    let mut import_raw = raw.replacen("@import", "", 1);
    import_raw = import_raw.trim().to_string();

    // Drop the semicolon and the quotes
    let import = import_raw[1..import_raw.len()-2].to_string();

    start_point.join(import)
}

pub fn extract_imports_paths(imports: &Vec<String>, start_point: &PathBuf) -> Vec<PathBuf> {
    imports.iter().map(|m| {
        extract_import_path(m, start_point)
    }).collect()
}