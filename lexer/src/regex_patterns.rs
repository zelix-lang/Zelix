use fancy_regex::Regex;
use lazy_static::lazy_static;

lazy_static! {
    pub static ref NUMBER_REGEX: Regex = Regex::new(r#"^\d+((\.\d+)?)$"#)
        .unwrap();
}