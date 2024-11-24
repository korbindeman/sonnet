package window

import (
	"fmt"

	"github.com/korbindeman/sonnet/internal/buffer"
	"github.com/korbindeman/sonnet/internal/render"
)

type Window struct {
	height       int
	width        int
	Buffer       buffer.Buffer
	anchor       render.Coord
	scrollOffset int
}

func NewWindow(height, width int, buffer *buffer.Buffer, anchor render.Coord) *Window {
	window := &Window{height, width, *buffer, anchor, 0}

	window.SetCursor()
	window.DisplayBuffer()

	return window
}

func NewFullscreenWindow(buffer *buffer.Buffer) *Window {
	width, height, _ := render.GetTerminalSize()

	window := &Window{height - 1, width - 1, *buffer, render.Coord{Line: 1, Col: 1}, 0}

	window.SetCursor()
	window.DisplayBuffer()

	return window
}

func (w *Window) GetSize() (int, int) {
	return w.width, w.height
}

func (w *Window) SetCursor() {
	cursor := w.Buffer.GetCursor()

	line, col := cursor.Line-w.scrollOffset, cursor.Col

	w.moveCursorWithLineNumbers(render.Coord{Line: line, Col: col})
}

func (w *Window) moveCursor(coord render.Coord) {
	if coord.Line < 1 {
		coord.Line = 1
		w.ScrollUp()
	}
	if coord.Line > w.height {
		coord.Line = w.height
		w.ScrollDown()
	}

	// TODO: handle col

	line := coord.Line + w.anchor.Line
	col := coord.Col + w.anchor.Col

	render.MoveCursor(render.NewCoord(line, col))
}

func (w *Window) moveCursorWithLineNumbers(coord render.Coord) {
	lineNumLen := 5
	line := coord.Line
	col := coord.Col + lineNumLen
	w.moveCursor(render.Coord{Line: line, Col: col})
}

func (w *Window) renderLine(line string, lineNum int, termination string) {
	fmt.Printf("%4d ", lineNum)
	if len(line) > w.width-5 {
		line = line[:w.width-5]
	}
	fmt.Print(line, termination)
}

func (w *Window) DisplayBuffer() {
	render.ClearScreen()
	for i := 0; i < w.height; i++ {
		indexWithOffset := i + w.scrollOffset
		if indexWithOffset >= len(w.Buffer.Rows()) {
			break
		}
		line := w.Buffer.Rows()[indexWithOffset]
		w.moveCursor(render.Coord{Line: i + 1, Col: 1})
		w.renderLine(line, indexWithOffset+1, "\r\n")
	}
	w.SetCursor()
}

func (w *Window) LoadBuffer(buffer *buffer.Buffer) {
	w.Buffer = *buffer
	w.DisplayBuffer()
}

func (w *Window) ScrollUp() {
	if w.scrollOffset > 0 {
		w.scrollOffset--
	}
	w.DisplayBuffer()
}

func (w *Window) ScrollDown() {
	if w.scrollOffset < len(w.Buffer.Rows())-w.height+1 {
		w.scrollOffset++
	}
	w.DisplayBuffer()
}
