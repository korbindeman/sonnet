package render

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type Coord struct {
	Line int
	Col  int
}

func NewCoord(line, col int) *Coord {
	return &Coord{line, col}
}

func (coord *Coord) Move(deltaLine, deltaCol int) {
	coord.Line += deltaLine
	coord.Col += deltaCol
}

func (coord *Coord) IsValid(maxLines, maxCols int) bool {
	return coord.Line >= 0 && coord.Line < maxLines && coord.Col >= 0 && coord.Col < maxCols
}

func ClearScreen() {
	fmt.Print("\x1b[2J")
}

func ClearLine() {
	fmt.Print("\x1b[K")
}

func GetTerminalSize() (int, int, error) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, 0, fmt.Errorf("error getting window size: %w", err)
	}
	return width, height, nil
}

func MoveCursor(coord *Coord) {
	fmt.Printf("\x1b[%d;%dH", coord.Line, coord.Col)
}
