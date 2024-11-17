pub mod example_program {
    use lazy_static::lazy_static;

    lazy_static! {
        pub static ref EXAMPLE_PROGRAM: String = format!(
            "{}\n\n{}\n{}\n{}",
            "import from \"@Surf:standard/io/println\";",
            "fun main() -> nothing {",
            "    println(\"Hello, world!\");",
            "}"
        );

        pub static ref EXAMPLE_GIT_IGNORE : String = "/out".to_string();
    }
}