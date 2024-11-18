package main

import (
	"fmt"

	"github.com/korbindeman/sonnet/internal/buffer"
	"github.com/korbindeman/sonnet/internal/keymaps"
	"github.com/korbindeman/sonnet/internal/utils"
	"github.com/korbindeman/sonnet/internal/window"
)

type InputHandler func(x, y int, width, height int) (int, int)

func handleInput(keyBindings keymaps.KeyBindings, input byte) keymaps.InputHandler {
	if handler, exists := keyBindings[input]; exists {
		return handler
	}
	return func(x, y, width, height int, buffer *buffer.Buffer) (int, int) { return x, y }
}

func main() {
	utils.EnableRawMode()

	keyBindings := keymaps.NewDefaultKeyBindings()

	window := window.NewWindow()

	buffer := buffer.NewBuffer()
	buffer.InsertRow(0, "Hello, world!")

	window.DisplayBuffer(buffer)

	window.SetCursor()

	x, y := window.Cursor.GetPosition()

	for {
		input, err := utils.ReadInput()
		if err != nil {
			fmt.Println(err)
			break
		}

		handler := handleInput(keyBindings, input)
		width, height := window.GetSize()
		x, y = handler(x, y, width, height, buffer)
		window.Cursor.Move(y, x)
	}
}
