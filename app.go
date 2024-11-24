package main

import (
	"os"

	"github.com/korbindeman/sonnet/internal/render"
	"golang.org/x/term"
)

type sonnet struct {
	oldState *term.State
}

func Initialize() *sonnet {
	s := &sonnet{}
	s.oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
	render.ClearScreen()
	return s
}

func (s *sonnet) Close() {
	term.Restore(int(os.Stdin.Fd()), s.oldState)
	render.ClearScreen()
	os.Exit(0)
}
