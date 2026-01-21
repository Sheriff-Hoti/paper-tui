package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/Sheriff-Hoti/paper-tui/data"
)

type Config struct {
	Init_hook             string `json:"init_hook"`
	Wallpaper_select_hook string `json:"wallpaper_select_hook"`
	Wallpaper_dir         string `json:"wallpaper_dir"`
	Data_dir              string `json:"data_dir"`
}

// const allowedImageExtensions = ".jpg,.jpeg,.png,.gif"

func ReadConfigFile(config_path string) (*Config, error) {
	//TODO check in main.go if config file is passed via flag if yes and the file does not exist return error, if its not passed via flag return default config
	config := GetDefaultConfigVals()

	if _, err := os.Stat(config_path); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist and if it does not exists just return the defaults

		return config, nil
	}

	file, err := os.Open(config_path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func GetDefaultConfigPath() string {
	const (
		xdgConfigHome = "XDG_CONFIG_HOME"
	)

	if val, ok := os.LookupEnv(xdgConfigHome); ok {
		return filepath.Join(val, "paper-tui", "config.json")
	}

	// fallback to $HOME/.config/paper-tui/config.json
	home, err := os.UserHomeDir()
	if err != nil {
		// if home can't be resolved, fallback to current working directory
		return filepath.Join(".", "config.json")
	}

	return filepath.Join(home, ".config", "paper-tui", "config.json")
}

func GetDefaultConfigVals() *Config {

	return &Config{
		Init_hook:             "",
		Wallpaper_select_hook: "",
		Wallpaper_dir:         "",
		Data_dir:              data.GetDefaultDataPath(),
	}
}
