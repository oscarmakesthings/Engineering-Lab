# Hangman in Go

Small terminal Hangman game written in Go.

The game uses ANSI escape codes for colored terminal output.

## Run

```bash
go run .
```

## Build

```bash
go build
```

## How it works

- Picks a random word from a built-in word list
- Accepts one-letter guesses from the terminal
- Tracks correct and incorrect guesses
- Uses colored output for the board, guesses, and win/loss messages
- Ends when the player guesses the word or reaches 6 wrong guesses
