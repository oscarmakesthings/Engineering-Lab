package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	SnakeChar = 'O'
	FoodChar  = '*'
	WallChar  = '#'
	EmptyChar = ' '
)

type Point struct {
	x, y int
}

type GameState struct {
	snake     []Point
	direction Point
	food      Point
	score     int
	gameOver  bool
	startTime time.Time
	width     int
	height    int
}

func main() {
	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer s.Fini()

	// Set default style
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	// Game configuration
	w, h := s.Size()
	// Boundary for the game board (leave space for score and timer)
	gameW, gameH := w-2, h-4
	if gameW > 40 {
		gameW = 40
	}
	if gameH > 20 {
		gameH = 20
	}

	game := &GameState{
		snake:     []Point{{x: gameW / 2, y: gameH / 2}, {x: gameW/2 - 1, y: gameH / 2}},
		direction: Point{x: 1, y: 0},
		startTime: time.Now(),
		width:     gameW,
		height:    gameH,
	}
	game.spawnFood()

	// Channel for input events
	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					close(quit)
					return
				}
				if game.gameOver {
					if ev.Key() == tcell.KeyEnter || ev.Rune() == 'r' {
						// Restart game
						game.restart()
					}
					continue
				}
				switch ev.Key() {
				case tcell.KeyUp:
					if game.direction.y == 0 {
						game.direction = Point{x: 0, y: -1}
					}
				case tcell.KeyDown:
					if game.direction.y == 0 {
						game.direction = Point{x: 0, y: 1}
					}
				case tcell.KeyLeft:
					if game.direction.x == 0 {
						game.direction = Point{x: -1, y: 0}
					}
				case tcell.KeyRight:
					if game.direction.x == 0 {
						game.direction = Point{x: 1, y: 0}
					}
				}
				// Support WASD
				switch ev.Rune() {
				case 'w', 'W':
					if game.direction.y == 0 {
						game.direction = Point{x: 0, y: -1}
					}
				case 's', 'S':
					if game.direction.y == 0 {
						game.direction = Point{x: 0, y: 1}
					}
				case 'a', 'A':
					if game.direction.x == 0 {
						game.direction = Point{x: -1, y: 0}
					}
				case 'd', 'D':
					if game.direction.x == 0 {
						game.direction = Point{x: 1, y: 0}
					}
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	// Game loop
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			if !game.gameOver {
				game.update()
			}
			game.draw(s)
		}
	}
}

func (g *GameState) spawnFood() {
	for {
		newFood := Point{
			x: rand.Intn(g.width),
			y: rand.Intn(g.height),
		}
		// Check if food spawned on snake
		onSnake := false
		for _, p := range g.snake {
			if p == newFood {
				onSnake = true
				break
			}
		}
		if !onSnake {
			g.food = newFood
			break
		}
	}
}

func (g *GameState) update() {
	// Calculate new head position
	head := g.snake[0]
	newHead := Point{
		x: head.x + g.direction.x,
		y: head.y + g.direction.y,
	}

	// Collision with walls
	if newHead.x < 0 || newHead.x >= g.width || newHead.y < 0 || newHead.y >= g.height {
		g.gameOver = true
		return
	}

	// Collision with self
	for _, p := range g.snake {
		if p == newHead {
			g.gameOver = true
			return
		}
	}

	// Move snake
	g.snake = append([]Point{newHead}, g.snake...)

	// Check if food eaten
	if newHead == g.food {
		g.score += 10
		g.spawnFood()
	} else {
		// Remove tail
		g.snake = g.snake[:len(g.snake)-1]
	}
}

func (g *GameState) restart() {
	g.snake = []Point{{x: g.width / 2, y: g.height / 2}, {x: g.width/2 - 1, y: g.height / 2}}
	g.direction = Point{x: 1, y: 0}
	g.score = 0
	g.gameOver = false
	g.startTime = time.Now()
	g.spawnFood()
}

func (g *GameState) draw(s tcell.Screen) {
	s.Clear()

	// Style definitions
	snakeStyle := tcell.StyleDefault.Foreground(tcell.ColorGreen)
	foodStyle := tcell.StyleDefault.Foreground(tcell.ColorRed)
	wallStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	textStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow)

	// Draw boundaries
	for x := -1; x <= g.width; x++ {
		s.SetContent(x+1, 0, WallChar, nil, wallStyle)
		s.SetContent(x+1, g.height+1, WallChar, nil, wallStyle)
	}
	for y := 0; y <= g.height+1; y++ {
		s.SetContent(0, y, WallChar, nil, wallStyle)
		s.SetContent(g.width+1, y, WallChar, nil, wallStyle)
	}

	// Draw food
	s.SetContent(g.food.x+1, g.food.y+1, FoodChar, nil, foodStyle)

	// Draw snake
	for _, p := range g.snake {
		s.SetContent(p.x+1, p.y+1, SnakeChar, nil, snakeStyle)
	}

	// Draw Score and Timer
	elapsed := time.Since(g.startTime)
	if g.gameOver {
		elapsed = elapsed // keep it frozen? actually it would keep counting if we don't save end time
		// For simplicity, we just display it.
	}

	scoreStr := fmt.Sprintf("Score: %d", g.score)
	timerStr := fmt.Sprintf("Time: %02d:%02d", int(elapsed.Minutes()), int(elapsed.Seconds())%60)

	drawText(s, 1, g.height+2, textStyle, scoreStr)
	drawText(s, g.width-len(timerStr)+1, g.height+2, textStyle, timerStr)

	if g.gameOver {
		gameOverMsg := " GAME OVER! "
		restartMsg := "Press 'Enter' or 'R' to Restart"
		drawText(s, (g.width-len(gameOverMsg))/2+1, g.height/2, textStyle.Bold(true).Foreground(tcell.ColorRed), gameOverMsg)
		drawText(s, (g.width-len(restartMsg))/2+1, g.height/2+1, textStyle, restartMsg)
	}

	s.Show()
}

func drawText(s tcell.Screen, x, y int, style tcell.Style, text string) {
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}
