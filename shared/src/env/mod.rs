use std::{env::var, path::PathBuf};

use lazy_static::lazy_static;

use util::result::try_unwrap;

lazy_static! {
    pub static ref STANDARD_LIBRARY_LOCATION : PathBuf = 
        PathBuf::from(
            try_unwrap(
                var("SURF_STANDARD_PATH"),
                "Failed to locate standard libraries"
            )
        );
}