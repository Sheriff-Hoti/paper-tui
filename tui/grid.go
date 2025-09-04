package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sheriff-Hoti/paper-tui/backend"
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
	COLS_SPACING = 2
	TOP_SPACING  = 1
	LEFT_SPACING = 1
)

type grid struct {
	keys          *gridKeyMap
	cells         [][]*cell
	page_index    int
	cell_index    int
	window_width  uint32
	window_height uint32
	paginator     paginator.Model
	backend       backend.WallpaperBackend
}

func NewGrid(abs_files []string, selected_file string, init_term_width int, init_term_height int, back backend.WallpaperBackend) *grid {

	backen := backend.InitBackend()
	cells := make([][]*cell, 0, (len(abs_files)+PAGE_SIZE-1)/PAGE_SIZE)

	p := paginator.New()
	p.Type = paginator.Dots
	// p.PerPage = 10
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")

	idCounter := 1

	for i := 0; i < len(abs_files); i += PAGE_SIZE {
		end := min(i+PAGE_SIZE, len(abs_files))
		// chunked = append(chunked, abs_files[i:end])
		cell_page := make([]*cell, 0, PAGE_SIZE)

		for idx, file := range abs_files[i:end] {

			row_idx := uint32(idx / COLS)
			col_idx := uint32(idx % COLS)
			img_width := uint32((init_term_width / COLS) - 2)
			img_height := uint32((init_term_height / ROWS) - 2)

			cell_page = append(cell_page, &cell{
				filename:   file,
				id:         uint32(idCounter),
				row_idx:    row_idx,
				col_idx:    col_idx,
				row_cell:   (row_idx * img_height) + TOP_SPACING + (ROWS_SPACING * row_idx),
				col_cell:   (col_idx * img_width) + LEFT_SPACING + (COLS_SPACING * col_idx),
				img_width:  img_width,
				img_height: img_height,
			})

			idCounter++

		}

		cells = append(cells, cell_page)
	}

	p.SetTotalPages(len(cells))

	//debbuginn
	// s := ""
	// for idx, cell_page := range cells {
	// 	s += fmt.Sprint("page:", idx, "\n")
	// 	for _, cell := range cell_page {
	// 		s += fmt.Sprint(cell.id, ":", cell.filename, ":  :col:", cell.col_idx, " ", "row:", cell.row_idx, " widht:", cell.img_width, " height:", cell.img_height, "\n")
	// 	}
	// 	s += "\n"
	// }

	// fmt.Print(s)

	return &grid{
		keys:          gridKeyMaps(),
		cells:         cells,
		window_width:  uint32(init_term_width),
		window_height: uint32(init_term_height),
		paginator:     p,
		backend:       backen,
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
	}
}

func (g *grid) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	first_page := g.cells[0]
	for _, cell := range first_page {
		fmt.Fprintf(os.Stdout, "\x1b[%d;%dH", cell.row_cell+1, cell.col_cell+1)
		cell.RenderImage(os.Stdout, KittyImgOpts{
			DstCols: uint32(cell.img_width),
			DstRows: uint32(cell.img_height),
			// 100 pixels =1 row
			// SrcX:      109,
			// SrcY:      109,
			// Cursor: 1,
			// SrcWidth:  uint32(800),
			// SrcHeight: uint32(500),
			//TODO: fix it
			ImageId:     cell.id,
			PlacementId: cell.id,
		})
		cell.initialized = true
	}
	return nil
}

func (g *grid) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// var cmd []tea.Cmd
	switch msg := msg.(type) {
	case page_change_msg:
		for _, cell := range msg.old_cells {
			cell.Hide(os.Stdout, KittyImgOpts{
				ImageId:     cell.id,
				PlacementId: cell.id,
			})
		}
		for _, cell := range msg.cells {

			img_opts := KittyImgOpts{
				DstCols: uint32(cell.img_width),
				DstRows: uint32(cell.img_height),
				// 100 pixels =1 row
				// SrcX:      109,
				// SrcY:      109,
				// Cursor: 1,
				// SrcWidth:  uint32(800),
				// SrcHeight: uint32(500),
				//TODO: fix it
				ImageId:     cell.id,
				PlacementId: cell.id,
			}

			fmt.Fprintf(os.Stdout, "\x1b[%d;%dH", cell.row_cell+1, cell.col_cell+1)
			cell.RenderImage(os.Stdout, img_opts)

			//show-render strategy doesnt work very well because kittys storage quotas,
			//  if the quota is reached then the show func wont show images because they hae been cleared
			// if cell.initialized {
			// 	cell.Show(os.Stdout, img_opts)
			// } else {

			cell.initialized = true
			// }

		}

	// 	return g, nil
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch {

		// These keys should exit the program.
		case key.Matches(msg, g.keys.quit):
			return g, tea.Quit

		case key.Matches(msg, g.keys.left):
			cmd := g.go_left()
			return g, cmd
		case key.Matches(msg, g.keys.right):
			cmd := g.go_right()
			return g, cmd
		case key.Matches(msg, g.keys.up):
			g.go_up()
			return g, nil
		case key.Matches(msg, g.keys.down):
			g.go_down()
			return g, nil

		case key.Matches(msg, g.keys.select_cell):
			g.backend.SetImage(g.cells[g.page_index][g.cell_index].filename)
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

	s := fmt.Sprintf("%s %d %t", file, cell.id, cell.initialized)
	square := lipgloss.NewStyle().
		Width(int(cell.img_width)).
		Height(int(cell.img_height)).
		// Background(lipgloss.Color("12")).        // blue square
		MarginTop(int(cell.row_cell)). // y position
		MarginLeft(int(cell.col_cell-1)).
		Border(lipgloss.RoundedBorder(), true).
		Render(s)

	background := lipgloss.NewStyle().
		Width(int(g.window_width)).
		Height(int(g.window_height)).
		Background(lipgloss.Color("235")).Render(square)

	g.paginator.Page = g.page_index
	return lipgloss.JoinVertical(lipgloss.Center, background, g.paginator.View())
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

func (g *grid) go_left() tea.Cmd {
	// guardrails
	if len(g.cells) == 0 {
		return nil
	}
	if g.page_index < 0 || g.page_index >= len(g.cells) {
		return nil
	}

	page := g.cells[g.page_index]
	if len(page) == 0 {
		return nil
	}
	if g.cell_index < 0 || g.cell_index >= len(page) {
		return nil
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
			return nil
		}
		// target = same row, last column
		target := row*COLS + (COLS - 1)
		if target >= len(prevPage) {
			// clamp to last cell on the page if the target index doesn't exist
			target = len(prevPage) - 1
		}
		g.page_index = prev
		g.cell_index = target
		return change_page_cmd(page, prevPage)

	}

	// normal left move
	g.cell_index--
	return nil
}

func (g *grid) go_right() tea.Cmd {
	// guardrails
	if len(g.cells) == 0 {
		return nil
	}
	if g.page_index < 0 || g.page_index >= len(g.cells) {
		return nil
	}

	page := g.cells[g.page_index]
	if len(page) == 0 {
		return nil
	}
	if g.cell_index < 0 || g.cell_index >= len(page) {
		return nil
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
			return nil
		}
		// target = same row, first column
		target := row * COLS
		if target >= len(nextPage) {
			// clamp to last cell on the page if the target index doesn't exist
			target = len(nextPage) - 1
		}
		g.page_index = next
		g.cell_index = target
		return change_page_cmd(page, nextPage)

	}

	// normal right move; if this would overflow the page (partial last page), wrap to next page
	if g.cell_index >= len(page)-1 {
		next := g.page_index + 1
		if next >= len(g.cells) {
			next = 0
		}
		nextPage := g.cells[next]
		if len(nextPage) == 0 {
			return nil
		}
		g.page_index = next
		g.cell_index = 0
		return change_page_cmd(page, nextPage)
	}

	g.cell_index++
	return nil
}
