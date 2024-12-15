pub mod value_name {
    use fancy_regex::Regex;
    use lazy_static::lazy_static;

    lazy_static! {
        pub static ref VALUE_NAME_REGEX: Regex = Regex::new(r"^[a-zA-Z_][a-zA-Z0-9_]*$").unwrap();
        // A collection of all C++ keywords that match the regex above
        // This is used to prevent users from using C++ keywords as function or variable names
        // which will cause a conflict with the generated code
        pub static ref CPP_KEYWORDS : Vec<&'static str> = vec![
            // Skipping "bool" and "break", they're not unknown tokens
            // Since only unknown tokens are allowed, adding them here is redundant
            // because the analyzer will catch them anyway
            "alignas", "alignof", "and", "and_eq", "asm", "auto", "bitand", "bitor",
            "case", "catch", "char", "char8_t", "char16_t", "char32_t", "class", "compl",
            // Skipped "const" and "continue" as well
            "concept", "constexpr", "const_cast", "co_await", "co_return", "co_yield", 
            "decltype", "default", "delete", "do", "double", "dynamic_cast",
            // Skipped "else", "false" and "for"
            "else", "enum", "explicit", "export", "extern", "float", "friend", 
            "goto", "inline", "int", "long", "mutable", "namespace", "new",
            // Skipped "if"
            "noexcept", "not", "not_eq", "nullptr", "operator", "or", "or_eq",
            "private", "protected", "public", "register", "reinterpret_cast", "requires",
            // Skipped "return"
            "short", "signed", "sizeof", "static", "static_assert", "static_cast", "struct",
            // Skipped "true"
            "switch", "template", "this", "thread_local", "throw", "try", "typedef",
            "typeid", "typename", "union", "unsigned", "using", "virtual", "void",
            // Skipped "while"
            "volatile", "wchar_t", "xor", "xor_eq"
        ];
    }

}