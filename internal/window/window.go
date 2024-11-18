package window

import (
	"github.com/korbindeman/sonnet/internal/buffer"
	"github.com/korbindeman/sonnet/internal/editor"
	"github.com/korbindeman/sonnet/internal/utils"
)

type Window struct {
	height int
	width  int
	Cursor Cursor
}

func NewWindow() *Window {
	width, height, _ := utils.GetWindowSize()

	window := &Window{height, width, *NewCursor()}

	window.SetCursor()

	return window
}

func (w *Window) GetSize() (int, int) {
	return w.width, w.height
}

func (w *Window) SetCursor() {
	utils.MoveCursor(w.Cursor.line, w.Cursor.col)
}

func (w *Window) DisplayBuffer(buffer *buffer.Buffer) {
	editor.DisplayBuffer(buffer, w.width, w.height)
}
