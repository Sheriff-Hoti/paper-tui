package tui

import "github.com/charmbracelet/bubbles/key"

type gridKeyMap struct {
	quit key.Binding
	hide key.Binding
}

func gridKeyMaps() *gridKeyMap {
	return &gridKeyMap{
		quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("ctrl+c/q/esc", "quit"),
		),
		hide: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "hide"),
		),
	}
}
