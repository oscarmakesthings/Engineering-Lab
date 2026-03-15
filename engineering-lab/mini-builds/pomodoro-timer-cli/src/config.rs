use serde::Deserialize;
use std::fs;
use std::path::Path;

#[derive(Debug, Clone, Deserialize, PartialEq, Eq)]
pub struct Config {
    pub presets: Vec<Preset>,
}

#[derive(Debug, Clone, Deserialize, PartialEq, Eq)]
pub struct Preset {
    pub name: String,
    pub blocks: Vec<Block>,
}

#[derive(Debug, Clone, Deserialize, PartialEq, Eq)]
pub struct Block {
    pub label: String,
    pub minutes: u64,
}

pub fn load_config(path: &Path) -> Result<Config, String> {
    let contents = fs::read_to_string(path)
        .map_err(|error| format!("failed to read config {}: {error}", path.display()))?;
    load_config_from_str(&contents)
        .map_err(|error| format!("failed to parse config {}: {error}", path.display()))
}

pub fn load_config_from_str(contents: &str) -> Result<Config, toml::de::Error> {
    toml::from_str(contents)
}

pub fn find_preset(config: &Config, preset_name: &str) -> Result<Vec<Block>, String> {
    config
        .presets
        .iter()
        .find(|preset| preset.name == preset_name)
        .map(|preset| preset.blocks.clone())
        .ok_or_else(|| format!("preset '{preset_name}' not found"))
}

pub fn default_preset_name(config: &Config) -> Result<&str, String> {
    config
        .presets
        .first()
        .map(|preset| preset.name.as_str())
        .ok_or_else(|| "config must contain at least one preset".to_string())
}

pub fn validate_config(config: &Config) -> Result<(), String> {
    if config.presets.is_empty() {
        return Err("config must contain at least one preset".to_string());
    }

    for preset in &config.presets {
        if preset.name.trim().is_empty() {
            return Err("preset names must not be empty".to_string());
        }

        if preset.blocks.is_empty() {
            return Err(format!("preset '{}' must contain at least one block", preset.name));
        }

        for block in &preset.blocks {
            validate_block(block)
                .map_err(|error| format!("preset '{}': {error}", preset.name))?;
        }
    }

    Ok(())
}

pub fn parse_inline_blocks(input: &str) -> Result<Vec<Block>, String> {
    let mut blocks = Vec::new();

    for raw_part in input.split(',') {
        let part = raw_part.trim();
        if part.is_empty() {
            continue;
        }

        let (label, minutes) = part
            .rsplit_once(':')
            .ok_or_else(|| format!("invalid block '{part}', expected Label:Minutes"))?;
        let block = Block {
            label: label.trim().to_string(),
            minutes: minutes
                .trim()
                .parse::<u64>()
                .map_err(|_| format!("invalid minute value in block '{part}'"))?,
        };

        validate_block(&block).map_err(|error| format!("block '{part}': {error}"))?;
        blocks.push(block);
    }

    if blocks.is_empty() {
        return Err("no blocks were provided".to_string());
    }

    Ok(blocks)
}

fn validate_block(block: &Block) -> Result<(), String> {
    if block.label.trim().is_empty() {
        return Err("labels must not be empty".to_string());
    }

    if block.minutes == 0 {
        return Err("minutes must be at least 1".to_string());
    }

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::{
        default_preset_name, find_preset, load_config_from_str, parse_inline_blocks,
        validate_config, Block,
    };

    #[test]
    fn parses_inline_blocks() {
        let blocks = parse_inline_blocks("Focus:25,Break:5,Review:15").unwrap();
        assert_eq!(
            blocks,
            vec![
                Block {
                    label: "Focus".to_string(),
                    minutes: 25
                },
                Block {
                    label: "Break".to_string(),
                    minutes: 5
                },
                Block {
                    label: "Review".to_string(),
                    minutes: 15
                }
            ]
        );
    }

    #[test]
    fn rejects_zero_minute_inline_block() {
        let error = parse_inline_blocks("Focus:0").unwrap_err();
        assert!(error.contains("minutes must be at least 1"));
    }

    #[test]
    fn loads_and_validates_config() {
        let config = load_config_from_str(
            r#"
            [[presets]]
            name = "classic"
            blocks = [
              { label = "Focus", minutes = 25 },
              { label = "Break", minutes = 5 }
            ]
            "#,
        )
        .unwrap();

        validate_config(&config).unwrap();
        let blocks = find_preset(&config, "classic").unwrap();
        assert_eq!(blocks.len(), 2);
    }

    #[test]
    fn rejects_empty_preset_names() {
        let config = load_config_from_str(
            r#"
            [[presets]]
            name = " "
            blocks = [
              { label = "Focus", minutes = 25 }
            ]
            "#,
        )
        .unwrap();

        let error = validate_config(&config).unwrap_err();
        assert!(error.contains("preset names must not be empty"));
    }

    #[test]
    fn returns_first_preset_as_default() {
        let config = load_config_from_str(
            r#"
            [[presets]]
            name = "classic-mode"
            blocks = [
              { label = "Focus", minutes = 25 }
            ]

            [[presets]]
            name = "deep-work"
            blocks = [
              { label = "Deep Work", minutes = 50 }
            ]
            "#,
        )
        .unwrap();

        assert_eq!(default_preset_name(&config).unwrap(), "classic-mode");
    }
}
