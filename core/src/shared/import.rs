use std::path::PathBuf;

#[derive(Debug, Clone)]
pub struct Import {

    /// What to import
    name: String,

    /// Where to import from
    from: PathBuf,

    /// The trace of the import chain
    trace: String

}

pub trait Importable {

    fn new(name: String, from: PathBuf, trace: String) -> Self;
    fn get_name(&self) -> String;
    fn get_from(&self) -> PathBuf;
    fn get_trace(&self) -> String;
    
}

impl Importable for Import {

    fn new(name: String, from: PathBuf, trace: String) -> Import {
        Import {
            name,
            from,
            trace
        }
    }

    fn get_name(&self) -> String {
        self.name.clone()
    }

    fn get_from(&self) -> PathBuf {
        self.from.clone()
    }

    fn get_trace(&self) -> String {
        self.trace.clone()
    }

}