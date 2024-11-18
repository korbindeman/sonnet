package window

import (
	"fmt"

	"github.com/korbindeman/sonnet/internal/buffer"
	"github.com/korbindeman/sonnet/internal/utils"
)

type Window struct {
	height int
	width  int
	Buffer buffer.Buffer
}

func NewWindow() *Window {
	width, height, _ := utils.GetWindowSize()

	window := &Window{height, width, *buffer.NewBuffer()}

	window.SetCursor()

	return window
}

func (w *Window) GetSize() (int, int) {
	return w.width, w.height
}

func (w *Window) SetCursor() {
	line, col := w.Buffer.GetCursor()
	utils.MoveCursor(line, col)
}

func (w *Window) DisplayBuffer() {
	utils.ClearScreen()
	utils.MoveCursor(1, 1)
	for i, line := range w.Buffer.Rows() {
		if i >= w.height-1 {
			break
		}
		fmt.Print(line, "\r\n")
	}
	w.SetCursor()
}

func (w *Window) LoadBuffer(buffer *buffer.Buffer) {
	w.Buffer = *buffer
	w.DisplayBuffer()
}
