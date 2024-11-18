package main

import (
	"os"

	"github.com/korbindeman/sonnet/internal/utils"
	"golang.org/x/term"
)

type sonnet struct {
	oldState *term.State
}

func Initialize() *sonnet {
	s := &sonnet{}
	s.oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
	utils.ClearScreen()
	return s
}

func (s *sonnet) Close() {
	term.Restore(int(os.Stdin.Fd()), s.oldState)
	utils.ClearScreen()
	os.Exit(0)
}
