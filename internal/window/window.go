package window

import (
	"fmt"

	"github.com/korbindeman/sonnet/internal/buffer"
	"github.com/korbindeman/sonnet/internal/render"
)

type Window struct {
	height int
	width  int
	Buffer buffer.Buffer
	anchor render.Coord
}

func NewWindow(height, width int, buffer *buffer.Buffer, anchor render.Coord) *Window {
	window := &Window{height, width, *buffer, anchor}

	window.SetCursor()

	window.DisplayBuffer()

	return window
}

func NewFullscreenWindow(buffer *buffer.Buffer) *Window {
	width, height, _ := render.GetTerminalSize()

	window := &Window{height, width, *buffer, render.Coord{Line: 1, Col: 1}}

	window.SetCursor()

	window.DisplayBuffer()

	return window
}

func (w *Window) GetSize() (int, int) {
	return w.width, w.height
}

func (w *Window) SetCursor() {
	cursor := w.Buffer.GetCursor()
	w.moveCursorWithLineNumbers(cursor)
}

func (w *Window) moveCursor(coord render.Coord) {
	if coord.Line < 0 {
		coord.Line = 0
	}
	if coord.Col < 0 {
		coord.Col = 0
	}
	if coord.Line >= w.height {
		coord.Line = w.height - 1
	}
	if coord.Col >= w.width {
		coord.Col = w.width - 1
	}
	line := coord.Line + w.anchor.Line
	col := coord.Col + w.anchor.Col
	render.MoveCursor(render.NewCoord(line, col))
}

func (w *Window) moveCursorWithLineNumbers(coord render.Coord) {
	lineNumLen := 5
	line := coord.Line + w.anchor.Line - 1
	col := coord.Col + w.anchor.Col + lineNumLen
	render.MoveCursor(render.NewCoord(line, col))
}

func renderLine(line string, lineNum int, termination string) {
	fmt.Printf("%4d ", lineNum)
	fmt.Print(line, termination)
}

func (w *Window) DisplayBuffer() {
	render.ClearScreen()
	for i, line := range w.Buffer.Rows() {
		if i >= w.height-1 {
			break
		}
		w.moveCursor(render.Coord{Line: i, Col: 1})
		renderLine(line, i+1, "\r\n")
	}
	w.SetCursor()
}

func (w *Window) LoadBuffer(buffer *buffer.Buffer) {
	w.Buffer = *buffer
	w.DisplayBuffer()
}
