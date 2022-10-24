package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const PaddleHeight = 4
const InitialBallVelocityRow = 1
const InitialBallVelocityCol = 2

const PaddleSymbol = 0x2588
const BallSymbol = 0x25Cf

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int
	symbol                  rune
}

var screen tcell.Screen
var player1Paddle *GameObject
var player2Paddle *GameObject
var ball *GameObject
var isGamePaused bool
var debugLog string

var gameObjects []*GameObject

func main() {
	InitScreen()
	InitGameState()
	inputChan := UserInput()

	for {
		key := ReadInput(inputChan)
		DrawState()
		UpdateState()
		HandleUserInput(key)
		time.Sleep(300 * time.Millisecond)
	}

}

func PrintString(row, column int, str string) {
	for _, c := range str {
		screen.SetContent(column, row, c, nil, tcell.StyleDefault)
		column += 1
	}
}
func Print(row, column, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(column+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func UpdateState() {
	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}

	if ClollidesWithWall(ball) {
		ball.velRow = -ball.velRow
	}

	if CollidesWithPaddle(ball, player1Paddle) || CollidesWithPaddle(ball, player2Paddle) {
		ball.velCol = -ball.col
	}
}

func DrawState() {
	screen.Clear()
	PrintString(0, 0, debugLog)
	for _, obj := range gameObjects {
		Print(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}
	screen.Show()

}

func ClollidesWithWall(obj *GameObject) bool {
	_, screenHeight := screen.Size()
	return !(obj.row+obj.velRow >= 0 && obj.row+obj.velRow < screenHeight)
}

func CollidesWithPaddle(ball *GameObject, paddle *GameObject) bool {
	var collidesCol bool
	if ball.col < paddle.col {
		collidesCol = ball.col+ball.velCol >= paddle.col

	} else {
		collidesCol = ball.col+ball.velCol <= paddle.col
	}
	return collidesCol &&
		ball.row+ball.velRow >= paddle.row &&
		ball.row+ball.velRow < paddle.row+paddle.height
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

func HandleUserInput(key string) {
	_, screenHeight := screen.Size()
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[w]" && player1Paddle.row > 0 {
		player1Paddle.row--
	} else if key == "Rune[s]" && player1Paddle.row+player1Paddle.height < screenHeight {
		player1Paddle.row++
	} else if key == "Up" && player2Paddle.row > 0 {
		player2Paddle.row--
	} else if key == "Down" && player2Paddle.row+player2Paddle.height < screenHeight {
		player2Paddle.row++
	} else if key == "Esc" {
		isGamePaused = true
	}
}

func UserInput() chan string {
	inputChannel := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChannel <- ev.Name()
			}
		}
	}()
	return inputChannel
}

func InitGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - PaddleHeight/2

	player1Paddle = &GameObject{
		row: paddleStart, col: 0, width: 1, height: PaddleHeight,
		velRow: 0, velCol: 0,
		symbol: PaddleSymbol,
	}
	player2Paddle = &GameObject{
		row: paddleStart, col: width - 1, width: 1, height: PaddleHeight,
		velRow: 0, velCol: 0,
		symbol: PaddleSymbol,
	}
	ball = &GameObject{
		row: height / 2, col: width / 2, width: 1, height: 1,
		velRow: InitialBallVelocityRow, velCol: InitialBallVelocityCol,
		symbol: BallSymbol,
	}
	gameObjects = []*GameObject{
		player1Paddle, player2Paddle, ball,
	}
}

func ReadInput(inputChannel chan string) string {
	var key string
	select {
	case key = <-inputChannel:
	default:
		key = ""
	}
	return key
}
