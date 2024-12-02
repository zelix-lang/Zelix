use fancy_regex::Regex;
use lazy_static::lazy_static;
use crate::result::try_unwrap;

lazy_static! {
    pub static ref IMPORT_REGEX : Regex = try_unwrap(
        // This regex also excludes the standard library imports
        // which are then processed by the lexer
        Regex::new(r#"import\s+"(?!@Surf:standard)[\s\S]*?"\s*;"#),
        "Failed to compile regex pattern for imports"
    );
}

pub fn extract_import_matches(input: String) -> Vec<String> {
    let matches: Vec<_> = IMPORT_REGEX.find_iter(input.as_str()).collect();

    matches.iter()
        .map(|m|
            try_unwrap(
                m.clone(),
                "Failed to extract import from code"
            ).as_str().to_string()
        )
        .collect()
}