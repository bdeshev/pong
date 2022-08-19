package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

const paddleHeight = 4
const paddleSymbol = 0x2588

func PrintString(s tcell.Screen, col, row int, str string) {
	for _, c := range str {
		s.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}

func Print(s tcell.Screen, row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			s.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func displayPaddles(screen tcell.Screen) {
	screen.Clear()
	width, height := screen.Size()
	paddleStart := height/2 - paddleHeight/2
	Print(screen, paddleStart, 0, 1, paddleHeight, paddleSymbol)
	Print(screen, paddleStart, width-1, 1, paddleHeight, paddleSymbol)
	screen.Show()
}

func main() {
	screen := InitScreen()
	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyESC {
				screen.Fini()
				os.Exit(0)
			}
		}
	}
}

func InitScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if e := screen.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)

	displayPaddles(screen)

	return screen
}
