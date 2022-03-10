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

// var debugLog string
var isGamePaused bool

var gameObjects []*GameObject

func main() {
	initScreen()
	initGameState()
	inputChan := initUserInput()

	for !isGameOver() {
		handleUserInput(readInput(inputChan))
		updateState()
		drawState()

		time.Sleep(100 * time.Millisecond)
	}

	screenWidth, screenHeight := screen.Size()
	winner := getWinner()
	printStringCenter(screenHeight/2-1, screenWidth/2, "Game over")
	printStringCenter(screenHeight/2, screenWidth/2, fmt.Sprintf("%s wins...", winner))
	screen.Show()
	time.Sleep(3 * time.Second)
	screen.Fini()
}

func drawState() {
	if isGamePaused {
		return
	}

	screen.Clear()

	// printString(0, 0, debugLog)
	for _, obj := range gameObjects {
		print(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}

	screen.Show()
}

func collidesWithWall(obj *GameObject) bool {
	_, screenHeight := screen.Size()

	return obj.row+obj.velRow < 0 || obj.row+obj.velRow >= screenHeight
}

func collidesWithPaddle(ball, paddle *GameObject) bool {
	var collidesOnColumn bool

	if ball.col < paddle.col {
		collidesOnColumn = ball.col+ball.velCol >= paddle.col
	} else {
		collidesOnColumn = ball.col+ball.velCol <= paddle.col
	}
	return collidesOnColumn &&
		ball.row >= paddle.row &&
		ball.row < paddle.row+paddle.height
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
	} else if key == "Rune[p]" {
		isGamePaused = !isGamePaused
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

func isGameOver() bool {
	return getWinner() != ""
}

func getWinner() string {
	screenWidth, _ := screen.Size()
	if ball.col < 0 {
		return "Player 2"
	} else if ball.col >= screenWidth {
		return "Player 1"
	} else {
		return ""
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
	if isGamePaused {
		return
	}

	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}

	// debugLog = fmt.Sprintf(
	// 	"ball: row=%d, col=%d\npaddle1: row=%d, col=%d\npaddle2: row=%d, col=%d",
	// 	ball.row, ball.col,
	// 	player1Paddle.row, player1Paddle.col,
	// 	player2Paddle.row, player2Paddle.col,
	// )

	if collidesWithWall(ball) {
		ball.velRow = -ball.velRow
	}

	if collidesWithPaddle(ball, player1Paddle) || collidesWithPaddle(ball, player2Paddle) {
		ball.velCol = -ball.velCol
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

func printStringCenter(row, col int, str string) {
	col -= len(str) / 2

	printString(row, col, str)
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
