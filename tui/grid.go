package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	COLS      = 3
	ROWS      = 3
	PAGE_SIZE = COLS * ROWS
)

type grid struct {
	keys          *gridKeyMap
	cells         [][]cell
	page_index    int
	cell_index    int
	window_width  uint32
	window_height uint32
	paginator     paginator.Model
}

func NewGrid(abs_files []string, selected_file string, init_term_width int, init_term_height int) *grid {

	// chunked := make([][]string, 0, (len(abs_files)+PAGE_SIZE-1)/PAGE_SIZE)

	cells := make([][]cell, 0, (len(abs_files)+PAGE_SIZE-1)/PAGE_SIZE)

	p := paginator.New()
	p.Type = paginator.Dots
	// p.PerPage = 10
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(10)
	for i := 0; i < len(abs_files); i += PAGE_SIZE {
		end := min(i+PAGE_SIZE, len(abs_files))
		// chunked = append(chunked, abs_files[i:end])
		cell_page := make([]cell, 0, PAGE_SIZE)

		for idx, file := range abs_files[i:end] {
			cell_page = append(cell_page, cell{
				filename: file,
				id:       uint32(idx),
				row_idx:  uint32(idx / COLS),
				col_idx:  uint32(idx % COLS),
				Width:    uint32((init_term_width / COLS) - 2),
				Height:   uint32((init_term_height / ROWS) - 2),
			})
		}

		cells = append(cells, cell_page)
	}

	//debbuginn
	s := ""
	for idx, cell_page := range cells {
		s += fmt.Sprint("page:", idx, "\n")
		for _, cell := range cell_page {
			s += fmt.Sprint(cell.id, ":", cell.filename, ":  :col:", cell.col_idx, " ", "row:", cell.row_idx, " widht:", cell.Width, " height:", cell.Height, "\n")
		}
		s += "\n"
	}

	fmt.Print(s)

	return &grid{
		keys:          gridKeyMaps(),
		cells:         cells,
		window_width:  uint32(init_term_width),
		window_height: uint32(init_term_height),
		paginator:     p,
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

		case key.Matches(msg, g.keys.left):
			g.cell_index--
			return g, nil
		case key.Matches(msg, g.keys.right):
			g.cell_index++
			return g, nil
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return g, nil
}

func (g *grid) View() string {
	// The header
	page := g.cells[0]
	cell := page[g.cell_index]
	square := lipgloss.NewStyle().
		Width(int(cell.Width)).
		Height(int(cell.Height)).
		// Background(lipgloss.Color("12")).        // blue square
		MarginTop(int(cell.row_idx)). // y position
		MarginLeft(int(cell.col_idx)).
		Border(lipgloss.RoundedBorder(), true).
		Render("")

	background := lipgloss.NewStyle().
		Width(int(g.window_width)).
		Height(int(g.window_height)).
		Background(lipgloss.Color("235"))

	// view := lipgloss.JoinVertical(lipgloss.Center, background, g.paginator.View())
	return background.Render(square)

}
