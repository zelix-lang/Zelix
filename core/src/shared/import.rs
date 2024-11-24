use std::path::PathBuf;

#[derive(Debug, Clone)]
pub struct Import {

    /// What to import
    name: String,

    /// Where to import from
    from: PathBuf

}

pub trait Importable {

    fn new(name: String, from: PathBuf) -> Self;
    fn get_name(&self) -> String;
    fn get_from(&self) -> PathBuf;
    
}

impl Importable for Import {

    fn new(name: String, from: PathBuf) -> Self {
        Import {
            name,
            from
        }
    }

    fn get_name(&self) -> String {
        self.name.clone()
    }

    fn get_from(&self) -> PathBuf {
        self.from.clone()
    }

}