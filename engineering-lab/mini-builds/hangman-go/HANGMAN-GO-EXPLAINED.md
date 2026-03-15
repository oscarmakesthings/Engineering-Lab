# Hangman Go Walkthrough

This document explains the Hangman project in:

- `go.mod`
- `main.go`

The goal is not just to tell you what each line does, but to teach the Go ideas behind it and the general shape of a terminal application.

## What This Program Is

This is a small command-line application, often shortened to CLI app.

A CLI app:

- runs in the terminal
- reads text input from the user
- prints text output back to the terminal
- usually follows a loop of input -> update state -> render output

This Hangman game follows exactly that pattern:

1. Start the game state.
2. Print the current board.
3. Ask the player for a guess.
4. Update the game state.
5. Repeat until the player wins or loses.

That loop is one of the most important patterns in terminal app development.

## File 1: `go.mod`

### Source

```go
module hangman-go

go 1.22
```

### Line-by-line explanation

`module hangman-go`

- A Go module is the project unit in modern Go.
- It tells Go the name of this module.
- In larger projects, this is often a full path like `github.com/username/projectname`.
- In this small local project, `hangman-go` is enough.

`go 1.22`

- This says the project is intended for Go 1.22 behavior.
- It helps the Go toolchain know which language and module features to expect.

### CLI app lesson

Most Go command-line programs start with a `go.mod` file because it gives the project a clean module boundary. Once this file exists, you can usually run the app with:

```bash
go run .
```

or build it with:

```bash
go build
```

## File 2: `main.go`

## Lines 1-12: Package and Imports

```go
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"
)
```

### Line 1

`package main`

- Every Go file belongs to a package.
- `main` is special.
- If a package is named `main` and contains a `main()` function, Go can build it as an executable program.
- That is why CLI apps usually use `package main`.

### Lines 3-12

The `import` block brings in standard library packages.

`"bufio"`

- Used for buffered input.
- Here it helps us read full lines from the terminal.

`"fmt"`

- Go's formatting package.
- Used for printing text like `fmt.Println` and formatted strings like `fmt.Sprintf`.

`"math/rand"`

- Used to choose a random word for the game.

`"os"`

- Gives access to operating system features.
- Here it is used for standard input/output behavior and exiting on input errors.

`"sort"`

- Used to sort guessed letters alphabetically.

`"strings"`

- Used for string helpers like lowercasing, trimming whitespace, joining pieces together, and checking whether a letter exists in the word.

`"time"`

- Used to seed randomness with the current time.

`"unicode"`

- Used to work safely with letters as Unicode characters.
- This matters because Go strings are byte-based, but characters may be more than one byte.

### Go lesson

Go prefers small, focused packages. Instead of one giant utility library, you import exactly what you need.

### CLI lesson

Terminal apps often use:

- `fmt` for output
- `bufio` or `os` for input
- `strings` for cleaning user input

That combination is extremely common.

## Lines 14-24: Constants

```go
const maxWrongGuesses = 6

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)
```

### Line 14

`const maxWrongGuesses = 6`

- `const` defines a constant value.
- Constants do not change while the program runs.
- Hangman traditionally allows a fixed number of wrong guesses, so a constant is appropriate.

### Lines 16-24

This is a grouped `const` block for ANSI terminal color codes.

Examples:

- `"\033[31m"` means switch text color to red
- `"\033[32m"` means green
- `"\033[0m"` resets formatting back to normal

### Go lesson

Use `const` for values that truly should not change. This makes your code safer and easier to read.

### CLI lesson

Many terminal apps use ANSI escape codes for color. They are just strings printed to the terminal. The terminal interprets them as formatting instructions instead of visible text.

## Lines 26-37: Word List

```go
var words = []string{
	"gopher",
	"compiler",
	"terminal",
	"variable",
	"function",
	"package",
	"pointer",
	"interface",
	"hangman",
	"concurrency",
}
```

### Explanation

`var words`

- `var` declares a variable.
- Unlike a constant, this value could change.

`[]string`

- This means "slice of strings".
- A slice is one of the most common Go data structures.
- It is a dynamic view over a sequence of values.

### Go lesson: array vs slice

In Go:

- an array has a fixed size, like `[10]string`
- a slice has flexible length, like `[]string`

Slices are used much more often than arrays in everyday Go.

### Why a slice here?

The game needs a collection of possible words. A slice is the natural choice because:

- it stores many values of the same type
- it works well with indexing
- it can grow later if needed

## Lines 39-43: The `game` Struct

```go
type game struct {
	word         string
	guessed      map[rune]bool
	wrongGuesses int
}
```

### Explanation

`type game struct { ... }`

- A `struct` groups related data together.
- This is similar to a class's fields in other languages, but without inheritance.

This struct stores the entire state of one Hangman game.

`word string`

- The secret word the player is trying to guess.

`guessed map[rune]bool`

- A `map` stores key-value pairs.
- The key type is `rune`, which represents a character.
- The value type is `bool`, meaning `true` or `false`.
- So this map answers a question like:
  "Has the player guessed the letter `a`?"

`wrongGuesses int`

- Tracks how many incorrect guesses the player has made.

### Go lesson: why `rune`?

Go strings are sequences of bytes. A `rune` is Go's type for a Unicode code point, which is closer to what people think of as a character. For letter-based logic, `rune` is usually safer than raw bytes.

### CLI app lesson

A CLI game should keep its state in one place. That makes the update loop easier to reason about. The `game` struct is that state container.

## Lines 45-51: Constructor-Like Setup

```go
func newGame(word string) *game {
	return &game{
		word:         strings.ToLower(word),
		guessed:      make(map[rune]bool),
		wrongGuesses: 0,
	}
}
```

### Explanation

`func newGame(word string) *game`

- Defines a function named `newGame`.
- It takes one parameter: `word string`
- It returns `*game`, which means "pointer to a game"

### Why a pointer?

Pointers let functions and methods work with the same underlying game object instead of copying it. Since the game state changes over time, returning a pointer is appropriate.

`return &game{ ... }`

- Creates a new `game` struct value
- `&` returns its address, giving a pointer

`strings.ToLower(word)`

- Normalizes the secret word to lowercase.
- This makes guessing logic simpler because guesses are also converted to lowercase later.

`make(map[rune]bool)`

- `make` is used to initialize certain built-in Go data structures, including maps, slices, and channels.
- A map must be initialized before use.

`wrongGuesses: 0`

- Starts the player with zero mistakes.

### Go lesson

Go does not require constructors, but functions like `newGame()` are a common pattern when object setup should be controlled or normalized.

## Lines 53-70: Displaying the Hidden Word

```go
func (g *game) displayWord() string {
	var builder strings.Builder

	for _, letter := range g.word {
		if g.guessed[letter] {
			builder.WriteString(colorGreen)
			builder.WriteRune(letter)
			builder.WriteString(colorReset)
		} else {
			builder.WriteString(colorYellow)
			builder.WriteRune('_')
			builder.WriteString(colorReset)
		}
		builder.WriteRune(' ')
	}

	return strings.TrimSpace(builder.String())
}
```

### Explanation

`func (g *game) displayWord() string`

- This is a method, not a plain function.
- The receiver is `(g *game)`, meaning this method belongs to `game`.
- It returns a `string`.

### Receiver concept

In Go, methods are functions with a receiver.

Example shape:

```go
func (x SomeType) methodName() {}
```

That is how Go attaches behavior to data types.

`var builder strings.Builder`

- `strings.Builder` is an efficient way to build strings piece by piece.
- This is better than repeated string concatenation inside a loop.

`for _, letter := range g.word`

- `range` loops over the string.
- The first value is the index, but `_` discards it.
- The second value is the current rune.

`if g.guessed[letter]`

- Looks up the letter in the guessed map.
- If it is `true`, the player already guessed it.

`builder.WriteRune(letter)`

- Appends the actual letter.

`builder.WriteRune('_')`

- Appends an underscore for an unguessed letter.

`builder.WriteRune(' ')`

- Adds a space between letters for readability.

`strings.TrimSpace(builder.String())`

- Converts the builder into a string
- removes the extra trailing space at the end

### Go lesson

This function shows three useful Go ideas:

- methods with receivers
- iterating with `range`
- building strings efficiently

### CLI lesson

Rendering is a core terminal-app job. In this project, `displayWord()` is part of the render layer. It turns game state into something the user can see.

## Lines 72-91: Processing a Guess

```go
func (g *game) guess(letter rune) (bool, string) {
	letter = unicode.ToLower(letter)

	if !unicode.IsLetter(letter) {
		return false, "Please enter a letter."
	}

	if g.guessed[letter] {
		return false, fmt.Sprintf("You already guessed %q.", letter)
	}

	g.guessed[letter] = true

	if strings.ContainsRune(g.word, letter) {
		return true, fmt.Sprintf("Nice. %q is in the word.", letter)
	}

	g.wrongGuesses++
	return false, fmt.Sprintf("Nope. %q is not in the word.", letter)
}
```

### Explanation

This method handles the rules for one player guess.

`(bool, string)`

- This method returns two values.
- The `bool` tells us whether the guess was correct.
- The `string` is the user-facing message to print.

Multiple return values are a very common Go pattern.

`letter = unicode.ToLower(letter)`

- Normalizes input to lowercase.
- This keeps logic simple and consistent.

`if !unicode.IsLetter(letter)`

- Rejects non-letter input.
- The `!` means "not".

`if g.guessed[letter]`

- Prevents the same guess from being used again.

`g.guessed[letter] = true`

- Records the guessed letter in the map.

`strings.ContainsRune(g.word, letter)`

- Checks whether the secret word contains the guessed letter.

`g.wrongGuesses++`

- Increments the wrong guess counter if the letter is not in the word.

### Go lesson

Go often favors direct, readable control flow:

- validate input
- reject bad cases early
- update state
- return a result

That style is visible here.

### CLI lesson

In a CLI program, it is helpful to separate input reading from game rules. `readGuess()` reads raw input; `guess()` applies business logic. That separation keeps the code cleaner.

## Lines 93-104: Win/Loss Checks

```go
func (g *game) won() bool {
	for _, letter := range g.word {
		if !g.guessed[letter] {
			return false
		}
	}
	return true
}

func (g *game) lost() bool {
	return g.wrongGuesses >= maxWrongGuesses
}
```

### `won()`

- Loops through every letter in the secret word.
- If even one letter has not been guessed, the player has not won.
- If the loop finishes, every letter was guessed, so return `true`.

### `lost()`

- Returns whether the number of wrong guesses has reached the allowed maximum.

### Go lesson

Small single-purpose methods are a good habit. Instead of burying this logic inside `main()`, it is named and isolated.

### CLI lesson

Terminal apps benefit from explicit state checks. `won()` and `lost()` clearly represent two end conditions for the loop.

## Lines 106-118: Showing Guessed Letters

```go
func (g *game) guessedLetters() string {
	letters := make([]string, 0, len(g.guessed))
	for letter := range g.guessed {
		letters = append(letters, string(letter))
	}

	if len(letters) == 0 {
		return "(none)"
	}

	sort.Strings(letters)
	return strings.Join(letters, ", ")
}
```

### Explanation

`letters := make([]string, 0, len(g.guessed))`

- Creates a slice of strings.
- Length starts at `0`.
- Capacity starts at `len(g.guessed)`.
- This is a small optimization because we already know roughly how many elements we might add.

`for letter := range g.guessed`

- Iterates over the keys in the map.

`append(letters, string(letter))`

- Adds the guessed rune as a string.

`if len(letters) == 0`

- If nothing has been guessed yet, return a friendly placeholder.

`sort.Strings(letters)`

- Sorts the guessed letters alphabetically.

`strings.Join(letters, ", ")`

- Turns the slice into one display string.

### Go lesson

This function teaches:

- how to create slices with `make`
- how to loop over map keys
- how to append to slices
- how to use library helpers like `sort.Strings` and `strings.Join`

## Lines 120-195: ASCII Art Renderer

```go
func hangmanArt(wrongGuesses int) string {
	stages := []string{
		`...`,
	}

	if wrongGuesses < 0 {
		wrongGuesses = 0
	}
	if wrongGuesses >= len(stages) {
		wrongGuesses = len(stages) - 1
	}

	return stages[wrongGuesses]
}
```

### Explanation

The real code contains all seven ASCII art stages.

`stages := []string{ ... }`

- Creates a slice of strings.
- Each string is one visual stage of the Hangman drawing.

The multi-line strings use backticks:

```go
`line 1
line 2`
```

These are raw string literals in Go.

### Why raw strings?

They are perfect for ASCII art because:

- newlines are preserved
- backslashes are easier to write
- the shape is more readable in the source code

### Bounds checks

`if wrongGuesses < 0`

- Prevents invalid negative indexes.

`if wrongGuesses >= len(stages)`

- Prevents indexes that are too large.

`return stages[wrongGuesses]`

- Returns the drawing matching the current game state.

### CLI lesson

This is a rendering helper. CLI apps often need functions that convert internal state into a display string. That is exactly what this one does.

## Lines 197-210: Color Helpers

```go
func colorize(text, color string) string {
	return color + text + colorReset
}

func formatRemainingGuesses(remaining int) string {
	switch {
	case remaining >= 4:
		return colorize(fmt.Sprintf("%d", remaining), colorGreen)
	case remaining >= 2:
		return colorize(fmt.Sprintf("%d", remaining), colorYellow)
	default:
		return colorize(fmt.Sprintf("%d", remaining), colorRed)
	}
}
```

### `colorize`

- Takes text and a color code.
- Wraps the text in ANSI start and reset sequences.

### `formatRemainingGuesses`

- Uses a `switch` statement with conditions.
- Green for safer values
- Yellow for medium danger
- Red for danger

### Go lesson

Go's `switch` can be used without a value. In that form, each `case` is just a boolean condition. That makes it useful as a clean alternative to `if/else if/else`.

## Lines 212-216: Random Word Selection

```go
func randomWord() string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	return words[rng.Intn(len(words))]
}
```

### Explanation

`time.Now().UnixNano()`

- Gets the current time in nanoseconds.
- This is used as the random seed.

`rand.NewSource(...)`

- Creates a random source.

`rand.New(source)`

- Builds a random number generator from that source.

`rng.Intn(len(words))`

- Picks a random integer from `0` up to but not including `len(words)`.

`words[...]`

- Uses that random index to return a word from the slice.

### Go lesson

Indexing into a slice looks like:

```go
slice[index]
```

That is how we access one specific item in a collection.

## Lines 218-236: Reading User Input

```go
func readGuess(reader *bufio.Reader) (rune, error) {
	fmt.Print("Enter a letter: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return 0, nil
	}

	letters := []rune(input)
	if len(letters) != 1 {
		return 0, nil
	}

	return letters[0], nil
}
```

### Explanation

This function is the input layer.

`reader *bufio.Reader`

- Accepts a buffered reader connected to terminal input.

`fmt.Print("Enter a letter: ")`

- Shows the prompt.

`reader.ReadString('\n')`

- Reads input until the user presses Enter.

`input, err := ...`

- Go commonly returns a value plus an error.
- This is a core Go pattern.

`if err != nil`

- `nil` means "no value" or "no error" for certain types.
- If `err` is not nil, something went wrong while reading.

`strings.TrimSpace(input)`

- Removes trailing newline and surrounding whitespace.

`if input == ""`

- Rejects empty input.

`letters := []rune(input)`

- Converts the string into a slice of runes.
- This lets the function count actual characters rather than raw bytes.

`if len(letters) != 1`

- Requires exactly one character of input.

`return letters[0], nil`

- Returns the guessed rune and no error.

### Go lesson: error handling

Go does not use exceptions in the typical Java/C++ sense. Instead, many functions return an `error` value explicitly. You check it right away.

That pattern looks like:

```go
value, err := someFunction()
if err != nil {
	return err
}
```

You will see that everywhere in Go.

### CLI lesson

Separating input parsing into its own function makes the main loop much easier to read.

## Lines 238-276: The Main Program Loop

```go
func main() {
	reader := bufio.NewReader(os.Stdin)
	game := newGame(randomWord())

	fmt.Println(colorize("Hangman", colorBold+colorCyan))
	fmt.Println(colorize("Guess the word before the drawing is complete.", colorBlue))

	for !game.won() && !game.lost() {
		fmt.Println(colorize(hangmanArt(game.wrongGuesses), colorCyan))
		fmt.Printf("%s %s\n", colorize("Word:", colorBold), game.displayWord())
		fmt.Printf("%s %s\n", colorize("Guessed:", colorBold), colorize(game.guessedLetters(), colorBlue))
		fmt.Printf("%s %s\n\n", colorize("Wrong guesses left:", colorBold), formatRemainingGuesses(maxWrongGuesses-game.wrongGuesses))

		letter, err := readGuess(reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "input error: %v\n", err)
			os.Exit(1)
		}

		if letter == 0 {
			fmt.Println(colorize("Please enter one character.\n", colorYellow))
			continue
		}

		correct, message := game.guess(letter)
		if correct {
			fmt.Printf("%s\n\n", colorize(message, colorGreen))
		} else {
			fmt.Printf("%s\n\n", colorize(message, colorRed))
		}
	}

	fmt.Println(colorize(hangmanArt(game.wrongGuesses), colorCyan))
	if game.won() {
		fmt.Printf("%s\n", colorize(fmt.Sprintf("You won. The word was %q.", game.word), colorBold+colorGreen))
	} else {
		fmt.Printf("%s\n", colorize(fmt.Sprintf("You lost. The word was %q.", game.word), colorBold+colorRed))
	}
}
```

### This is the heart of the CLI app

If you only remember one thing about terminal game architecture, remember this section.

### Step 1: Setup input

`reader := bufio.NewReader(os.Stdin)`

- Creates a buffered reader around standard input.
- `os.Stdin` is the terminal input stream.

### Step 2: Setup game state

`game := newGame(randomWord())`

- Picks a random word
- creates a new game state object

### Step 3: Print intro text

These `fmt.Println` calls display the title and instructions.

### Step 4: Run the game loop

`for !game.won() && !game.lost() {`

- This loop keeps running while the player has neither won nor lost.
- In plain English:
  "Keep playing until one end condition becomes true."

Inside the loop, the program follows a very standard CLI game pattern:

1. Render current state
2. Read input
3. Validate input
4. Update state
5. Print feedback

### Rendering phase

`hangmanArt(game.wrongGuesses)`

- Draw current hangman stage.

`game.displayWord()`

- Show guessed letters and hidden blanks.

`game.guessedLetters()`

- Show what the user already typed.

`formatRemainingGuesses(...)`

- Show remaining attempts in a color based on urgency.

### Input phase

`letter, err := readGuess(reader)`

- Reads one input guess from the user.

### Error phase

```go
if err != nil {
	fmt.Fprintf(os.Stderr, "input error: %v\n", err)
	os.Exit(1)
}
```

- Prints an error to standard error
- exits the program with a non-zero status code

### Validation phase

```go
if letter == 0 {
	fmt.Println(colorize("Please enter one character.\n", colorYellow))
	continue
}
```

- `0` here is used as a sentinel value meaning "no valid character was provided."
- `continue` skips the rest of the loop and starts the next iteration.

### Update phase

`correct, message := game.guess(letter)`

- Applies the guess to the game state.
- Returns whether it was correct and what message to show.

### Feedback phase

If the guess was correct, print a green message. Otherwise print a red message.

### Final result

After the loop ends, the program prints the final board and then prints either the win or loss message.

### Go lesson

`main()` in Go should usually feel like orchestration, not clutter.

Good `main()` functions:

- create dependencies
- initialize state
- call other functions
- handle top-level control flow

That is what this one does.

## Big Picture: How CLI Apps Are Structured

This program is small, but it already uses a strong architecture pattern:

### 1. State

Stored in:

- `game.word`
- `game.guessed`
- `game.wrongGuesses`

### 2. Input

Handled by:

- `readGuess()`

### 3. Rules / business logic

Handled by:

- `guess()`
- `won()`
- `lost()`

### 4. Rendering

Handled by:

- `displayWord()`
- `guessedLetters()`
- `hangmanArt()`
- terminal printing inside `main()`

That separation is exactly how many larger terminal apps and games are organized.

## Important Go Syntax You Saw In This Project

### Variable declaration

```go
reader := bufio.NewReader(os.Stdin)
```

- `:=` means "declare and initialize"
- Go infers the type automatically

### Explicit variable declaration

```go
var builder strings.Builder
```

- Used when you want to declare a variable with a clear type

### Function declaration

```go
func randomWord() string
```

- function name: `randomWord`
- no parameters
- returns a `string`

### Method declaration

```go
func (g *game) won() bool
```

- method on `game`
- returns a `bool`

### Slice literal

```go
[]string{"a", "b", "c"}
```

### Map creation

```go
make(map[rune]bool)
```

### For loop

```go
for !game.won() && !game.lost() {
}
```

Go only has one looping keyword: `for`.

### Range loop

```go
for _, letter := range g.word
```

Used to iterate through collections.

### Conditional

```go
if err != nil {
}
```

### Switch

```go
switch {
case remaining >= 4:
}
```

## What You Should Learn From This Project

If you are new to Go, this project teaches several practical habits:

### 1. Keep state together

Use a struct when several values belong to the same concept.

### 2. Separate reading input from applying rules

This prevents `main()` from becoming a mess.

### 3. Use helper functions for rendering

If the screen output has meaning, move it into named functions.

### 4. Normalize input early

Converting to lowercase immediately keeps later logic simpler.

### 5. Treat errors explicitly

This is a core Go habit. Check errors right away.

## How You Could Improve This CLI App Next

If you want to grow this into a more advanced terminal game, natural next steps are:

- load words from a file instead of hardcoding them
- add replay support
- add difficulty levels
- add a no-color mode for terminals that do not support ANSI colors
- clear and redraw the terminal between turns for a cleaner UI
- split the program into multiple files such as `game.go`, `ui.go`, and `main.go`

## Final Mental Model

You can think of this program as three layers:

### Input layer

- `readGuess()`

### Game logic layer

- `newGame()`
- `guess()`
- `won()`
- `lost()`
- `randomWord()`

### Display layer

- `displayWord()`
- `guessedLetters()`
- `hangmanArt()`
- `colorize()`
- the `fmt.Print*` calls in `main()`

That is a solid foundation for both Go programs and terminal apps in general.

If you want, the next useful step would be for me to add a second teaching file that explains:

- how to refactor this into multiple Go files
- how to add tests in Go with `_test.go`
- how `go run`, `go build`, and `go mod` fit together
