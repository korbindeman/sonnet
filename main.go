package main

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"os/signal"
	"syscall"
)

type InputHandler func(x, y int, width, height int) (int, int)

var keyBindings map[byte]InputHandler

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

func handleInput(input byte) InputHandler {
	if handler, exists := keyBindings[input]; exists {
		return handler
	}
	return func(x, y, width, height int) (int, int) {
		return x, y
	}
}

func addKeyBinding(key byte, handler InputHandler) {
	keyBindings[key] = handler
}

func readFileContent(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return lines, nil
}

func loadAndDisplayFile(filePath string, width, height int) {
	lines, err := readFileContent(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	clearScreen()
	moveCursor(1, 1)
	for i, line := range lines {
		if i >= height-1 {
			break
		}
		fmt.Print(line, "\r\n")
	}
	moveCursor(1, 1)
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

	// Initialize key bindings map
	keyBindings = make(map[byte]InputHandler)

	// Define key bindings
	addKeyBinding('q', func(x, y, width, height int) (int, int) {
		clearScreen()
		os.Exit(0)
		return x, y
	})
	addKeyBinding('h', func(x, y, width, height int) (int, int) {
		if x > 1 {
			x--
		}
		return x, y
	})
	addKeyBinding('j', func(x, y, width, height int) (int, int) {
		if y < height {
			y++
		}
		return x, y
	})
	addKeyBinding('k', func(x, y, width, height int) (int, int) {
		if y > 1 {
			y--
		}
		return x, y
	})
	addKeyBinding('l', func(x, y, width, height int) (int, int) {
		if x < width {
			x++
		}
		return x, y
	})
	addKeyBinding('r', func(x, y, width, height int) (int, int) {
		filePath := "example.txt"
		loadAndDisplayFile(filePath, width, height)
		return x, y
	})
	addKeyBinding(':', func(x, y, width, height int) (int, int) {
		moveCursor(height, 1)
		// clear the line
		fmt.Print("\x1b[K")
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
		moveCursor(y, x)
		loadAndDisplayFile(filename, width, height)
		return x, y
	})
	addKeyBinding('i', func(x, y, width, height int) (int, int) {
		// insert mode
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

	// _, height, err := getWindowSize()
	// for i := 0; i < height; i++ {
	// 	fmt.Print("\r\n", i+1)
	// }

	moveCursor(y, x)

	for {
		width, height, err := getWindowSize()
		if err != nil {
			fmt.Println(err)
			break
		}

		input, err := readInput()
		if err != nil {
			fmt.Println(err)
			break
		}

		handler := handleInput(input)
		x, y = handler(x, y, width, height)
		moveCursor(y, x)
	}
}
