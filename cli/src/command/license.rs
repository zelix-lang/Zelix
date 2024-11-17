use shared::logger::{Logger, LoggerImpl};

pub fn license_command(show_full: Option<bool>) {

    if !show_full.is_none() && show_full.unwrap() {
        Logger::log(&[
            "You may find the complete license notice in",
            "  <yellow_bright>|</yellow_bright> <black_bright>the LICENSE file of the GitHub repository</black_bright>",
            "  <blue_bright>|</blue_bright> <black_bright>the official GNU website:</black_bright>",
            "    <blue_bright>https://www.gnu.org/licenses/gpl-3.0.en.html</blue_bright>"
        ]);

        return;
    }

    Logger::log(&[
        "<blue_bright>Surf</blue_bright> Copyright (C) 2024 Rodrigo R. & all Surf contributors",
        "This program comes with ABSOLUTELY NO WARRANTY; for details type `surf license'.",
        "This is free software, and you are welcome to redistribute it",
        "under certain conditions; type `surf license --s true' for details."
    ]);
    
}