package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/korbindeman/sonnet/internal/buffer"
	"github.com/korbindeman/sonnet/internal/keymaps"
	"golang.org/x/term"
)

type InputHandler func(x, y int, width, height int) (int, int)

func enableRawMode() (*term.State, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	return oldState, nil
}

func disableRawMode(oldState *term.State) {
	term.Restore(int(os.Stdin.Fd()), oldState)
}

func clearScreen() {
	fmt.Print("\x1b[2J")
}

func moveCursor(row, col int) {
	fmt.Printf("\x1b[%d;%dH", row, col)
}

func readInput() (byte, error) {
	buf := make([]byte, 1)
	_, err := os.Stdin.Read(buf)
	if err != nil {
		return 0, fmt.Errorf("error reading input: %w", err)
	}
	return buf[0], nil
}

func getWindowSize() (int, int, error) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, 0, fmt.Errorf("error getting window size: %w", err)
	}
	return width, height, nil
}

func handleInput(keyBindings keymaps.KeyBindings, input byte) keymaps.InputHandler {
	if handler, exists := keyBindings[input]; exists {
		return handler
	}
	return func(x, y, width, height int, buffer *buffer.Buffer) (int, int) { return x, y }
}

func displayBuffer(buffer *buffer.Buffer, width, height int) {
	for i, line := range buffer.Rows() {
		if i >= height-1 {
			break
		}
		fmt.Print(line, "\r\n")
	}
}

func main() {
	oldState, err := enableRawMode()
	if err != nil {
		fmt.Println("Error enabling raw mode:", err)
		return
	}
	defer disableRawMode(oldState)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		disableRawMode(oldState)
		os.Exit(0)
	}()

	clearScreen()

	x := 1
	y := 1

	keyBindings := keymaps.NewDefaultKeyBindings()

	keyBindings.Add('q', func(x, y, width, height int, buffer *buffer.Buffer) (int, int) {
		clearScreen()
		disableRawMode(oldState)
		moveCursor(0, 0)
		os.Exit(0)
		return x, y
	})
	keyBindings.Add(':', func(x, y, width, height int, curbuffer *buffer.Buffer) (int, int) {
		moveCursor(height, 1)
		fmt.Print("\x1b[K") // clear the line
		fmt.Print(":")
		filename := ""
		for {
			input, err := readInput()
			if err != nil {
				fmt.Println(err)
				return x, y
			}
			if input == 13 {
				break
			}
			if input == 27 {
				moveCursor(height, 1)
				fmt.Print("\x1b[K")
				moveCursor(y, x)
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
		clearScreen()
		moveCursor(y, x)
		newbuffer, _ := buffer.LoadFile(filename)
		curbuffer.Replace(newbuffer)
		displayBuffer(curbuffer, width, height)
		return x, y
	})
	keyBindings.Add('i', func(x, y, width, height int, buffer *buffer.Buffer) (int, int) {
		for {
			input, err := readInput()
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

	moveCursor(y, x)

	buffer := buffer.NewBuffer()
	buffer.InsertRow(0, "Hello, world!")
	width, height, err := getWindowSize()
	displayBuffer(buffer, width, height)

	moveCursor(y, x)

	for {
		width, _, err := getWindowSize()
		if err != nil {
			fmt.Println(err)
			break
		}

		input, err := readInput()
		if err != nil {
			fmt.Println(err)
			break
		}

		handler := handleInput(keyBindings, input)
		x, y = handler(x, y, width, height, buffer)
		moveCursor(y, x)
	}
}
