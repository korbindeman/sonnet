package keymaps

import (
	"fmt"
	"os"

	"github.com/korbindeman/sonnet/internal/buffer"
	"github.com/korbindeman/sonnet/internal/io"
	"github.com/korbindeman/sonnet/internal/render"
	"github.com/korbindeman/sonnet/internal/window"
)

type InputHandler func(*window.Window)

type KeyBindings map[byte]InputHandler

const (
	KeyEnter     = 13
	KeyEscape    = 27
	KeyBackspace = 127
)

func HandleInput(keyBindings KeyBindings, input byte) InputHandler {
	if handler, exists := keyBindings[input]; exists {
		return handler
	}
	return func(win *window.Window) {}
}

func NewKeyBindings() KeyBindings {
	return make(KeyBindings)
}

func (k KeyBindings) Add(key byte, handler InputHandler) {
	k[key] = handler
}

func NewDefaultKeyBindings() KeyBindings {
	keyBindings := NewKeyBindings()

	keyBindings.Add('h', func(win *window.Window) {
		win.Buffer.MoveLeft()
		win.SetCursor()
	})
	keyBindings.Add('j', func(win *window.Window) {
		win.Buffer.MoveDown()
		win.SetCursor()
	})
	keyBindings.Add('k', func(win *window.Window) {
		win.Buffer.MoveUp()
		win.SetCursor()
	})
	keyBindings.Add('l', func(win *window.Window) {
		win.Buffer.MoveRight()
		win.SetCursor()
	})

	keyBindings.Add('q', func(win *window.Window) {
		render.ClearScreen()
		render.MoveCursor(render.NewCoord(0, 0))
		os.Exit(0)
	})

	keyBindings.Add(':', func(win *window.Window) {
		_, height := win.GetSize()
		render.MoveCursor(render.NewCoord(height, 1))
		render.ClearLine()
		fmt.Print(":")
		filename := ""
		for {
			input, err := io.ReadInput()
			if err != nil {
				fmt.Println("Input error:", err)
				continue
			}
			switch input {
			case KeyEnter:
				if filename == "" {
					fmt.Println("No filename provided")
					return
				}
				newBuffer, err := buffer.LoadFile(filename)
				if err != nil {
					fmt.Printf("Error loading file: %v\n", err)
					return
				}
				win.LoadBuffer(newBuffer)
				return
			case KeyEscape:
				render.MoveCursor(render.NewCoord(height, 1))
				render.ClearLine()
				win.SetCursor()
				return
			case KeyBackspace:
				if len(filename) > 0 {
					filename = filename[:len(filename)-1]
					fmt.Print("\b \b")
				}
			default:
				filename += string(input)
				fmt.Print(string(input))
			}
		}
	})

	keyBindings.Add('i', func(win *window.Window) {
		// TODO: Implement insert mode functionality
	})

	return keyBindings
}
