package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"thujareader/internal/config"
	"thujareader/internal/reader"
	"thujareader/internal/state"
	"thujareader/internal/ui"
)

func main() {
	// Resolve configuration and state file paths.
	paths, err := config.DefaultPaths()
	if err != nil {
		log.Fatal(err)
	}

	// Load configuration; on error, fall back to defaults but continue.
	cfg, err := config.Load(paths.ConfigFile)
	if err != nil {
		log.Printf("warning: failed to load config: %v", err)
	}

	// Load persisted application state (bookmarks, positions, recent files).
	store := state.NewFileStore(paths.StateFile)
	appState, err := store.Load()
	if err != nil {
		log.Printf("warning: failed to load state: %v", err)
	}

	var initialBook *reader.LoadedBook
	if len(os.Args) > 1 {
		unified := reader.NewDefaultUnifiedReader()
		book, err := unified.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		initialBook = &book
	}

	// Adapt stored bookmarks (keyed by string) to the UI's map keyed by
	// reader.BookID.
	loadedBookmarks := make(map[reader.BookID][]reader.Bookmark)
	for k, v := range appState.Bookmarks {
		loadedBookmarks[reader.BookID(k)] = v
	}

	model := ui.NewModelWithInitialBookAndBookmarks(initialBook, loadedBookmarks)
	// Apply configuration options that the UI currently understands.
	if cfg.RecentListSize > 0 {
		model.SetRecentLimit(cfg.RecentListSize)
	}

	program := tea.NewProgram(model, tea.WithOutput(os.Stdout))

	finalModel, err := program.Run()
	if err != nil {
		log.Fatal(err)
	}

	// On normal exit, persist updated bookmarks (and leave room for
	// positions/recent files when those are wired in).
	if m, ok := finalModel.(ui.Model); ok {
		bookmarks := m.ExportBookmarks()
		appState.Bookmarks = make(map[string][]reader.Bookmark)
		for k, v := range bookmarks {
			appState.Bookmarks[string(k)] = v
		}
		if err := store.Save(appState); err != nil {
			log.Printf("warning: failed to save state: %v", err)
		}
	}
}
