use shared::result::try_unwrap;

use crate::regex_patterns::IMPORT_REGEX;

pub fn extract_import_matches(input: &String) -> Vec<&str> {
    let matches: Vec<_> = IMPORT_REGEX.find_iter(input).collect();

    matches.iter().map(|m| 
        try_unwrap(
            m.clone(),
            "Failed to unwrap import"
        ).as_str()
    ).collect()
}