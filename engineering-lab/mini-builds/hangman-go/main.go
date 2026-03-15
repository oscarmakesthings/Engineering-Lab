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

type game struct {
	word        string
	guessed     map[rune]bool
	wrongGuesses int
}

func newGame(word string) *game {
	return &game{
		word:         strings.ToLower(word),
		guessed:      make(map[rune]bool),
		wrongGuesses: 0,
	}
}

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

func hangmanArt(wrongGuesses int) string {
	stages := []string{
		`
 +---+
 |   |
     |
     |
     |
     |
=========
`,
		`
 +---+
 |   |
 O   |
     |
     |
     |
=========
`,
		`
 +---+
 |   |
 O   |
 |   |
     |
     |
=========
`,
		`
 +---+
 |   |
 O   |
/|   |
     |
     |
=========
`,
		`
 +---+
 |   |
 O   |
/|\\  |
     |
     |
=========
`,
		`
 +---+
 |   |
 O   |
/|\\  |
/    |
     |
=========
`,
		`
 +---+
 |   |
 O   |
/|\\  |
/ \\  |
     |
=========
`,
	}

	if wrongGuesses < 0 {
		wrongGuesses = 0
	}
	if wrongGuesses >= len(stages) {
		wrongGuesses = len(stages) - 1
	}

	return stages[wrongGuesses]
}

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

func randomWord() string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	return words[rng.Intn(len(words))]
}

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
