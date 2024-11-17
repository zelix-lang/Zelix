use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct SurfBinds {
    pub name: String,
    pub value: String
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SurfConfigFile {
    pub name: String,
    pub description: String,
    pub version: String,
    pub author: String,
    pub license: String,
    pub git: String,
    pub main_file: String,
    pub repositories: Vec<String>,
    pub dependencies: Vec<String>,
    pub bind: Vec<SurfBinds>
}