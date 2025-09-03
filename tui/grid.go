package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type grid struct {
	keys *gridKeyMap
}

func NewGrid(absfiles []string, selected_file string, init_term_width int, init_term_height int) *grid {
	return &grid{
		keys: gridKeyMaps(),
		// Our to-do list is a grocery list

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
	}
}

func (g *grid) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (g *grid) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch {

		// These keys should exit the program.
		case key.Matches(msg, g.keys.quit):
			return g, tea.Quit

		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return g, nil
}

func (g *grid) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
