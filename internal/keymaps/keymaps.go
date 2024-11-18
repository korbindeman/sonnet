package keymaps

import (
	"fmt"
	"math"
	"os"

	"github.com/korbindeman/sonnet/internal/buffer"
	"github.com/korbindeman/sonnet/internal/editor"
	"github.com/korbindeman/sonnet/internal/utils"
)

type InputHandler func(x, y, width, height int, buffer *buffer.Buffer) (int, int)

type KeyBindings map[byte]InputHandler

func NewKeyBindings() KeyBindings {
	return make(KeyBindings)
}

func (k KeyBindings) Add(key byte, handler InputHandler) {
	k[key] = handler
}

func NewDefaultKeyBindings() KeyBindings {
	keyBindings := NewKeyBindings()

	keyBindings.Add('h', func(x, y, width, height int, buffer *buffer.Buffer) (int, int) {
		if x > 1 {
			x--
		}
		return x, y
	})
	keyBindings.Add('j', func(x, y, width, height int, buffer *buffer.Buffer) (int, int) {
		if y < buffer.Length() {
			y++
			if buffer.LineLength(y-1) < x {
				x = int(math.Max(float64(buffer.LineLength(y-1)), 1))
			}
		}
		return x, y
	})
	keyBindings.Add('k', func(x, y, width, height int, buffer *buffer.Buffer) (int, int) {
		if y > 1 {
			y--
			if buffer.LineLength(y-1) < x {
				x = int(math.Max(float64(buffer.LineLength(y-1)), 1))
			}
		}
		return x, y
	})
	keyBindings.Add('l', func(x, y, width, height int, buffer *buffer.Buffer) (int, int) {
		if x < buffer.LineLength(y-1) {
			x++
		}
		return x, y
	})

	keyBindings.Add('q', func(x, y, width, height int, buffer *buffer.Buffer) (int, int) {
		utils.ClearScreen()
		utils.MoveCursor(0, 0)
		os.Exit(0)
		return x, y
	})
	keyBindings.Add(':', func(x, y, width, height int, curbuffer *buffer.Buffer) (int, int) {
		utils.MoveCursor(height, 1)
		fmt.Print("\x1b[K") // clear the line
		fmt.Print(":")
		filename := ""
		for {
			input, err := utils.ReadInput()
			if err != nil {
				fmt.Println(err)
				return x, y
			}
			if input == 13 {
				break
			}
			if input == 27 {
				utils.MoveCursor(height, 1)
				fmt.Print("\x1b[K")
				utils.MoveCursor(y, x)
				return x, y
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
		x, y = 1, 1
		utils.ClearScreen()
		utils.MoveCursor(y, x)
		newbuffer, _ := buffer.LoadFile(filename)
		curbuffer.Replace(newbuffer)
		editor.DisplayBuffer(curbuffer, width, height)
		return x, y
	})
	keyBindings.Add('i', func(x, y, width, height int, buffer *buffer.Buffer) (int, int) {
		for {
			input, err := utils.ReadInput()
			if err != nil {
				fmt.Println(err)
			}
			if input == 27 {
				break
			}
			if input == 127 {
				fmt.Print("\b \b")
				continue
			}

			fmt.Print(string(input))
		}
		return x, y
	})

	return keyBindings
}
