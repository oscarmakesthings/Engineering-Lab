# Focus Timer

A small Rust CLI focus timer with configurable time blocks and desktop notifications.

## Features

- Named presets loaded from `pomodoro.toml`
- Inline custom block sequences
- Optional sequence repetition
- Config validation for presets and blocks
- Desktop notifications via system tools:
  - macOS: `osascript`
  - Linux: `notify-send`
  - Windows: PowerShell toast notifications

## Project Structure

```text
src/
├── app.rs
├── cli.rs
├── config.rs
├── lib.rs
├── main.rs
├── notify.rs
└── timer.rs
```

## Usage

List configured presets:

```bash
cargo run -- list
```

Run the default preset:

```bash
cargo run -- run
```

Run a named preset:

```bash
cargo run -- run --preset deep-work
```

Run a fully custom sequence:

```bash
cargo run -- run --blocks "Focus:30,Break:10,Focus:45,Long Break:20"
```

Repeat the sequence twice and print updates every 30 seconds:

```bash
cargo run -- run --preset study-sprint --repeat 2 --tick-seconds 30
```

Disable notifications:

```bash
cargo run -- run --blocks "Code:45,Tea:10" --no-notify
```

Run from any directory after installing the binary by putting your config at:

```text
~/.config/focus-timer/pomodoro.toml
```

You can still override the path explicitly:

```bash
focus-timer run --config "/absolute/path/to/pomodoro.toml"
```

## SDLC Notes

- Requirements are documented in `docs/requirements.md`.
- The initial test plan is in `docs/test-plan.md`.
- Unit tests currently target parsing, config validation, and time formatting.
