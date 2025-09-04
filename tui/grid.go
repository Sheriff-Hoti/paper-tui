package tui

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	COLS         = 3
	ROWS         = 3
	PAGE_SIZE    = COLS * ROWS
	ROWS_SPACING = 1
	COLS_SPACING = 1
	TOP_SPACING  = 1
	LEFT_SPACING = 1
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

			row_idx := uint32(idx / COLS)
			col_idx := uint32(idx % COLS)
			img_width := uint32((init_term_width / COLS) - 2)
			img_height := uint32((init_term_height / ROWS) - 2)

			cell_page = append(cell_page, cell{
				filename:   file,
				id:         uint32(idx),
				row_idx:    row_idx,
				col_idx:    col_idx,
				row_cell:   (row_idx * img_height) + TOP_SPACING + (ROWS_SPACING * row_idx),
				col_cell:   (col_idx * img_width) + LEFT_SPACING + (COLS_SPACING * col_idx),
				img_width:  img_width,
				img_height: img_height,
			})
		}

		cells = append(cells, cell_page)
	}

	//debbuginn
	s := ""
	for idx, cell_page := range cells {
		s += fmt.Sprint("page:", idx, "\n")
		for _, cell := range cell_page {
			s += fmt.Sprint(cell.id, ":", cell.filename, ":  :col:", cell.col_idx, " ", "row:", cell.row_idx, " widht:", cell.img_width, " height:", cell.img_height, "\n")
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
			g.go_left()
			return g, nil
		case key.Matches(msg, g.keys.right):
			g.go_right()
			return g, nil
		case key.Matches(msg, g.keys.up):
			g.go_up()
			return g, nil
		case key.Matches(msg, g.keys.down):
			g.go_down()
			return g, nil
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return g, nil
}

func (g *grid) View() string {
	// The header
	page := g.cells[g.page_index]
	cell := page[g.cell_index]

	file := filepath.Base(cell.filename)
	square := lipgloss.NewStyle().
		Width(int(cell.img_width)).
		Height(int(cell.img_height)).
		// Background(lipgloss.Color("12")).        // blue square
		MarginTop(int(cell.row_cell)). // y position
		MarginLeft(int(cell.col_cell)).
		Border(lipgloss.RoundedBorder(), true).
		Render(file)

	background := lipgloss.NewStyle().
		Width(int(g.window_width)).
		Height(int(g.window_height)).
		Background(lipgloss.Color("235"))

	// view := lipgloss.JoinVertical(lipgloss.Center, background, g.paginator.View())
	return background.Render(square)

}

func (g *grid) go_up() {
	//guardrails
	if len(g.cells) == 0 {
		return
	}
	//guardrails
	if g.page_index < 0 && g.page_index >= len(g.cells) {
		return
	}

	page := g.cells[g.page_index]
	//guardrails
	if len(page) == 0 {
		return
	}
	//guardrails
	if g.cell_index < 0 && g.cell_index >= len(page) {
		return
	}

	if (g.cell_index - COLS) < 0 {
		return
	}

	g.cell_index -= COLS

}

func (g *grid) go_down() {
	//guardrails
	if len(g.cells) == 0 {
		return
	}
	//guardrails
	if g.page_index < 0 && g.page_index >= len(g.cells) {
		return
	}

	page := g.cells[g.page_index]
	//guardrails
	if len(page) == 0 {
		return
	}
	//guardrails
	if g.cell_index < 0 && g.cell_index >= len(page) {
		return
	}

	if (g.cell_index + COLS) >= len(page) {
		return
	}

	g.cell_index += COLS

}

func (g *grid) go_left() {
	// guardrails
	if len(g.cells) == 0 {
		return
	}
	if g.page_index < 0 || g.page_index >= len(g.cells) {
		return
	}

	page := g.cells[g.page_index]
	if len(page) == 0 {
		return
	}
	if g.cell_index < 0 || g.cell_index >= len(page) {
		return
	}

	// get current cell's column/row
	cur := page[g.cell_index]
	col := int(cur.col_idx)
	row := int(cur.row_idx)

	// if we're in the first column, move to the previous page keeping the same row
	if col == 0 {
		prev := g.page_index - 1
		if prev < 0 {
			prev = len(g.cells) - 1
		}
		prevPage := g.cells[prev]
		if len(prevPage) == 0 {
			return
		}
		// target = same row, last column
		target := row*COLS + (COLS - 1)
		if target >= len(prevPage) {
			// clamp to last cell on the page if the target index doesn't exist
			target = len(prevPage) - 1
		}
		g.page_index = prev
		g.cell_index = target
		return
	}

	// normal left move
	g.cell_index--
}

func (g *grid) go_right() {
	// guardrails
	if len(g.cells) == 0 {
		return
	}
	if g.page_index < 0 || g.page_index >= len(g.cells) {
		return
	}

	page := g.cells[g.page_index]
	if len(page) == 0 {
		return
	}
	if g.cell_index < 0 || g.cell_index >= len(page) {
		return
	}

	// get current cell's column/row
	cur := page[g.cell_index]
	col := int(cur.col_idx)
	row := int(cur.row_idx)

	// if we're in the last column, move to the next page keeping the same row (first column)
	if col >= COLS-1 {
		next := g.page_index + 1
		if next >= len(g.cells) {
			next = 0
		}
		nextPage := g.cells[next]
		if len(nextPage) == 0 {
			return
		}
		// target = same row, first column
		target := row * COLS
		if target >= len(nextPage) {
			// clamp to last cell on the page if the target index doesn't exist
			target = len(nextPage) - 1
		}
		g.page_index = next
		g.cell_index = target
		return
	}

	// normal right move; if this would overflow the page (partial last page), wrap to next page
	if g.cell_index >= len(page)-1 {
		next := g.page_index + 1
		if next >= len(g.cells) {
			next = 0
		}
		nextPage := g.cells[next]
		if len(nextPage) == 0 {
			return
		}
		g.page_index = next
		g.cell_index = 0
		return
	}

	g.cell_index++
}
