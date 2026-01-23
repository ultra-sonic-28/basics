package apple2

import (
	"fmt"
)

type ConsoleRenderer struct{}

func NewConsoleRenderer() *ConsoleRenderer {
	return &ConsoleRenderer{}
}

func (r *ConsoleRenderer) RenderText(lines []string) {
	// Clear screen + home
	fmt.Print("\033[H\033[2J")
	for _, line := range lines {
		fmt.Println(line)
	}
}

func (r *ConsoleRenderer) RenderChar(x, y int, ch rune) {
	// y,x sont 0-based, ANSI est 1-based
	fmt.Printf("\033[%d;%dH%c", y+1, x+1, ch)
}
