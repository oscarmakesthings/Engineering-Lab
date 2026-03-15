use crate::config::Block;
use crate::notify::notify;
use std::io::{self, Write};
use std::thread;
use std::time::Duration;

pub fn run_block(
    block: &Block,
    current: u32,
    total: u32,
    tick_seconds: u64,
    notifications_enabled: bool,
) -> Result<(), String> {
    if tick_seconds == 0 {
        return Err("tick-seconds must be at least 1".to_string());
    }

    let seconds_total = block.minutes.saturating_mul(60);
    println!("[{current}/{total}] {} for {} minute(s)", block.label, block.minutes);
    notify(
        notifications_enabled,
        &format!("Starting {}", block.label),
        &format!("{} minute(s)", block.minutes),
    );

    let mut remaining = seconds_total;
    while remaining > 0 {
        print!("\r{}", format_time_remaining(&block.label, remaining));
        io::stdout().flush().map_err(|error| error.to_string())?;
        let sleep_for = remaining.min(tick_seconds);
        thread::sleep(Duration::from_secs(sleep_for));
        remaining = remaining.saturating_sub(sleep_for);
    }

    println!();
    println!("{} complete.", block.label);
    notify(
        notifications_enabled,
        &format!("{} complete", block.label),
        "Move to the next block.",
    );

    Ok(())
}

pub fn format_time_remaining(label: &str, remaining_seconds: u64) -> String {
    let minutes = remaining_seconds / 60;
    let seconds = remaining_seconds % 60;
    format!(
        "███████  ██████   ██████ ██    ██ ███████\n\
         ██      ██    ██ ██      ██    ██ ██     \n\
         █████   ██    ██ ██      ██    ██ ███████\n\
         ██      ██    ██ ██      ██    ██      ██\n\
         ██       ██████   ██████  ██████  ███████\n\
                                                    \n\
                    {:02}:{:02}                    \n\
                                                    \n\
         ------------------------------------------",
        minutes, seconds
    )
}

#[cfg(test)]
mod tests {
    use super::format_time_remaining;

    #[test]
    fn formats_time_remaining() {
        let result = format_time_remaining("Focus", 125);
        assert!(result.contains("02:05"));
    }
}
