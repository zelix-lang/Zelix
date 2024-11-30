use fancy_regex::Regex;
use lazy_static::lazy_static;
use shared::result::try_unwrap;

lazy_static! {
    pub static ref IMPORT_REGEX : Regex = try_unwrap(
        Regex::new(r#"@import\s\s*"[\s\S]*?"\s*;"#),
        "Failed to compile regex pattern for imports"
    );

    pub static ref NUMBER_REGEX: Regex = Regex::new(r#"^\d+((\.\d+)?)$"#)
        .unwrap();
}