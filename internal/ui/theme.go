package ui

import "os"

// Theme describes colors and basic pseudo-graphics characters used by
// the TUI. It intentionally stays very small so it can be wired to a
// configuration system in later phases.
type Theme struct {
	// ANSI escape sequences (without reset) for the major regions.
	menuBarPrefix   string
	statusBarPrefix string
	reset           string

	// Box-drawing characters. For very limited terminals these can fall
	// back to ASCII characters.
	borderTopLeft     rune
	borderTopRight    rune
	borderBottomLeft  rune
	borderBottomRight rune
	borderHorizontal  rune
	borderVertical    rune
}

// DefaultTheme returns a theme approximating the classic DOS edit.exe
// look: blue background with a cyan menu bar and highlighted status
// bar. It uses ANSI colors and assumes a terminal with basic color
// support.
func DefaultTheme() Theme {
	return Theme{
		// Cyan menu bar on blue background with bright white text.
		menuBarPrefix:   "\x1b[1;37;46m",
		statusBarPrefix: "\x1b[1;37;44m",
		reset:           "\x1b[0m",

		borderTopLeft:     '┌',
		borderTopRight:    '┐',
		borderBottomLeft:  '└',
		borderBottomRight: '┘',
		borderHorizontal:  '─',
		borderVertical:    '│',
	}
}

// NoColorTheme provides a safe fallback for terminals without color
// support. It keeps the same layout but omits ANSI sequences and
// replaces box-drawing characters with ASCII where possible.
func NoColorTheme() Theme {
	return Theme{
		menuBarPrefix:   "",
		statusBarPrefix: "",
		reset:           "",

		borderTopLeft:     '+',
		borderTopRight:    '+',
		borderBottomLeft:  '+',
		borderBottomRight: '+',
		borderHorizontal:  '-',
		borderVertical:    '|',
	}
}

// ThemeFromEnv chooses a theme based on environment hints. This forms
// a minimal configuration hook that can later be replaced or
// augmented by a full configuration system.
//
// If THUJAREADER_NO_COLOR is set to any non-empty value, a no-color
// theme is returned; otherwise the default ANSI-based theme is used.
func ThemeFromEnv() Theme {
	if v := os.Getenv("THUJAREADER_NO_COLOR"); v != "" {
		return NoColorTheme()
	}
	return DefaultTheme()
}

// applyMenuBar colors a menu bar line according to the theme.
func (t Theme) applyMenuBar(line string) string {
	if t.menuBarPrefix == "" {
		return line
	}
	return t.menuBarPrefix + line + t.reset
}

// applyStatusBar colors a status bar line according to the theme.
func (t Theme) applyStatusBar(line string) string {
	if t.statusBarPrefix == "" {
		return line
	}
	return t.statusBarPrefix + line + t.reset
}
