package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

const paddleHeight = 4
const paddleSymbol = 0x2588

type Paddle struct {
	row, col, width, height int
}

var screen tcell.Screen
var player1 *Paddle
var player2 *Paddle

func PrintString(col, row int, str string) {
	for _, c := range str {
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}

func Print(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func DrawState() {
	screen.Clear()
	Print(player1.row, player1.col, player1.width, player1.height, paddleSymbol)
	Print(player2.row, player2.col, player2.width, player2.height, paddleSymbol)
	screen.Show()
}

func main() {
	InitScreen()
	InitGameState()

	for {
		DrawState()
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventKey:
			if ev.Rune() == 'q' {
				screen.Fini()
				os.Exit(0)
			} else if ev.Key() == tcell.KeyUp {
				player2.row--
			} else if ev.Key() == tcell.KeyDown {
				player2.row++
			} else if ev.Rune() == 'w' {
				player1.row--
			} else if ev.Rune() == 's' {
				player1.row++
			}
		}
	}
}

func InitScreen() {
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

func InitGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - paddleHeight/2

	player1 = &Paddle{
		row: paddleStart, col: 0, width: 1, height: paddleHeight,
	}
	player2 = &Paddle{
		row: paddleStart, col: width - 1, width: 1, height: paddleHeight,
	}
}
