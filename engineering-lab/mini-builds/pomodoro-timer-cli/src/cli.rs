use clap::{Parser, Subcommand};
use std::path::PathBuf;

#[derive(Parser, Debug)]
#[command(name = "focus-timer")]
#[command(about = "A configurable focus timer with desktop notifications")]
pub struct Cli {
    #[command(subcommand)]
    pub command: Commands,
}

#[derive(Subcommand, Debug)]
pub enum Commands {
    /// Run a timer sequence from a preset or inline block list
    Run {
        /// Preset name from the config file
        #[arg(long, conflicts_with = "blocks")]
        preset: Option<String>,

        /// Inline blocks like "Focus:25,Break:5,Focus:25,Long Break:15"
        #[arg(long, conflicts_with = "preset")]
        blocks: Option<String>,

        /// Optional path to a TOML config file
        #[arg(long)]
        config: Option<PathBuf>,

        /// Repeat the selected block sequence this many times
        #[arg(long, default_value_t = 1)]
        repeat: u32,

        /// Print countdown updates every N seconds
        #[arg(long, default_value_t = 60)]
        tick_seconds: u64,

        /// Disable desktop notifications
        #[arg(long)]
        no_notify: bool,
    },
    /// Show available presets from the config file
    List {
        /// Optional path to a TOML config file
        #[arg(long)]
        config: Option<PathBuf>,
    },
}
