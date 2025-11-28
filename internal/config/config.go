package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// Config holds user-editable settings loaded from a JSON file. The
// structure is intentionally minimal for Phase 5 and can be extended in
// later phases without breaking existing configs (unknown fields are
// ignored on load).
type Config struct {
	// ThemeOverride allows selecting an alternate theme if supported by
	// the UI. For now this is a free-form string.
	ThemeOverride string `json:"theme_override,omitempty"`

	// RecentListSize limits the number of recent files remembered. If
	// zero or negative, a sensible default is used.
	RecentListSize int `json:"recent_list_size,omitempty"`

	// DefaultLibraryPath, when set, can be used as a starting directory
	// for file-open dialogs or path prompts.
	DefaultLibraryPath string `json:"default_library_path,omitempty"`
}

// DefaultConfig returns a Config populated with built-in defaults.
func DefaultConfig() Config {
	return Config{
		ThemeOverride:      "",
		RecentListSize:     10,
		DefaultLibraryPath: "",
	}
}

// Paths groups the resolved locations of the configuration and state
// files on disk so callers do not need to repeat this logic.
type Paths struct {
	ConfigFile string
	StateFile  string
}

// DefaultPaths computes per-user paths for the config and state JSON
// files. On Windows it uses %APPDATA%\thujareader; on Unix-like systems
// it uses $XDG_CONFIG_HOME/thujareader or ~/.config/thujareader.
func DefaultPaths() (Paths, error) {
	var base string
	if runtime.GOOS == "windows" {
		base = os.Getenv("APPDATA")
		if base == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return Paths{}, err
			}
			base = filepath.Join(home, "AppData", "Roaming")
		}
		base = filepath.Join(base, "thujareader")
	} else {
		base = os.Getenv("XDG_CONFIG_HOME")
		if base == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return Paths{}, err
			}
			base = filepath.Join(home, ".config")
		}
		base = filepath.Join(base, "thujareader")
	}

	return Paths{
		ConfigFile: filepath.Join(base, "config.json"),
		StateFile:  filepath.Join(base, "state.json"),
	}, nil
}

// Load reads configuration from the given path. If the file does not
// exist, DefaultConfig is returned with a nil error. If the file is
// present but invalid, a non-nil error is returned so callers can
// decide how to proceed.
func Load(path string) (Config, error) {
	if path == "" {
		return DefaultConfig(), errors.New("config path is empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return DefaultConfig(), err
	}
	if len(data) == 0 {
		return DefaultConfig(), nil
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}
	return cfg, nil
}

// Save writes the provided configuration to disk as JSON, creating the
// parent directory if needed.
func Save(path string, cfg Config) error {
	if path == "" {
		return errors.New("config path is empty")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
