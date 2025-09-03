package tui

import "github.com/charmbracelet/bubbles/key"

type gridKeyMap struct {
	quit  key.Binding
	hide  key.Binding
	left  key.Binding
	right key.Binding
	up    key.Binding
	down  key.Binding
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
		left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("<-", "left"),
		),
		right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("->", "right"),
		),
		up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("^", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("d", "down"),
		),
	}
}
