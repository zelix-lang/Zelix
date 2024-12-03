use fancy_regex::Regex;
use inquire::Text;

use logger::{Logger, LoggerImpl};

pub fn question(prompt: &str, default_value: &str, rule_regex: Regex) -> String {
    loop {
        match Text::new(prompt).prompt() {
            Ok(value) => {
                if rule_regex.is_match(&value).unwrap() {
                    return value;
                }
                
                Logger::log(&[&"<red_bright>Invalid input! Please try again.</red_bright>"]);
            },
            Err(_) => {
                return default_value.to_string();
            }
        }
    }
}