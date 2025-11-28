package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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

// Model holds UI state for the Phase 2 TUI shell emulating DOS edit.exe.
type Model struct {
	width  int
	height int

	theme Theme

	menus       []menu
	activeMenu  int  // index into menus, -1 when no menu is active
	activeItem  int  // index into items of the active menu
	menuOpen    bool // whether menu bar interaction is active
	statusLine  string
	statusDirty bool
}

// NewModel constructs the initial UI model.
func NewModel() Model {
	m := Model{
		// Start with a reasonable default size so that the UI can render
		// even if no WindowSizeMsg is delivered (which can happen on some
		// terminals, especially on Windows). Resize events will override
		// these values when they arrive.
		width:  80,
		height: 25,
		theme:  ThemeFromEnv(),
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
		m.setStatus("Open: not yet implemented (will open file dialog in later phase).")
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

				innerWidth := max(0, m.width-2)
				if len(line) > innerWidth {
					line = line[:innerWidth]
				}
				if len(line) < innerWidth {
					line += strings.Repeat(" ", innerWidth-len(line))
				}
				b.WriteString(line)
			} else {
				b.WriteString(strings.Repeat(" ", max(0, m.width-2)))
			}
		} else {
			b.WriteString(strings.Repeat(" ", max(0, m.width-2)))
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
