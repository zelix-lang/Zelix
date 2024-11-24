use super::import::Import;

pub struct ImportGroup {

    imports: Vec<Import>,
    skipped_tokens: i32

}

pub trait ImportGroupImpl {

    fn new(imports: Vec<Import>, skipped_tokens: i32) -> Self;

    fn get_imports(&self) -> &Vec<Import>;
    fn get_skipped_tokens(&self) -> i32;

}

impl ImportGroupImpl for ImportGroup {

    fn new(imports: Vec<Import>, skipped_tokens: i32) -> Self {
        ImportGroup {
            imports,
            skipped_tokens
        }
    }

    fn get_imports(&self) -> &Vec<Import> {
        &self.imports
    }

    fn get_skipped_tokens(&self) -> i32 {
        self.skipped_tokens
    }

}