package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Data struct {
	Current_wallpaper string `json:"current_wallpaper"`
	Init              bool   `json:"init"`
}

func GetDefaultDataPath() string {
	const (
		xdgConfigHome = "XDG_DATA_HOME"
	)

	if val, ok := os.LookupEnv(xdgConfigHome); ok {
		return filepath.Join(val, "paper-tui", "data.json")
	}

	// fallback to $HOME/.config/paper-tui/config.json
	home, err := os.UserHomeDir()
	if err != nil {
		// if home can't be resolved, fallback to current working directory
		return filepath.Join(".", "data.json")
	}

	return filepath.Join(home, ".local", "share", "paper-tui", "data.json")
}

func ReadDataFile(dataPath string) (*Data, error) {
	dataPath = os.ExpandEnv(dataPath)

	// If file does not exist → return empty Data
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		return &Data{}, nil
	}

	// Open file (read only)
	file, err := os.Open(dataPath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	// If empty → return empty Data
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat file: %w", err)
	}
	if stat.Size() == 0 {
		return &Data{}, nil
	}

	// Decode JSON
	var data Data
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("decoding json: %w", err)
	}

	return &data, nil
}

func WriteToDataFile(data *Data, dataPath string) error {
	dataPath = os.ExpandEnv(dataPath)

	// Ensure parent dirs exist
	if err := os.MkdirAll(filepath.Dir(dataPath), 0755); err != nil {
		return fmt.Errorf("creating parent dirs: %w", err)
	}

	// Open file with truncation
	file, err := os.OpenFile(dataPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	// Encode as JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // pretty-print for readability
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("encoding json: %w", err)
	}

	return nil
}
