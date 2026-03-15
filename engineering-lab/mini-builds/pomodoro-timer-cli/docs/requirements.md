# Requirements

## Product Goal

Provide a local CLI Pomodoro timer that supports named presets and fully customizable time blocks without requiring network access or a GUI application.

## Functional Requirements

- The user can list timer presets from a TOML configuration file.
- The user can run a timer sequence from a preset name.
- The user can run a timer sequence from inline block definitions.
- The user can repeat a selected sequence one or more times.
- The user can control countdown update frequency with `--tick-seconds`.
- The user can disable desktop notifications with `--no-notify`.
- The application validates preset and block configuration before running.

## Non-Functional Requirements

- The application must run entirely offline.
- Invalid CLI input and invalid config files must produce clear errors.
- The default workflow should rely on a single config file, `pomodoro.toml`.
- The CLI should remain cross-platform at the notification layer by using system tooling where available.

## Out of Scope

- Persistent session history
- Pause and resume
- Sound effects
- TUI or GUI mode
