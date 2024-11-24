package buffer

import (
	"github.com/korbindeman/sonnet/internal/render"
	"github.com/korbindeman/sonnet/internal/utils"
)

func (b *Buffer) GetCursor() render.Coord {
	return b.cursor
}

func (b *Buffer) SetCursor(line, col int) {
	b.cursor.Line = line
	b.cursor.Col = col
	b.virtualCol = col
}

func (b *Buffer) MoveDown() {
	if b.cursor.Line+1 <= b.Length() {
		b.cursor.Line++
	}
	if b.CurrentLineLength() < b.cursor.Col {
		b.cursor.Col = b.CurrentLineLength()
	}
	if b.virtualCol >= b.cursor.Col {
		b.cursor.Col = utils.Min(b.virtualCol, b.CurrentLineLength())
	}
	b.cursor.Col = utils.Max(b.cursor.Col, 1)
}

func (b *Buffer) MoveUp() {
	if b.cursor.Line > 1 {
		b.cursor.Line--
	}
	if b.CurrentLineLength() < b.cursor.Col {
		b.cursor.Col = b.CurrentLineLength()
	}
	if b.virtualCol >= b.cursor.Col {
		b.cursor.Col = utils.Min(b.virtualCol, b.CurrentLineLength())
	}
	b.cursor.Col = utils.Max(b.cursor.Col, 1)
}

func (b *Buffer) MoveRight() {
	if b.cursor.Col < b.CurrentLineLength() {
		b.cursor.Col++
	}
	b.virtualCol = b.cursor.Col
}

func (b *Buffer) MoveLeft() {
	if b.cursor.Col > 1 {
		b.cursor.Col--
	}
	b.virtualCol = b.cursor.Col
}
