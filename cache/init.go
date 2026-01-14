package cache

import (
	"os"
	"path/filepath"
)

func GetCacheDir() string {
	const (
		xdgCacheHome = "XDG_CACHE_HOME"
	)

	if val, ok := os.LookupEnv(xdgCacheHome); ok {
		return filepath.Join(val, "paper-tui")
	}

	// fallback to $HOME/.cache/paper-tui
	home, err := os.UserHomeDir()
	if err != nil {
		// if home can't be resolved, fallback to current working directory
		return filepath.Join(".", "cache")
	}

	return filepath.Join(home, ".cache", "paper-tui")
}
