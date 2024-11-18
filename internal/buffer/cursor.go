package buffer

import "github.com/korbindeman/sonnet/internal/utils"

func (b *Buffer) GetCursor() (int, int) {
	return b.line, b.col
}

func (b *Buffer) SetCursor(line, col int) {
	b.line = line
	b.col = col
	b.virtualCol = col
}

func (b *Buffer) MoveDown() {
	if b.line < len(b.rows) {
		b.line++
	}
	if b.CurrentLineLength() < b.col {
		b.col = b.CurrentLineLength()
	}
	if b.virtualCol >= b.col {
		b.col = utils.Min(b.virtualCol, b.CurrentLineLength())
	}
}

func (b *Buffer) MoveUp() {
	if b.line > 1 {
		b.line--
	}
	if b.CurrentLineLength() < b.col {
		b.col = b.CurrentLineLength()
	}
	if b.virtualCol >= b.col {
		b.col = utils.Min(b.virtualCol, b.CurrentLineLength())
	}
}

func (b *Buffer) MoveRight() {
	if b.col < b.CurrentLineLength() {
		b.col++
	}
	b.virtualCol = b.col
}

func (b *Buffer) MoveLeft() {
	if b.col > 1 {
		b.col--
	}
	b.virtualCol = b.col
}
