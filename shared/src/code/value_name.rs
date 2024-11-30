pub mod value_name {
    use fancy_regex::Regex;
    use lazy_static::lazy_static;


    lazy_static! {
        pub static ref VALUE_NAME_REGEX: Regex = Regex::new(r"^[a-zA-Z_][a-zA-Z0-9_]*$").unwrap();
    }

}