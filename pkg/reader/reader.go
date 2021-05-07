package reader

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type LineReader struct {
	b *bufio.Reader
}

func New(rd io.Reader) *LineReader {
	return &LineReader{
		b: bufio.NewReader(rd),
	}
}

func (lr *LineReader) MustString() string {
	str, err := lr.b.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(str, "\n")
}

func (lr *LineReader) MustInt() int {
	str, err := lr.b.ReadString('\n')
	if err != nil {
		return 0
	}
	str = strings.TrimSuffix(str, "\n")
	n, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return n
}

func (lr *LineReader) MustBool() bool {
	str, err := lr.b.ReadString('\n')
	if err != nil {
		return false
	}
	str = strings.TrimSuffix(str, "\n")
	n, err := strconv.Atoi(str)
	if err != nil {
		return false
	}
	return n != 0
}
