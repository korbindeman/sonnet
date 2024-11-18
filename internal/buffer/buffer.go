package buffer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Buffer struct {
	rows       []string
	line       int
	col        int
	virtualCol int
}

func NewBuffer() *Buffer {
	return &Buffer{line: 1, col: 1, virtualCol: 1}
}

func (b *Buffer) Rows() []string {
	return b.rows
}

func (b *Buffer) InsertRow(row int, line string) {
	b.rows = append(b.rows[:row], append([]string{line}, b.rows[row:]...)...)
}

func (b *Buffer) DeleteRow(row int) {
	b.rows = append(b.rows[:row], b.rows[row+1:]...)
}

func (b *Buffer) UpdateRow(row int, line string) {
	b.rows[row] = line
}

func (b *Buffer) InsertChar(row, col int, ch rune) {
	b.rows[row] = b.rows[row][:col] + string(ch) + b.rows[row][col:]
}

func (b *Buffer) DeleteChar(row, col int) {
	b.rows[row] = b.rows[row][:col] + b.rows[row][col+1:]
}

func (b *Buffer) Length() int {
	return len(b.rows)
}

func (b *Buffer) LineLength(row int) int {
	return len(b.rows[row])
}

func (b *Buffer) CurrentLineLength() int {
	return len(b.rows[b.line-1])
}

func (b *Buffer) Replace(buffer *Buffer) {
	b.rows = buffer.rows
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

func LoadFile(filename string) (*Buffer, error) {
	buffer := NewBuffer()

	lines, err := readFileContent(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	spaces := "    "
	for i, line := range lines {
		output := strings.ReplaceAll(line, "\t", spaces)
		buffer.InsertRow(i, output)
	}

	return buffer, nil
}
