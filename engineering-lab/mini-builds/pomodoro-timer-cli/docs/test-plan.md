# Test Plan

## Scope

The first test pass covers deterministic logic that can be verified without waiting for real timers.

## Unit Tests

- Inline block parsing succeeds for valid `Label:Minutes` sequences.
- Inline block parsing rejects zero-minute or malformed blocks.
- TOML config parsing succeeds for valid presets.
- Config validation rejects empty preset names and empty block lists.
- Preset lookup returns the expected block sequence.
- Countdown formatting renders `MM:SS` correctly.

## Deferred Tests

- End-to-end CLI tests once the Rust toolchain is available in the environment.
- Notification behavior on macOS, Linux, and Windows.
- Runtime tests for the timer loop using an injectable sleeper abstraction.

## Exit Criteria

- All unit tests pass under `cargo test`.
- README commands match the actual CLI.
- Sample config passes validation and can be listed successfully.
