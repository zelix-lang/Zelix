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

    /// Compiles the project
    Compile {
        /// The path to the project
        #[arg(short, long)]
        path: Option<PathBuf>
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

