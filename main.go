package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

func printScreen(s tcell.Screen, row, col int, str string) {
	for _, c := range str {
		s.SetContent(col, row, c, nil, tcell.StyleDefault)
		col++
	}
}

func print(s tcell.Screen, row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			s.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func displayHelloWorld(screen tcell.Screen) {
	screen.Clear()
	printScreen(screen, 1, 5, "Hello, World!")
	print(screen, 0, 0, 5, 5, '*')
	screen.Show()
}

func main() {
	screen := initScreen()

	displayHelloWorld(screen)

	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
			displayHelloWorld(screen)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				screen.Fini()
				os.Exit(0)
			}
		}
	}

	// Draw paddles
	// Player movement
	// Take care of paddle boundaries
	// Draw ball
	// Update ball movement
	// Handle collisions
	// Handle game over

}

func initScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
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

	return screen
}
