use fancy_regex::Regex;
use lazy_static::lazy_static;
use shared::{logger::{Logger, LoggerImpl}, result::try_unwrap};

lazy_static! {
    // Used to print warnings for cammel case variable names
    // Surf encourages snake case variable names!
    pub static ref CAMMEL_CASE_REGEX: Regex = 
        Regex::new(r"^[a-zA-Z][a-zA-Z0-9]*$").unwrap();
}

fn check_variable_name(var_name: &String, trace: &String) {
    if try_unwrap(
        CAMMEL_CASE_REGEX.is_match(var_name),
        "Failed to validate a variable name"
    ) {
        Logger::warn(
            "Consider using snake case for variable names",
            &[
                format!(
                    "Consider converting {} to snake case",
                    var_name
                ).as_str(),
                trace.as_str()
            ],
        );
    }
}

pub fn check_variable() {
    
}