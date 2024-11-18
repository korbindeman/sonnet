package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

func EnableRawMode() (*term.State, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	return oldState, nil
}

func DisableRawMode(oldState *term.State) {
	term.Restore(int(os.Stdin.Fd()), oldState)
}

func ClearScreen() {
	fmt.Print("\x1b[2J")
}

func MoveCursor(row, col int) {
	fmt.Printf("\x1b[%d;%dH", row, col)
}

func ReadInput() (byte, error) {
	buf := make([]byte, 1)
	_, err := os.Stdin.Read(buf)
	if err != nil {
		return 0, fmt.Errorf("error reading input: %w", err)
	}
	return buf[0], nil
}

func GetWindowSize() (int, int, error) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, 0, fmt.Errorf("error getting window size: %w", err)
	}
	return width, height, nil
}

func Initialize() {
	oldState, err := EnableRawMode()
	if err != nil {
		fmt.Println("Error enabling raw mode:", err)
		return
	}
	defer DisableRawMode(oldState)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		DisableRawMode(oldState)
		os.Exit(0)
	}()

	ClearScreen()
}
