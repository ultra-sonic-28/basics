package apple2

import (
	"fmt"
)

type ConsoleRenderer struct{}

func NewConsoleRenderer() *ConsoleRenderer {
	return &ConsoleRenderer{}
}

func (r *ConsoleRenderer) RenderText(lines []string) {
	fmt.Print("\033[H\033[2J")
	for _, line := range lines {
		fmt.Println(line)
	}
}
