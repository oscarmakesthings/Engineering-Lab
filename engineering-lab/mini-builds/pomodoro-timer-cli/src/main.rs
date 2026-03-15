use clap::Parser;
use focus_timer::cli::{Cli, Commands};

fn main() {
    let cli = Cli::parse();

    let result = match cli.command {
        Commands::Run {
            preset,
            blocks,
            config,
            repeat,
            tick_seconds,
            no_notify,
        } => focus_timer::app::run_command(
            preset,
            blocks,
            config.as_deref(),
            repeat,
            tick_seconds,
            !no_notify,
        ),
        Commands::List { config } => focus_timer::app::list_command(config.as_deref()),
    };

    if let Err(error) = result {
        eprintln!("error: {error}");
        std::process::exit(1);
    }
}
