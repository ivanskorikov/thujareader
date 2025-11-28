package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

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
	lines       []string
	topLine     int

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
		activeMenu: -1,
		activeItem: 0,
		statusLine: "Press F10 or Alt key combinations to open menus. F1 for Help.",
	}

	if book != nil {
		m.setBook(*book)
	}

	return m
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
	case tea.KeyF3:
		m.executeCommand(cmdOpen)
		return true
	}

	// Alt+<letter> opens corresponding menu (e.g., Alt+F for File).
	if msg.Alt && len(msg.Runes) == 1 {
		m.openMenuByAltKey(msg.Runes[0])
		return true
	}

	if !m.menuOpen {
		// When the menu is not open, allow basic scrolling of the
		// currently loaded book using arrow keys. More advanced
		// navigation is handled in later phases.
		switch msg.Type {
		case tea.KeyUp:
			if m.topLine > 0 {
				m.topLine--
			}
			return true
		case tea.KeyDown:
			if m.currentBook != nil && m.topLine < len(m.lines)-1 {
				m.topLine++
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
		m.setStatus("Find: not yet implemented (search dialog will appear in later phase).")
	case cmdToc:
		m.setStatus("TOC: not yet implemented (table of contents view will appear in later phase).")
	case cmdBookmarks:
		m.setStatus("Bookmarks: not yet implemented (bookmark dialog will appear in later phase).")
	case cmdRecentFiles:
		m.setStatus("Recent files: not yet implemented (recent list will appear in later phase).")
	case cmdHelp:
		m.setStatus("Help: not yet implemented (help screen will appear in later phase).")
	default:
		return
	}
}

func (m *Model) setStatus(text string) {
	m.statusLine = text
	m.statusDirty = true
}

// setBook installs a newly loaded book into the model and prepares a
// simple line-based view over its text.
func (m *Model) setBook(book reader.LoadedBook) {
	m.currentBook = &book
	if book.Text != "" {
		m.lines = strings.Split(book.Text, "\n")
	} else {
		m.lines = nil
	}
	m.topLine = 0
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

				if len(line) > innerWidth {
					line = line[:innerWidth]
				}
				if len(line) < innerWidth {
					line += strings.Repeat(" ", innerWidth-len(line))
				}
				b.WriteString(line)
			} else {
				b.WriteString(strings.Repeat(" ", innerWidth))
			}
		} else if m.inputMode && i == 0 {
			// Show a simple one-line input prompt at the top of the main
			// area when collecting a file path.
			line := m.inputPrompt + string(m.inputBuffer)
			if len(line) > innerWidth {
				line = line[:innerWidth]
			}
			if len(line) < innerWidth {
				line += strings.Repeat(" ", innerWidth-len(line))
			}
			b.WriteString(line)
		} else if m.currentBook != nil {
			// Render book text starting from topLine.
			idx := m.topLine + i
			if idx >= 0 && idx < len(m.lines) {
				line := m.lines[idx]
				if len(line) > innerWidth {
					line = line[:innerWidth]
				}
				if len(line) < innerWidth {
					line += strings.Repeat(" ", innerWidth-len(line))
				}
				b.WriteString(line)
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
	if len(line) > m.width {
		return line[:m.width]
	}
	if len(line) < m.width {
		line += strings.Repeat(" ", m.width-len(line))
	}
	return line
}

func (m Model) renderStatusBar() string {
	text := m.statusLine
	if len(text) > m.width {
		text = text[:m.width]
	}
	if len(text) < m.width {
		text += strings.Repeat(" ", m.width-len(text))
	}
	return text
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
