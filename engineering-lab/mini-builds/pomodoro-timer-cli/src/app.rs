use crate::config::{
    default_preset_name, find_preset, load_config, parse_inline_blocks, validate_config,
};
use crate::notify::notify;
use crate::timer::run_block;
use std::env;
use std::path::{Path, PathBuf};

pub fn run_command(
    preset_name: Option<String>,
    inline_blocks: Option<String>,
    config_path: Option<&Path>,
    repeat: u32,
    tick_seconds: u64,
    notifications_enabled: bool,
) -> Result<(), String> {
    if repeat == 0 {
        return Err("repeat must be at least 1".to_string());
    }

    let blocks = if let Some(blocks) = inline_blocks {
        parse_inline_blocks(&blocks)?
    } else {
        let config_path = resolve_config_path(config_path)?;
        let config = load_config(&config_path)?;
        validate_config(&config)?;
        let preset_name = match preset_name {
            Some(name) => name,
            None => default_preset_name(&config)?.to_string(),
        };
        find_preset(&config, &preset_name)?
    };

    let total_blocks = blocks.len() as u32 * repeat;
    println!("Starting timer with {total_blocks} block(s).");

    for round in 0..repeat {
        if repeat > 1 {
            println!("Round {}/{}", round + 1, repeat);
        }

        for (index, block) in blocks.iter().enumerate() {
            let current = round * blocks.len() as u32 + index as u32 + 1;
            run_block(block, current, total_blocks, tick_seconds, notifications_enabled)?;
        }
    }

    println!("Sequence complete.");
    notify(
        notifications_enabled,
        "Pomodoro complete",
        "Your configured timer sequence has finished.",
    );

    Ok(())
}

pub fn list_command(config_path: Option<&Path>) -> Result<(), String> {
    let config_path = resolve_config_path(config_path)?;
    let config = load_config(&config_path)?;
    validate_config(&config)?;

    for preset in config.presets {
        let description = preset
            .blocks
            .iter()
            .map(|block| format!("{}:{}m", block.label, block.minutes))
            .collect::<Vec<_>>()
            .join(", ");
        println!("{} -> {}", preset.name, description);
    }

    Ok(())
}

fn resolve_config_path(explicit_path: Option<&Path>) -> Result<PathBuf, String> {
    if let Some(path) = explicit_path {
        return Ok(path.to_path_buf());
    }

    let cwd_config = PathBuf::from("pomodoro.toml");
    if cwd_config.exists() {
        return Ok(cwd_config);
    }

    let home = env::var_os("HOME")
        .ok_or_else(|| "could not determine home directory for config lookup".to_string())?;
    let config_path = PathBuf::from(home)
        .join(".config")
        .join("focus-timer")
        .join("pomodoro.toml");

    if config_path.exists() {
        return Ok(config_path);
    }

    Err(
        "could not find pomodoro.toml. Pass --config or place it in the current directory or ~/.config/focus-timer/pomodoro.toml"
            .to_string(),
    )
}
