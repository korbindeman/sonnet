package keymaps

import "github.com/korbindeman/sonnet/internal/buffer"

type KeymapData struct {
	X      int
	Y      int
	Width  int
	Height int
	Buffer *buffer.Buffer
}

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
		}
		return x, y
	})
	keyBindings.Add('k', func(x, y, width, height int, buffer *buffer.Buffer) (int, int) {
		if y > 1 {
			y--
		}
		return x, y
	})
	keyBindings.Add('l', func(x, y, width, height int, buffer *buffer.Buffer) (int, int) {
		if x < width {
			x++
		}
		return x, y
	})

	return keyBindings
}
