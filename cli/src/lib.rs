use args::Commands;
use clap::Parser;
use command::{
    init::init_command, license::license_command, run::run_command,
    version::version_command,
};
mod args;
mod command;
mod example;
mod structs;
mod processor;

pub fn main() {
    let parsed_args = args::Args::parse();

    match parsed_args.command {
        Commands::Init { path } => {
            init_command(path);
        }
        Commands::Run { path } => {
            run_command(path);
        }
        Commands::License { show_full } => {
            license_command(show_full);
        }
        Commands::Version {} => {
            version_command();
        }
    }
}
