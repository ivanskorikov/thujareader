package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"thujareader/internal/reader"
	"thujareader/internal/ui"
)

func main() {
	var initialBook *reader.LoadedBook
	if len(os.Args) > 1 {
		unified := reader.NewDefaultUnifiedReader()
		book, err := unified.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		initialBook = &book
	}

	model := ui.NewModelWithInitialBook(initialBook)
	program := tea.NewProgram(model, tea.WithOutput(os.Stdout))

	if err := program.Start(); err != nil {
		log.Fatal(err)
	}
}
