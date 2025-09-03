package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	COLS      = 3
	ROWS      = 3
	PAGE_SIZE = COLS * ROWS
)

type grid struct {
	keys  *gridKeyMap
	cells [][]cell
}

func NewGrid(abs_files []string, selected_file string, init_term_width int, init_term_height int) *grid {

	chunked := make([][]string, 0, (len(abs_files)+PAGE_SIZE-1)/PAGE_SIZE)

	cells := make([][]cell, 0, (len(abs_files)+PAGE_SIZE-1)/PAGE_SIZE)

	for i := 0; i < len(abs_files); i += PAGE_SIZE {
		end := min(i+PAGE_SIZE, len(abs_files))
		chunked = append(chunked, abs_files[i:end])
		cell_page := make([]cell, 0, PAGE_SIZE)

		for idx, file := range abs_files[i:end] {
			cell_page = append(cell_page, cell{
				filename: file,
				id:       uint32(idx),
				RowCell:  uint32(idx / COLS),
				ColCell:  uint32(idx % COLS),
			})
		}

		cells = append(cells, cell_page)
	}

	//debbuginn
	s := ""
	for idx, cell_page := range cells {
		s += fmt.Sprint("page:", idx, "\n")
		for _, cell := range cell_page {
			s += fmt.Sprint(cell.id, ":", cell.filename, ":  :col:", cell.ColCell, " ", "row:", cell.RowCell, "\n")
		}
		s += "\n"
	}

	fmt.Print(s)

	return &grid{
		keys: gridKeyMaps(),
		// Our to-do list is a grocery list
		cells: cells,
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
	for idx, cell_page := range g.cells {
		for _, cell := range cell_page {
			s += fmt.Sprint(cell.id, ":", cell.filename, "\n")
		}
		s += "\n"
		s += fmt.Sprint("page:", idx, "\n")
	}
	// Send the UI for rendering
	return s
}
