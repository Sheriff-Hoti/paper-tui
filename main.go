package main

import (
	"flag"
	// "fmt"
	// "log"
	// "os"

	// "github.com/Sheriff-Hoti/paper-tui/backend"
	"github.com/Sheriff-Hoti/paper-tui/config"
	// "github.com/Sheriff-Hoti/paper-tui/data"
	// "github.com/Sheriff-Hoti/paper-tui/tui"
	// tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	// "golang.org/x/term"
)

func main() {
	config_file := flag.String("config", config.GetDefaultConfigPath(), "config file")
	// init := flag.Bool("init", false, "Init flag")

	flag.Parse()

	configS := viper.New()
	configS.SetConfigName("config")
	configS.SetConfigType("yaml")
	//TODO check if flag is not null and if it is not null then use setConfigFile
	configS.AddConfigPath(*config_file)
	configS.AddConfigPath("$XDG_CONFIG_HOME/paper-tui")
	configS.AddConfigPath("$HOME/paper-tui")

	// err := configS.ReadInConfig()

	// config_struct, err := config.ReadConfigFile(*config_file)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// data, err := data.ReadDataFile(config_struct.Data_dir)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// back := backend.InitBackend()

	// //check init here
	// if *init {
	// 	log.Print(data.Current_wallpaper)
	// 	//and check if data is already initialized
	// 	back.SetImage(data.Current_wallpaper)
	// 	return
	// }

	// files, err := config.GetWallpapers(config_struct.Wallpaper_dir)
	// if err != nil {
	// 	log.Fatal("Error trying to get wallpapers:", err)
	// }

	// width, height, err := term.GetSize(int(os.Stdout.Fd()))

	// if err != nil {
	// 	log.Fatal("Error trying to get terminal size:", err)
	// }

	// p := tea.NewProgram(
	// 	tui.NewGrid(files, config_struct, data, width, height, back), tea.WithAltScreen())
	// if _, err := p.Run(); err != nil {
	// 	fmt.Printf("Alas, there's been an error: %v", err)
	// 	os.Exit(1)
	// }
}
