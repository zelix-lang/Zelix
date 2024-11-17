use std::process::exit;

use crate::logger::{Logger, LoggerImpl};

pub fn try_unwrap<T, E>(result: Result<T, E>, error_message: &str) -> T
where E: std::fmt::Debug, {
    if result.is_err() {
        Logger::err(
            error_message,
            &[&"No help available, sorry"],
            &[&"No trace available, sorry"]
        );
        exit(1);
    }

    return result.unwrap();
}