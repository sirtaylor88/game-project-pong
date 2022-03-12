package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const PaddleSymbol = 0x2588
const BallSymbol = 0x25cf

const PaddleHeight = 8

const BoundaryOffsetRow = 5
const BoundaryOffsetCol = 10

const InitialBallVelocityRow = 1
const InitialBallVelocityCol = 2

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int
	symbol                  rune
	color                   tcell.Style
}

var screen tcell.Screen
var player1Paddle *GameObject
var player2Paddle *GameObject
var ball *GameObject

// var debugLog string
var isGamePaused bool
var netColor tcell.Style
var tableColor tcell.Style

var gameObjects []*GameObject

func main() {
	rand.Seed(time.Now().UnixNano())
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
	printStringCenter(
		screenHeight/2, screenWidth/2,
		fmt.Sprintf("%s wins...", winner),
	)
	screen.Show()
	time.Sleep(3 * time.Second)
	screen.Fini()
}

func clearScreen() {
	for _, obj := range gameObjects {
		screenWidth, screenHeight := screen.Size()
		print(
			0, obj.col, obj.width, screenHeight,
			0x20, tcell.StyleDefault,
		)
		print(
			obj.row-obj.velRow, obj.col-obj.velCol,
			screenWidth, screenHeight, 0x20, tcell.StyleDefault,
		)
	}
}

func drawState() {
	if isGamePaused {
		return
	}

	// screen.Clear()
	clearScreen()

	// printString(0, 0, debugLog)
	printPongTable()

	screenWidth, _ := screen.Size()
	for _, obj := range gameObjects {
		if obj.col > BoundaryOffsetCol &&
			obj.col < screenWidth-BoundaryOffsetCol-1 {
			if obj.col == screenWidth/2 {
				print(
					obj.row, obj.col, obj.width, obj.height,
					obj.symbol, netColor,
				)
			} else {
				print(
					obj.row, obj.col, obj.width, obj.height,
					obj.symbol, obj.color,
				)
			}
		}
	}

	screen.Show()
}

func collidesWithWall(ball *GameObject) bool {
	_, screenHeight := screen.Size()

	return ball.row+ball.velRow < BoundaryOffsetRow ||
		ball.row+ball.velRow >= screenHeight-BoundaryOffsetRow
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
		return playerPaddle.row > BoundaryOffsetRow
	} else {
		return playerPaddle.row+PaddleHeight < screenHeight-BoundaryOffsetRow
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
		Background(tcell.ColorDarkBlue).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)

	screenWidth, screenHeight := screen.Size()
	if screenWidth <= 100 || screenHeight <= 30 {
		fmt.Fprintln(os.Stderr, "The screen is too small.")
		os.Exit(1)
	}
}

func initGameState() {
	screenWidth, screenHeight := screen.Size()
	paddleStart := screenHeight/2 - PaddleHeight/2

	netColor = tcell.StyleDefault.Background(tcell.ColorWhite)
	tableColor = tcell.StyleDefault.Background(tcell.ColorLightGreen)
	player1Paddle = &GameObject{
		row:    paddleStart,
		col:    BoundaryOffsetCol + 5,
		width:  1,
		height: PaddleHeight,
		velRow: 0,
		velCol: 0,
		symbol: PaddleSymbol,
		color:  tcell.StyleDefault.Foreground(tcell.ColorRed),
	}

	player2Paddle = &GameObject{
		row:    paddleStart,
		col:    screenWidth - BoundaryOffsetCol - 6,
		width:  1,
		height: PaddleHeight,
		velRow: 0,
		velCol: 0,
		symbol: PaddleSymbol,
		color:  tcell.StyleDefault.Foreground(tcell.ColorBlue),
	}

	ball = &GameObject{
		row:    screenHeight / 2,
		col:    screenWidth / 2,
		width:  1,
		height: 1,
		velRow: InitialBallVelocityRow + rand.Intn(1),
		velCol: InitialBallVelocityCol + rand.Intn(1),
		symbol: BallSymbol,
		color:  tableColor.Foreground(tcell.ColorYellow),
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
	// 	"ball: row=%d, col=%d\npaddle1: row=%d, col=%d\npaddle2: row=%d, col=%d\n",
	// 	ball.row, ball.col,
	// 	player1Paddle.row, player1Paddle.col,
	// 	player2Paddle.row, player2Paddle.col,
	// )

	if collidesWithWall(ball) {
		ball.velRow = -ball.velRow
	}

	if collidesWithPaddle(ball, player1Paddle) ||
		collidesWithPaddle(ball, player2Paddle) {
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

func print(row, col, width, height int, ch rune, color tcell.Style) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, color)
		}
	}
}

func printPongTable() {
	screenWidth, screenHeight := screen.Size()
	for r := BoundaryOffsetRow; r < screenHeight-BoundaryOffsetRow; r++ {
		for c := BoundaryOffsetCol + 1; c < screenWidth-BoundaryOffsetCol-1; c++ {
			if c == screenWidth/2 {
				screen.SetContent(c, r, 0x20, nil, netColor)
			} else {
				screen.SetContent(c, r, 0x20, nil, tableColor)
			}

		}
	}
}
