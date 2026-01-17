package cache

import (
	"os"
	"path/filepath"
)

func GetCacheDir() (string, error) {
	const (
		xdgCacheHome = "XDG_CACHE_HOME"
	)

	if val, ok := os.LookupEnv(xdgCacheHome); ok {
		dir := filepath.Join(val, "paper-tui")
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
		return dir, nil
	}

	// fallback to $HOME/.cache/paper-tui
	home, err := os.UserHomeDir()
	if err != nil {
		// if home can't be resolved, fallback to current working directory
		dir := filepath.Join(".", "cache")
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
		return dir, nil
	}

	dir := filepath.Join(home, ".cache", "paper-tui")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}
