package main

import (
	"github.com/korbindeman/sonnet/internal/buffer"
	"github.com/korbindeman/sonnet/internal/io"
	"github.com/korbindeman/sonnet/internal/keymaps"
	"github.com/korbindeman/sonnet/internal/render"
	"github.com/korbindeman/sonnet/internal/window"
)

func main() {
	sonnet := Initialize()
	defer sonnet.Close()

	keyBindings := keymaps.NewDefaultKeyBindings()

	// buf, _ := buffer.LoadFile("main.go")
	buf, _ := buffer.LoadFile("internal/window/window.go")

	// window := window.NewFullscreenWindow(buf)
	window := window.NewWindow(10, 50, buf, render.Coord{Line: 5, Col: 5})

	for {
		input, _ := io.ReadInput()

		handler := keymaps.HandleInput(keyBindings, input)

		handler(window)
	}
}
