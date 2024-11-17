use fancy_regex::Regex;
use logger::add_prefix_to_text;

use crate::ansi::colors::{color_map, RESET};

mod logger {
    use lazy_static::lazy_static;
    use fancy_regex::Regex;

    lazy_static! {
        pub static ref CHAT_REGEX : Regex = Regex::new(r"<(\w+(:[#\w:]+)?)>(.*?)<\/\1>").unwrap();
    }

    pub fn add_prefix_to_text(prefix: Option<&str>, text: &[&str]) -> String {
        let mut result = String::new();
        let final_prefix = prefix.unwrap_or("<magenta_bright>   [info] | </magenta_bright>");

        for line in text {
            result.push_str(&format!("<black_bright>{}{}</black_bright>\n", final_prefix, line));
        }

        result.trim().to_string()
    }
}

pub fn colorize(message: &str) -> String {
    let re = Regex::new(r"<(\w+(:[#\w:]+)?)>(.*?)</\1>").unwrap();
    let mut result = String::from(message);

    while let Ok(Some(caps)) = re.captures(&result) {
        if let Some(color) = caps.get(1).map(|m| m.as_str()) {
            if let Some(text) = caps.get(3).map(|m| m.as_str()) {
                let reset_clone = RESET.clone();
                let color_code = color_map.get(color).unwrap_or(&reset_clone);
                let replacement = format!("{}{}{}", color_code, text, reset_clone);
                result = re.replace(&result, &replacement).to_string();
            }
        }
    }

    result
}

pub struct  Logger {}

pub trait LoggerImpl {
    fn log(message: &[&str]);
    fn warn(message: &str, details: &[&str]);
    fn err(why: &str, help: &[&str], details: &[&str]);
}

impl LoggerImpl for Logger {
    fn log(message: &[&str]) {
        for element in message {
            println!("{}", colorize(element));
        }
    }

    fn warn(message: &str, details: &[&str]) {
        Logger::log(&[
            format!(
                "<yellow_bright>[warning] | </yellow_bright>{}",
                message
            ).as_str(),
            add_prefix_to_text(
                None,
                details
            ).as_str(),
            // A new line is always convenient to separate the messages
            ""
        ]);
    }

    fn err(why: &str, help: &[&str], details: &[&str]) {
        Logger::log(&[
            format!(
                "<red_bright>[error] | </red_bright>{}",
                why
            ).as_str(),
            add_prefix_to_text(
                None,
                details
            ).as_str(),
            add_prefix_to_text(
                Some("<blue_bright>   [help] | </blue_bright>"),
                help
            ).as_str(),
            // A new line is always convenient to separate the messages
            ""
        ]);
    }
}