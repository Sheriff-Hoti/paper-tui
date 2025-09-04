package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type page_change_msg struct {
	cells     []cell
	old_cells []cell
}

func change_page_cmd(old_cells []cell, cells []cell) tea.Cmd {
	return func() tea.Msg {
		return page_change_msg{cells: cells, old_cells: old_cells}
	}
}
