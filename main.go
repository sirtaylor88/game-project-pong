package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const PaddleSymbol = 0x2588
const BallSymbol = 0x25cf

const PaddleHeight = 4

const InitialBallVelocityRow = 1
const InitialBallVelocityCol = 2

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int
	symbol                  rune
}

var screen tcell.Screen
var player1Paddle *GameObject
var player2Paddle *GameObject
var ball *GameObject
var debugLog string

var gameObjects []*GameObject

func main() {
	initScreen()
	initGameState()
	inputChan := initUserInput()

	for {
		handleUserInput(readInput(inputChan))
		updateState()
		drawState()

		time.Sleep(75 * time.Millisecond)
	}
	// Handle collisions
	// Handle game over

}

func drawState() {
	screen.Clear()

	printString(0, 0, debugLog)
	for _, obj := range gameObjects {
		print(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}

	screen.Show()
}

func handleUserInput(key string) {
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[z]" && player1Paddle.isPaddleInsideBoundary("up") {
		player1Paddle.row--
	} else if key == "Rune[s]" && player1Paddle.isPaddleInsideBoundary("down") {
		player1Paddle.row++
	} else if key == "Up" && player2Paddle.isPaddleInsideBoundary("up") {
		player2Paddle.row--
	} else if key == "Down" && player2Paddle.isPaddleInsideBoundary("down") {
		player2Paddle.row++
	}
}

func (playerPaddle *GameObject) isPaddleInsideBoundary(direction string) bool {
	_, screenHeight := screen.Size()
	if direction == "up" {
		return playerPaddle.row > 0
	} else {
		return playerPaddle.row+PaddleHeight < screenHeight
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

	player1Paddle = &GameObject{
		row:    paddleStart,
		col:    0,
		width:  1,
		height: PaddleHeight,
		velRow: 0,
		velCol: 0,
		symbol: PaddleSymbol,
	}

	player2Paddle = &GameObject{
		row:    paddleStart,
		col:    width - 1,
		width:  1,
		height: PaddleHeight,
		velRow: 0,
		velCol: 0,
		symbol: PaddleSymbol,
	}

	ball = &GameObject{
		row:    height / 2,
		col:    width / 2,
		width:  1,
		height: 1,
		velRow: InitialBallVelocityRow,
		velCol: InitialBallVelocityCol,
		symbol: BallSymbol,
	}

	gameObjects = []*GameObject{
		player1Paddle, player2Paddle, ball,
	}
}

func updateState() {
	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
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
