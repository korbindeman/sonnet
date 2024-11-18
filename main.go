package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/korbindeman/sonnet/internal/buffer"
	"github.com/korbindeman/sonnet/internal/editor"
	"github.com/korbindeman/sonnet/internal/keymaps"
	"github.com/korbindeman/sonnet/internal/utils"
)

type InputHandler func(x, y int, width, height int) (int, int)

func handleInput(keyBindings keymaps.KeyBindings, input byte) keymaps.InputHandler {
	if handler, exists := keyBindings[input]; exists {
		return handler
	}
	return func(x, y, width, height int, buffer *buffer.Buffer) (int, int) { return x, y }
}

func main() {
	oldState, err := utils.EnableRawMode()
	if err != nil {
		fmt.Println("Error enabling raw mode:", err)
		return
	}
	defer utils.DisableRawMode(oldState)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		utils.DisableRawMode(oldState)
		os.Exit(0)
	}()

	utils.ClearScreen()

	x := 1
	y := 1

	keyBindings := keymaps.NewDefaultKeyBindings()

	utils.MoveCursor(y, x)

	buffer := buffer.NewBuffer()
	buffer.InsertRow(0, "Hello, world!")
	width, height, err := utils.GetWindowSize()
	editor.DisplayBuffer(buffer, width, height)

	utils.MoveCursor(y, x)

	for {
		width, _, err := utils.GetWindowSize()
		if err != nil {
			fmt.Println(err)
			break
		}

		input, err := utils.ReadInput()
		if err != nil {
			fmt.Println(err)
			break
		}

		handler := handleInput(keyBindings, input)
		x, y = handler(x, y, width, height, buffer)
		utils.MoveCursor(y, x)
	}
}
