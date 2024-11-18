package main

import (
	"fmt"

	"github.com/korbindeman/sonnet/internal/buffer"
	"github.com/korbindeman/sonnet/internal/keymaps"
	"github.com/korbindeman/sonnet/internal/utils"
	"github.com/korbindeman/sonnet/internal/window"
)

func handleInput(keyBindings keymaps.KeyBindings, input byte) keymaps.InputHandler {
	if handler, exists := keyBindings[input]; exists {
		return handler
	}
	return func(win *window.Window) {}
}

func main() {
	sonnet := Initialize()
	defer sonnet.Close()

	keyBindings := keymaps.NewDefaultKeyBindings()

	buf, _ := buffer.LoadFile("main.go")

	window := window.NewWindow()
	window.LoadBuffer(buf)

	for {
		input, err := utils.ReadInput()
		if err != nil {
			fmt.Println(err)
			break
		}

		handler := handleInput(keyBindings, input)

		handler(window)
	}
}
