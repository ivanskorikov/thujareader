package ui

import (
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
	"golang.org/x/term"

	"thujareader/internal/reader"
)

// menuID identifies a top-level menu.
type menuID int

const (
	menuFile menuID = iota
	menuSearch
	menuView
	menuBookmarks
	menuHelp
)

// commandID represents a high-level command invoked from menus or keybindings.
type commandID int

const (
	cmdNone commandID = iota
	cmdOpen
	cmdExit
	cmdFind
	cmdToc
	cmdBookmarks
	cmdRecentFiles
	cmdHelp
	cmdAddBookmark
	cmdDeleteBookmark
)

// menuItem is a single item within a menu.
type menuItem struct {
	label   string
	command commandID
}

// menu describes a top-level menu.
type menu struct {
	id    menuID
	label string
	items []menuItem
}

// Model holds UI state for the TUI shell emulating DOS edit.exe.
type Model struct {
	width  int
	height int

	theme Theme

	// unifiedReader is the shared entry point for loading books from
	// disk. It is used both for CLI-argument opens and the in-app
	// File → Open flow.
	unifiedReader reader.UnifiedReader

	// currentBook holds the last successfully loaded book, along with a
	// pre-split view of its text for simple line-based rendering.
	currentBook *reader.LoadedBook
	// textRunes caches the book text as runes so that wrapping,
	// navigation, and search can operate on rune offsets (matching the
	// domain model's TotalCharacters semantics).
	textRunes []rune
	// lines holds the wrapped visual lines for the current viewport
	// width; lineOffsets maps each visual line to its starting rune
	// offset within the book's linear text.
	lines       []string
	lineOffsets []int
	topLine     int

	// currentPos tracks the logical position within the book. It is
	// updated when scrolling or jumping so that status/location display
	// can reflect the current chapter and percentage.
	currentPos reader.Position

	// TOC dialog state.
	tocOpen  bool
	tocIndex int

	// Bookmarks dialog state and in-memory storage.
	bookmarks     map[reader.BookID][]reader.Bookmark
	bookmarksOpen bool
	bookmarkIndex int

	// Recent files list and dialog state.
	recentFiles []string
	recentOpen  bool
	recentIndex int
	recentLimit int

	// Search state for Find / Find Next.
	lastSearch       string
	lastSearchOffset int // rune offset of last match start; -1 if none

	menus       []menu
	activeMenu  int  // index into menus, -1 when no menu is active
	activeItem  int  // index into items of the active menu
	menuOpen    bool // whether menu bar interaction is active
	statusLine  string
	statusDirty bool

	// inputMode indicates that the UI is currently collecting a single
	// line of text input from the user (e.g. for a file path).
	inputMode   bool
	inputPrompt string
	inputBuffer []rune
	// pendingCommand records which command should be executed when the
	// current line input is confirmed (e.g. cmdOpen).
	pendingCommand commandID
}

// NewModel constructs the initial UI model without a pre-loaded book.
func NewModel() Model {
	return NewModelWithInitialBook(nil)
}

// NewModelWithInitialBook constructs the initial UI model, optionally
// pre-populated with a book that was opened via CLI arguments.
func NewModelWithInitialBook(book *reader.LoadedBook) Model {
	m := Model{
		// Start with a reasonable default size so that the UI can render
		// even if no WindowSizeMsg is delivered (which can happen on some
		// terminals, especially on Windows). Resize events will override
		// these values when they arrive.
		width:         80,
		height:        25,
		theme:         ThemeFromEnv(),
		unifiedReader: reader.NewDefaultUnifiedReader(),
		menus: []menu{
			{
				id:    menuFile,
				label: "File",
				items: []menuItem{
					{label: "Open...  F3", command: cmdOpen},
					{label: "Recent Files", command: cmdRecentFiles},
					{label: "Exit      Alt+F X", command: cmdExit},
				},
			},
			{
				id:    menuSearch,
				label: "Search",
				items: []menuItem{
					{label: "Find...  F7", command: cmdFind},
					{label: "TOC", command: cmdToc},
				},
			},
			{
				id:    menuView,
				label: "View",
				items: []menuItem{},
			},
			{
				id:    menuBookmarks,
				label: "Bookmarks",
				items: []menuItem{
					{label: "Manage Bookmarks", command: cmdBookmarks},
					{label: "Add Bookmark  F2", command: cmdAddBookmark},
					{label: "Delete Bookmark", command: cmdDeleteBookmark},
				},
			},
			{
				id:    menuHelp,
				label: "Help",
				items: []menuItem{
					{label: "Help Topics  F1", command: cmdHelp},
				},
			},
		},
		activeMenu:  -1,
		activeItem:  0,
		statusLine:  "Press F10 or Alt key combinations to open menus. F1 for Help.",
		bookmarks:   make(map[reader.BookID][]reader.Bookmark),
		recentLimit: 10,
	}

	// Try to detect the actual terminal size at startup so that initial
	// wrapping uses the full window width/height even on platforms where
	// Bubble Tea may not immediately deliver a WindowSizeMsg.
	if w, h, ok := detectTerminalSize(); ok {
		if w > 0 {
			m.width = w
		}
		if h > 0 {
			m.height = h
		}
	}

	if book != nil {
		m.setBook(*book)
	}

	return m
}

// NewModelWithInitialBookAndBookmarks constructs the initial UI model
// and pre-populates it with a book (if any) and a set of bookmarks
// loaded from persisted state.
func NewModelWithInitialBookAndBookmarks(book *reader.LoadedBook, bookmarks map[reader.BookID][]reader.Bookmark) Model {
	m := NewModelWithInitialBook(book)
	if bookmarks != nil {
		m.bookmarks = bookmarks
	}
	return m
}

// detectTerminalSize returns the current terminal width and height in
// cells, if stdout is attached to a TTY and the size can be queried.
// It is a best-effort helper used to initialize the model before any
// WindowSizeMsg is received.
func detectTerminalSize() (int, int, bool) {
	fd := int(os.Stdout.Fd())
	if !term.IsTerminal(fd) {
		return 0, 0, false
	}
	w, h, err := term.GetSize(fd)
	if err != nil || w <= 0 || h <= 0 {
		return 0, 0, false
	}
	return w, h, true
}

// Init runs any startup commands. For now there are no asynchronous
// startup operations required.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages including resize events and
// keyboard input.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Recompute wrapping when the window size changes so that text
		// fits the new viewport width.
		m.reflowWrappedLines()
		return m, nil

	case tea.KeyMsg:
		// Always allow Ctrl+C to quit.
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		// When we are in a simple input mode (e.g. entering a file path
		// for the Open command), route all key presses through the input
		// handler instead of the normal menu/keybinding logic.
		if m.inputMode {
			if m.handleInputKey(msg) {
				return m, nil
			}
			return m, nil
		}

		if m.handleKey(msg) {
			return m, nil
		}
	}

	return m, nil
}

func (m *Model) openMenuByAltKey(ch rune) {
	for i, menu := range m.menus {
		if len(menu.label) == 0 {
			continue
		}
		if strings.ToLower(string(menu.label[0])) == strings.ToLower(string(ch)) {
			m.menuOpen = true
			m.activeMenu = i
			m.activeItem = 0
			return
		}
	}
}

func (m *Model) handleKey(msg tea.KeyMsg) bool {
	switch msg.Type {
	case tea.KeyF10:
		// Toggle menu bar interaction.
		if m.menuOpen {
			m.menuOpen = false
			m.activeMenu = -1
		} else {
			m.menuOpen = true
			if m.activeMenu < 0 {
				m.activeMenu = 0
			}
		}
		return true
	case tea.KeyF1:
		m.executeCommand(cmdHelp)
		return true
	case tea.KeyF2:
		// F2: add bookmark at current position.
		m.executeCommand(cmdAddBookmark)
		return true
	case tea.KeyF3:
		m.executeCommand(cmdOpen)
		return true
	case tea.KeyF7:
		// F7 either opens the Find dialog or, if a previous search term
		// exists, jumps to the next match.
		if !m.inputMode && m.lastSearch != "" {
			m.performSearch(m.lastSearch, false)
		} else {
			m.executeCommand(cmdFind)
		}
		return true
	}

	// Alt+<letter> opens corresponding menu (e.g., Alt+F for File).
	if msg.Alt && len(msg.Runes) == 1 {
		m.openMenuByAltKey(msg.Runes[0])
		return true
	}

	if !m.menuOpen {
		// When the menu is not open, either handle TOC navigation when
		// the TOC dialog is active or perform normal reading/view
		// navigation.
		if m.currentBook == nil {
			return false
		}

		// TOC dialog navigation when open.
		if m.tocOpen {
			switch msg.Type {
			case tea.KeyEsc:
				m.tocOpen = false
				return true
			case tea.KeyUp:
				if m.tocIndex > 0 {
					m.tocIndex--
				}
				return true
			case tea.KeyDown:
				if m.currentBook != nil {
					maxIdx := len(m.currentBook.TOC) - 1
					if maxIdx >= 0 && m.tocIndex < maxIdx {
						m.tocIndex++
					}
				}
				return true
			case tea.KeyEnter:
				if m.currentBook != nil && m.tocIndex >= 0 && m.tocIndex < len(m.currentBook.TOC) {
					entry := m.currentBook.TOC[m.tocIndex]
					m.jumpToPosition(entry.Pos)
				}
				m.tocOpen = false
				return true
			}
			return false
		}

		// Bookmarks dialog navigation when open.
		if m.bookmarksOpen {
			switch msg.Type {
			case tea.KeyEsc:
				m.bookmarksOpen = false
				return true
			case tea.KeyUp:
				if m.bookmarkIndex > 0 {
					m.bookmarkIndex--
				}
				return true
			case tea.KeyDown:
				current := m.currentBookmarks()
				if len(current) == 0 {
					return true
				}
				if m.bookmarkIndex < len(current)-1 {
					m.bookmarkIndex++
				}
				return true
			case tea.KeyEnter:
				current := m.currentBookmarks()
				if len(current) == 0 {
					m.bookmarksOpen = false
					return true
				}
				if m.bookmarkIndex < 0 || m.bookmarkIndex >= len(current) {
					m.bookmarksOpen = false
					return true
				}
				bm := current[m.bookmarkIndex]
				m.jumpToPosition(bm.Pos)
				m.bookmarksOpen = false
				m.setStatus("Jumped to bookmark: " + bm.Name)
				return true
			}
			return false
		}

		// Recent files dialog navigation when open.
		if m.recentOpen {
			switch msg.Type {
			case tea.KeyEsc:
				m.recentOpen = false
				return true
			case tea.KeyUp:
				if m.recentIndex > 0 {
					m.recentIndex--
				}
				return true
			case tea.KeyDown:
				if len(m.recentFiles) == 0 {
					return true
				}
				if m.recentIndex < len(m.recentFiles)-1 {
					m.recentIndex++
				}
				return true
			case tea.KeyEnter:
				if len(m.recentFiles) == 0 {
					m.recentOpen = false
					return true
				}
				if m.recentIndex < 0 || m.recentIndex >= len(m.recentFiles) {
					m.recentOpen = false
					return true
				}
				path := m.recentFiles[m.recentIndex]
				m.recentOpen = false
				m.openPath(path)
				return true
			}
			return false
		}

		// Normal reading navigation when no modal dialog (like TOC) is
		// active.
		switch msg.Type {
		case tea.KeyUp:
			if m.topLine > 0 {
				m.topLine--
				m.updateCurrentPositionFromTopLine()
			}
			return true
		case tea.KeyDown:
			if m.topLine < len(m.lines)-1 {
				m.topLine++
				m.updateCurrentPositionFromTopLine()
			}
			return true
		case tea.KeyPgUp:
			page := m.visibleLineCount()
			if page <= 0 {
				page = 1
			}
			if m.topLine > 0 {
				m.topLine -= page
				if m.topLine < 0 {
					m.topLine = 0
				}
				m.updateCurrentPositionFromTopLine()
			}
			return true
		case tea.KeyPgDown:
			page := m.visibleLineCount()
			if page <= 0 {
				page = 1
			}
			maxTop := max(0, len(m.lines)-1)
			if m.topLine < maxTop {
				m.topLine += page
				if m.topLine > maxTop {
					m.topLine = maxTop
				}
				m.updateCurrentPositionFromTopLine()
			}
			return true
		case tea.KeyHome:
			if m.topLine != 0 {
				m.topLine = 0
				m.updateCurrentPositionFromTopLine()
			}
			return true
		case tea.KeyEnd:
			maxTop := max(0, len(m.lines)-1)
			if m.topLine != maxTop {
				m.topLine = maxTop
				m.updateCurrentPositionFromTopLine()
			}
			return true
		}
		return false
	}

	switch msg.Type {
	case tea.KeyEsc:
		m.menuOpen = false
		m.activeMenu = -1
		return true
	case tea.KeyLeft:
		if m.activeMenu > 0 {
			m.activeMenu--
			m.activeItem = 0
		}
		return true
	case tea.KeyRight:
		if m.activeMenu < len(m.menus)-1 {
			m.activeMenu++
			m.activeItem = 0
		}
		return true
	case tea.KeyUp:
		if m.activeMenu >= 0 {
			if m.activeItem > 0 {
				m.activeItem--
			}
		}
		return true
	case tea.KeyDown:
		if m.activeMenu >= 0 {
			items := m.menus[m.activeMenu].items
			if len(items) > 0 && m.activeItem < len(items)-1 {
				m.activeItem++
			}
		}
		return true
	case tea.KeyEnter:
		if m.activeMenu >= 0 {
			items := m.menus[m.activeMenu].items
			if len(items) > 0 && m.activeItem >= 0 && m.activeItem < len(items) {
				cmd := items[m.activeItem].command
				m.executeCommand(cmd)
			}
		}
		return true
	}

	return false
}

func (m *Model) executeCommand(cmd commandID) {
	switch cmd {
	case cmdOpen:
		// When invoking the Open command from the menu, close the menu so
		// that, after opening a file, the main area can display the book
		// contents instead of leaving the dropdown visible.
		m.menuOpen = false
		m.activeMenu = -1

		// Enter a simple line-input mode where the user can type a file
		// path to open. This is a minimal stand-in for a full file
		// dialog and is sufficient for Phase 3.
		m.inputMode = true
		m.inputPrompt = "Open file: "
		m.inputBuffer = m.inputBuffer[:0]
		m.pendingCommand = cmdOpen
		m.setStatus("Enter path to EPUB/FB2 file and press Enter.")
	case cmdExit:
		m.setStatus("Exit: press Alt+F then X or Ctrl+C to quit.")
	case cmdFind:
		// Enter search input mode. Reuse the simple one-line input UI
		// but distinguish via pendingCommand.
		m.inputMode = true
		m.inputPrompt = "Find: "
		m.inputBuffer = m.inputBuffer[:0]
		m.pendingCommand = cmdFind
		m.setStatus("Enter search text and press Enter. Press Esc to cancel.")
	case cmdToc:
		if m.currentBook == nil || len(m.currentBook.TOC) == 0 {
			m.setStatus("TOC: no table of contents available for this book.")
			return
		}
		// Open TOC dialog starting at first entry.
		m.tocOpen = true
		m.tocIndex = 0
		m.menuOpen = false
		m.activeMenu = -1
		m.setStatus("TOC: Use ↑/↓ to select, Enter to jump, Esc to cancel.")
	case cmdBookmarks:
		if m.currentBook == nil {
			m.setStatus("Bookmarks: no book is currently open.")
			return
		}
		current := m.currentBookmarks()
		if len(current) == 0 {
			m.setStatus("Bookmarks: no bookmarks for this book.")
			return
		}
		m.bookmarksOpen = true
		m.bookmarkIndex = 0
		m.menuOpen = false
		m.activeMenu = -1
		m.setStatus("Bookmarks: Use ↑/↓ to select, Enter to jump, Esc to cancel.")
	case cmdAddBookmark:
		if m.currentBook == nil {
			m.setStatus("Cannot add bookmark: no book is open.")
			return
		}
		name := "Bookmark " + itoa(len(m.currentBookmarks())+1)
		bm := reader.Bookmark{
			Name:   name,
			BookID: m.currentBook.Book.ID,
			Pos:    m.currentPos,
		}
		list := m.currentBookmarks()
		list = append(list, bm)
		if m.bookmarks == nil {
			m.bookmarks = make(map[reader.BookID][]reader.Bookmark)
		}
		m.bookmarks[m.currentBook.Book.ID] = list
		m.setStatus("Added bookmark: " + name)
	case cmdDeleteBookmark:
		if !m.bookmarksOpen || m.currentBook == nil {
			return
		}
		current := m.currentBookmarks()
		if len(current) == 0 || m.bookmarkIndex < 0 || m.bookmarkIndex >= len(current) {
			return
		}
		name := current[m.bookmarkIndex].Name
		current = append(current[:m.bookmarkIndex], current[m.bookmarkIndex+1:]...)
		m.bookmarks[m.currentBook.Book.ID] = current
		if m.bookmarkIndex >= len(current) && m.bookmarkIndex > 0 {
			m.bookmarkIndex--
		}
		m.setStatus("Deleted bookmark: " + name)
	case cmdRecentFiles:
		if len(m.recentFiles) == 0 {
			m.setStatus("Recent files: list is empty.")
			return
		}
		m.recentOpen = true
		m.recentIndex = 0
		m.menuOpen = false
		m.activeMenu = -1
		m.setStatus("Recent files: Use ↑/↓ to select, Enter to open, Esc to cancel.")
	case cmdHelp:
		m.setStatus("Help: not yet implemented (help screen will appear in later phase).")
	default:
		return
	}
}

// currentBookmarks returns the slice of bookmarks for the currently
// open book. It never returns nil; when no book is open or there are no
// bookmarks for the book it returns an empty slice.
func (m *Model) currentBookmarks() []reader.Bookmark {
	if m.currentBook == nil {
		return nil
	}
	if m.bookmarks == nil {
		return nil
	}
	list, ok := m.bookmarks[m.currentBook.Book.ID]
	if !ok {
		return nil
	}
	return list
}

func (m *Model) setStatus(text string) {
	m.statusLine = text
	m.statusDirty = true
}

// SetRecentLimit updates the maximum number of recent files remembered
// in memory. Non-positive values are ignored.
func (m *Model) SetRecentLimit(limit int) {
	if limit <= 0 {
		return
	}
	m.recentLimit = limit
}

// ExportBookmarks returns a copy of the in-memory bookmarks map so that
// callers (e.g. main) can persist it to disk without mutating internal
// state.
func (m Model) ExportBookmarks() map[reader.BookID][]reader.Bookmark {
	if m.bookmarks == nil {
		return map[reader.BookID][]reader.Bookmark{}
	}
	out := make(map[reader.BookID][]reader.Bookmark, len(m.bookmarks))
	for k, v := range m.bookmarks {
		copySlice := make([]reader.Bookmark, len(v))
		copy(copySlice, v)
		out[k] = copySlice
	}
	return out
}

// setBook installs a newly loaded book into the model and prepares a
// wrapped view over its text based on the current viewport width.
func (m *Model) setBook(book reader.LoadedBook) {
	m.currentBook = &book
	m.textRunes = []rune(book.Text)
	m.topLine = 0
	m.currentPos = reader.Position{ChapterIndex: 0, OffsetInChapter: 0}
	m.lastSearch = ""
	m.lastSearchOffset = -1
	m.tocIndex = 0
	m.reflowWrappedLines()
	m.updateCurrentPositionFromTopLine()
}

// openPath attempts to load the given file via the unified reader and
// update the UI state accordingly.
func (m *Model) openPath(path string) {
	path = strings.TrimSpace(path)
	if path == "" {
		m.setStatus("No file path provided.")
		return
	}

	book, err := m.unifiedReader.Open(path)
	if err != nil {
		m.setStatus("Failed to open: " + err.Error())
		return
	}

	m.setBook(book)
	m.setStatus("Opened: " + book.Book.Title)
}

// handleInputKey processes key presses while the model is in a simple
// line-input mode (used for the Open command in Phase 3).
func (m *Model) handleInputKey(msg tea.KeyMsg) bool {
	switch msg.Type {
	case tea.KeyEsc:
		m.inputMode = false
		m.inputBuffer = nil
		m.pendingCommand = cmdNone
		return true
	case tea.KeyEnter:
		input := strings.TrimSpace(string(m.inputBuffer))
		pending := m.pendingCommand
		m.inputMode = false
		m.inputBuffer = nil
		m.pendingCommand = cmdNone

		if pending == cmdOpen {
			m.openPath(input)
		} else if pending == cmdFind {
			m.performSearch(input, true)
		}
		return true
	case tea.KeyBackspace:
		if len(m.inputBuffer) > 0 {
			m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
		}
		return true
	default:
		if len(msg.Runes) > 0 {
			m.inputBuffer = append(m.inputBuffer, msg.Runes...)
			return true
		}
	}

	return false
}

// performSearch executes a simple forward substring search over the
// book text. When newTerm is true, the previous search state is
// reset; otherwise, the search continues from the last match
// position. On success it jumps the viewport to the found position;
// on failure it updates the status bar with an explanatory message.
func (m *Model) performSearch(term string, newTerm bool) {
	if m.currentBook == nil || len(term) == 0 {
		m.setStatus("Find: empty search term.")
		return
	}

	text := string(m.textRunes)
	if newTerm || term != m.lastSearch {
		m.lastSearch = term
		m.lastSearchOffset = -1
	}

	start := m.lastSearchOffset + 1
	if start < 0 {
		start = 0
	}
	if start >= len(text) {
		m.setStatus("Find: no more matches.")
		return
	}

	idx := strings.Index(text[start:], term)
	if idx == -1 {
		if m.lastSearchOffset == -1 {
			m.setStatus("Find: no matches.")
		} else {
			m.setStatus("Find: no more matches.")
		}
		return
	}

	matchOffset := start + idx
	m.lastSearchOffset = matchOffset
	pos := m.absoluteOffsetToPosition(matchOffset)
	m.jumpToPosition(pos)
	m.setStatus("Find: match found.")
}

// reflowWrappedLines recomputes wrapped lines and their rune offsets
// based on the current window width.
func (m *Model) reflowWrappedLines() {
	if m.currentBook == nil || len(m.textRunes) == 0 {
		m.lines = nil
		m.lineOffsets = nil
		m.topLine = 0
		return
	}

	innerWidth := max(0, m.width-2)
	if innerWidth <= 0 {
		m.lines = nil
		m.lineOffsets = nil
		m.topLine = 0
		return
	}

	lines := make([]string, 0, len(m.textRunes)/innerWidth+1)
	offsets := make([]int, 0, len(lines))

	var (
		lineRunes       []rune
		col             int // display width in cells
		lineStartOffset int
	)

	flushLine := func() {
		lines = append(lines, string(lineRunes))
		offsets = append(offsets, lineStartOffset)
		lineRunes = lineRunes[:0]
		col = 0
		lineStartOffset = 0
	}

	currentOffset := 0
	for _, r := range m.textRunes {
		if r == '\n' {
			// End current visual line on explicit newline.
			flushLine()
			currentOffset++
			lineStartOffset = currentOffset
			continue
		}

		rw := runewidth.RuneWidth(r)
		if rw <= 0 {
			rw = 1
		}

		// If adding this rune would exceed the inner width, flush the
		// current line and start a new one at this rune offset.
		if col > 0 && col+rw > innerWidth {
			flushLine()
			lineStartOffset = currentOffset
		}

		lineRunes = append(lineRunes, r)
		col += rw
		currentOffset++
	}

	// Flush any remaining runes as the last line.
	if len(lineRunes) > 0 {
		flushLine()
	}

	m.lines = lines
	m.lineOffsets = offsets
	if m.topLine >= len(m.lines) {
		m.topLine = max(0, len(m.lines)-1)
	}
}

// visibleLineCount returns how many text lines fit inside the bordered
// main area.
func (m Model) visibleLineCount() int {
	innerHeight := m.height - 3
	if innerHeight < 1 {
		innerHeight = 1
	}
	// One line is used by the bottom border; the remaining lines are
	// available for content.
	return max(0, innerHeight-1)
}

// updateCurrentPositionFromTopLine updates the logical Position based
// on the current topLine and lineOffsets mapping.
func (m *Model) updateCurrentPositionFromTopLine() {
	if m.currentBook == nil || len(m.lineOffsets) == 0 {
		m.currentPos = reader.Position{}
		return
	}
	idx := m.topLine
	if idx < 0 {
		idx = 0
	}
	if idx >= len(m.lineOffsets) {
		idx = len(m.lineOffsets) - 1
	}
	abs := m.lineOffsets[idx]
	m.currentPos = m.absoluteOffsetToPosition(abs)
}

// jumpToPosition moves the viewport so that the given logical
// Position becomes visible near the top of the screen.
func (m *Model) jumpToPosition(pos reader.Position) {
	if m.currentBook == nil || len(m.lineOffsets) == 0 {
		return
	}
	abs := m.positionToAbsoluteOffset(pos)
	// Find the first visual line whose starting offset is at or after
	// the target offset.
	line := 0
	for i, off := range m.lineOffsets {
		if off >= abs {
			line = i
			break
		}
	}
	m.topLine = line
	m.updateCurrentPositionFromTopLine()
}

// positionToAbsoluteOffset converts a logical Position into a rune
// offset within the book's linear text stream.
func (m Model) positionToAbsoluteOffset(pos reader.Position) int {
	if m.currentBook == nil {
		return 0
	}
	if pos.ChapterIndex < 0 || pos.ChapterIndex >= len(m.currentBook.Book.Chapters) {
		return 0
	}
	ch := m.currentBook.Book.Chapters[pos.ChapterIndex]
	return ch.Offset + pos.OffsetInChapter
}

// absoluteOffsetToPosition converts a rune offset into a logical
// Position by finding the containing chapter and offset within it.
func (m Model) absoluteOffsetToPosition(offset int) reader.Position {
	if m.currentBook == nil || offset <= 0 {
		return reader.Position{}
	}
	chapters := m.currentBook.Book.Chapters
	if len(chapters) == 0 {
		return reader.Position{}
	}
	chapterIndex := 0
	for i, ch := range chapters {
		if offset < ch.Offset+ch.Length {
			chapterIndex = i
			break
		}
	}
	ch := chapters[chapterIndex]
	if offset < ch.Offset {
		offset = ch.Offset
	}
	offsetInChapter := offset - ch.Offset
	if offsetInChapter < 0 {
		offsetInChapter = 0
	}
	return reader.Position{ChapterIndex: chapterIndex, OffsetInChapter: offsetInChapter}
}

// View renders the full-screen layout with a menu bar, main area with
// pseudo-graphics borders, and a status bar at the bottom.
func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return "thujareader – initializing..."
	}

	if m.height < 3 {
		// Not enough space to render full layout; show a compact message.
		return "Terminal too small for thujareader UI. Resize the window."
	}

	var b strings.Builder

	// Top menu bar.
	b.WriteString(m.theme.applyMenuBar(m.renderMenuBar()))
	b.WriteRune('\n')

	// Main area bordered with pseudo-graphics.
	top := string(m.theme.borderTopLeft) + strings.Repeat(string(m.theme.borderHorizontal), max(0, m.width-2)) + string(m.theme.borderTopRight)
	bottom := string(m.theme.borderBottomLeft) + strings.Repeat(string(m.theme.borderHorizontal), max(0, m.width-2)) + string(m.theme.borderBottomRight)
	b.WriteString(top)
	b.WriteRune('\n')

	innerHeight := m.height - 3 // minus top menu and bottom border + status bar
	if innerHeight < 1 {
		innerHeight = 1
	}

	for i := 0; i < innerHeight-1; i++ {
		b.WriteRune(m.theme.borderVertical)

		innerWidth := max(0, m.width-2)
		// When a menu is open, render its items in the top lines of the
		// main area so that selecting a menu visibly opens a dropdown.
		if m.menuOpen && m.activeMenu >= 0 && m.activeMenu < len(m.menus) {
			items := m.menus[m.activeMenu].items
			if i < len(items) {
				label := items[i].label

				// Indicate the active item with a leading marker; other
				// items are padded with a space for alignment.
				var line string
				if i == m.activeItem {
					line = ">" + label
				} else {
					line = " " + label
				}

				b.WriteString(padOrTrim(line, innerWidth))
			} else {
				b.WriteString(strings.Repeat(" ", innerWidth))
			}
		} else if m.inputMode && i == 0 {
			// Show a simple one-line input prompt at the top of the main
			// area when collecting a file path.
			line := m.inputPrompt + string(m.inputBuffer)
			b.WriteString(padOrTrim(line, innerWidth))
		} else if m.tocOpen && m.currentBook != nil {
			// Render a simple TOC dialog: list of entries with the
			// currently selected one highlighted.
			idx := i
			if idx >= 0 && idx < len(m.currentBook.TOC) {
				entry := m.currentBook.TOC[idx]
				label := entry.Label
				if idx == m.tocIndex {
					label = "> " + label
				} else {
					label = "  " + label
				}
				b.WriteString(padOrTrim(label, innerWidth))
			} else {
				b.WriteString(strings.Repeat(" ", innerWidth))
			}
		} else if m.bookmarksOpen && m.currentBook != nil {
			// Render a simple bookmarks dialog: list of bookmark names with
			// the currently selected one highlighted.
			list := m.currentBookmarks()
			idx := i
			if idx >= 0 && idx < len(list) {
				entry := list[idx]
				label := entry.Name
				if idx == m.bookmarkIndex {
					label = "> " + label
				} else {
					label = "  " + label
				}
				b.WriteString(padOrTrim(label, innerWidth))
			} else {
				b.WriteString(strings.Repeat(" ", innerWidth))
			}
		} else if m.currentBook != nil {
			// Render wrapped book text starting from topLine.
			idx := m.topLine + i
			if idx >= 0 && idx < len(m.lines) {
				line := m.lines[idx]
				b.WriteString(padOrTrim(line, innerWidth))
			} else {
				b.WriteString(strings.Repeat(" ", innerWidth))
			}
		} else {
			b.WriteString(strings.Repeat(" ", innerWidth))
		}

		b.WriteRune(m.theme.borderVertical)
		b.WriteRune('\n')
	}

	b.WriteString(bottom)

	// Status bar on the last line.
	b.WriteRune('\n')
	b.WriteString(m.theme.applyStatusBar(m.renderStatusBar()))

	return b.String()
}

func (m Model) renderMenuBar() string {
	var segments []string
	for i, menu := range m.menus {
		label := " " + menu.label + " "
		if m.menuOpen && i == m.activeMenu {
			segments = append(segments, "["+label+"]")
		} else {
			segments = append(segments, " "+label+" ")
		}
	}
	line := strings.Join(segments, "")
	return padOrTrim(line, m.width)
}

func (m Model) renderStatusBar() string {
	text := m.statusLine
	location := ""
	if m.currentBook != nil && len(m.currentBook.Book.Chapters) > 0 {
		// Compute approximate progress percentage based on
		// TotalCharacters and current position.
		book := m.currentBook.Book
		if book.TotalCharacters > 0 {
			abs := m.positionToAbsoluteOffset(m.currentPos)
			if abs < 0 {
				abs = 0
			}
			if abs > book.TotalCharacters {
				abs = book.TotalCharacters
			}
			percent := (abs * 100) / book.TotalCharacters
			chapterIndex := m.currentPos.ChapterIndex
			chapterLabel := ""
			if chapterIndex >= 0 && chapterIndex < len(book.Chapters) {
				ch := book.Chapters[chapterIndex]
				if strings.TrimSpace(ch.Title) != "" {
					chapterLabel = ch.Title
				} else {
					chapterLabel = "Chapter " + itoa(chapterIndex+1)
				}
			}
			if chapterLabel != "" {
				location = chapterLabel + " "
			}
			location += itoa(percent) + "%"
		}
	}

	if location != "" {
		// Place location info at the right edge, trimming or padding the
		// main status text as needed. All widths are rune/column-aware.
		locWidth := runewidth.StringWidth(location)
		if locWidth > m.width {
			location = runewidth.Truncate(location, m.width, "")
			locWidth = runewidth.StringWidth(location)
		}
		available := m.width - locWidth - 1
		if available < 0 {
			available = 0
		}
		text = padOrTrim(text, available)
		text += " " + location
	} else {
		text = padOrTrim(text, m.width)
	}

	return text
}

// itoa is a small helper for integer-to-string conversion.
func itoa(i int) string {
	return strconv.Itoa(i)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// padOrTrim returns a version of s whose printable width (in terminal
// cells) is exactly width, using runewidth to account for multi-byte
// and wide runes. If s is wider, it is truncated; if it is narrower,
// it is padded with spaces on the right.
func padOrTrim(s string, width int) string {
	if width <= 0 {
		return ""
	}
	w := runewidth.StringWidth(s)
	if w > width {
		return runewidth.Truncate(s, width, "")
	}
	if w < width {
		return s + strings.Repeat(" ", width-w)
	}
	return s
}
