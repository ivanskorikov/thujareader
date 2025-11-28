package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"thujareader/internal/ui"
)

func main() {
	model := ui.NewModel()
	program := tea.NewProgram(model, tea.WithOutput(os.Stdout))

	if err := program.Start(); err != nil {
		log.Fatal(err)
	}
}
