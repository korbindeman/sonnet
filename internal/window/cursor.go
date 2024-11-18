package window

import "github.com/korbindeman/sonnet/internal/utils"

type Cursor struct {
	line int
	col  int
}

func NewCursor() *Cursor {
	return &Cursor{1, 1}
}

func (c *Cursor) Move(line, col int) {
	c.line = line
	c.col = col

	utils.MoveCursor(c.line, c.col)
}

func (c *Cursor) GetPosition() (int, int) {
	return c.line, c.col
}
