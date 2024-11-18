package editor

import (
	"fmt"
	"github.com/korbindeman/sonnet/internal/buffer"
)

func DisplayBuffer(buffer *buffer.Buffer, width, height int) {
	for i, line := range buffer.Rows() {
		if i >= height-1 {
			break
		}
		fmt.Print(line, "\r\n")
	}
}
