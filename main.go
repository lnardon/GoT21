package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Print("\033[H\033[2J") // Clean terminal

	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running game:", err)
	}
}
