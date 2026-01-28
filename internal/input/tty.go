package input

import (
	"bufio"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

type TTYInput struct {
	in  io.Reader
	out io.Writer
}

func NewTTYInput(in io.Reader, out io.Writer) *TTYInput {
	return &TTYInput{
		in:  in,
		out: out,
	}
}

func (t *TTYInput) ReadLine() (string, error) {
	reader := bufio.NewReader(t.in)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimRight(line, "\r\n"), err
}

func (t *TTYInput) GetChar() (rune, error) {
	fd := int(os.Stdin.Fd())

	// Sauvegarde de l'état du terminal
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return 0, err
	}
	defer term.Restore(fd, oldState)

	var buf [1]byte
	_, err = os.Stdin.Read(buf[:])
	if err != nil {
		return 0, err
	}

	return rune(buf[0]), nil
}
