use crate::logger::{Logger, LoggerImpl};

pub fn print_header() {
    Logger::log(&[&"<blue_bright>Surf Language</blue_bright>"]);
}