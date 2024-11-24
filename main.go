package main

import (
	"github.com/korbindeman/sonnet/internal/buffer"
	"github.com/korbindeman/sonnet/internal/io"
	"github.com/korbindeman/sonnet/internal/keymaps"
	"github.com/korbindeman/sonnet/internal/window"
)

func main() {
	sonnet := Initialize()
	defer sonnet.Close()

	keyBindings := keymaps.NewDefaultKeyBindings()

	buf, _ := buffer.LoadFile("main.go")

	window := window.NewFullscreenWindow(buf)

	for {
		input, _ := io.ReadInput()

		handler := keymaps.HandleInput(keyBindings, input)

		handler(window)
	}
}
