package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const PaddleSymbol = 0x2588
const PaddleHeight = 4

type Paddle struct {
	row, col, width, height int
}

var screen tcell.Screen
var player1 *Paddle
var player2 *Paddle
var debugLog string

func main() {
	initScreen()
	initGameState()
	inputChan := initUserInput()

	for {
		drawState()
		time.Sleep(50 * time.Millisecond)

		key := readInput(inputChan)
		handleUserInput(key)
	}
	// Draw ball
	// Update ball movement
	// Handle collisions
	// Handle game over

}

func drawState() {
	screen.Clear()

	printString(0, 0, debugLog)
	print(player1.row, player1.col, player1.width, player1.height, PaddleSymbol)
	print(player2.row, player2.col, player2.width, player2.height, PaddleSymbol)

	screen.Show()
}

func handleUserInput(key string) {
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[z]" && player1.isPaddleInsideBoundary("up") {
		player1.row--
	} else if key == "Rune[s]" && player1.isPaddleInsideBoundary("down") {
		player1.row++
	} else if key == "Up" && player2.isPaddleInsideBoundary("up") {
		player2.row--
	} else if key == "Down" && player2.isPaddleInsideBoundary("down") {
		player2.row++
	}
}

func (player *Paddle) isPaddleInsideBoundary(direction string) bool {
	_, screenHeight := screen.Size()
	if direction == "up" {
		return player.row > 0
	} else {
		return player.row+PaddleHeight < screenHeight
	}
}

func initScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func initGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - PaddleHeight/2

	player1 = &Paddle{
		row:    paddleStart,
		col:    0,
		width:  1,
		height: PaddleHeight,
	}

	player2 = &Paddle{
		row:    paddleStart,
		col:    width - 1,
		width:  1,
		height: PaddleHeight,
	}
}

func initUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {

			switch ev := screen.PollEvent().(type) {
			case *tcell.EventResize:
				screen.Sync()
				drawState()
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()

	return inputChan
}

func readInput(inputChan chan string) string {
	var key string
	select {
	case key = <-inputChan:
	default:
		key = ""
	}

	return key
}

func printString(row, col int, str string) {
	for _, c := range str {
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col++
	}
}

func print(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}
