package main

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"os/signal"
	"syscall"
)

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
		switch input {
		case 'q':
			clearScreen()
			os.Exit(0)
		case 'h':
			if x > 1 {
				x--
				moveCursor(y, x)
			}
		case 'j':
			if y < height {
				y++
				moveCursor(y, x)
			}
		case 'k':
			if y > 1 {
				y--
				moveCursor(y, x)
			}
		case 'l':
			if x < width {
				x++
				moveCursor(y, x)
			}
		default:
		}
	}
}
