package keymaps

import (
	"fmt"
	"os"

	"github.com/korbindeman/sonnet/internal/buffer"
	"github.com/korbindeman/sonnet/internal/utils"
	"github.com/korbindeman/sonnet/internal/window"
)

type InputHandler func(*window.Window)

type KeyBindings map[byte]InputHandler

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
		utils.ClearScreen()
		utils.MoveCursor(0, 0)
		os.Exit(0)
	})

	keyBindings.Add(':', func(win *window.Window) {
		_, height := win.GetSize()
		utils.MoveCursor(height, 1)
		fmt.Print("\x1b[K") // clear the line
		fmt.Print(":")
		filename := ""
		for {
			input, err := utils.ReadInput()
			if err != nil {
				fmt.Println(err)
			}
			if input == 13 {
				break
			}
			if input == 27 {
				utils.MoveCursor(height, 1)
				fmt.Print("\x1b[K")
				win.SetCursor()
			}
			if input == 127 {
				if len(filename) == 0 {
					continue
				}
				filename = filename[:len(filename)-1]
				fmt.Print("\b \b")
				continue
			}
			filename += string(input)
			fmt.Print(string(input))
		}
		newbuffer, _ := buffer.LoadFile(filename)
		win.LoadBuffer(newbuffer)
	})

	keyBindings.Add('i', func(win *window.Window) {
	})

	return keyBindings
}
