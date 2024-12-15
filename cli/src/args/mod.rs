use std::path::PathBuf;
use clap::{Parser, Subcommand};

#[derive(Subcommand)]
pub enum Commands {
    /// Initialize a new project
    Init {
        /// The path to the project
        #[arg(short, long)]
        path: Option<PathBuf>
    },

    /// Shows the license notice
    License {
        #[arg(short, long)]
        show_full: Option<bool>
    },

    /// Runs the project
    Run {
        /// The path to the project
        #[arg(short, long)]
        path: Option<PathBuf>
    },

    Version {}
}

#[derive(Parser)]
pub struct Args {

    #[command(subcommand)]
    pub command: Commands

}

