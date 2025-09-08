package tui

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

const (
	KITTY_IMG_HDR = "\x1b_G"
	KITTY_IMG_FTR = "\x1b\\"
)

type KittyImgOpts struct {
	SrcX        uint32 // x=
	SrcY        uint32 // y=
	SrcWidth    uint32 // w=
	SrcHeight   uint32 // h=
	CellOffsetX uint32 // X= (pixel x-offset inside terminal cell)
	CellOffsetY uint32 // Y= (pixel y-offset inside terminal cell)
	DstCols     uint32 // c= (display width in terminal columns)
	DstRows     uint32 // r= (display height in terminal rows)
	ZIndex      int32  // z=
	ImageId     uint32 // i=
	ImageNo     uint32 // I=
	PlacementId uint32 // p=
	Cursor      uint32
}

func (o KittyImgOpts) ToHeader(opts ...string) string {

	type fldmap struct {
		pv   *uint32
		code rune
	}
	sFld := []fldmap{
		{&o.SrcX, 'x'},
		{&o.SrcY, 'y'},
		{&o.SrcWidth, 'w'},
		{&o.SrcHeight, 'h'},
		{&o.CellOffsetX, 'X'},
		{&o.CellOffsetY, 'Y'},
		{&o.DstCols, 'c'},
		{&o.Cursor, 'C'},
		{&o.DstRows, 'r'},
		{&o.ImageId, 'i'},
		{&o.ImageNo, 'I'},
		{&o.PlacementId, 'p'},
	}

	for _, f := range sFld {
		if *f.pv != 0 {
			opts = append(opts, fmt.Sprintf("%c=%d", f.code, *f.pv))
		}
	}

	if o.ZIndex != 0 {
		opts = append(opts, fmt.Sprintf("z=%d", o.ZIndex))
	}

	return KITTY_IMG_HDR + strings.Join(opts, ",") + ";"
}

type cell struct {
	filename    string
	img_width   uint32
	img_height  uint32
	row_idx     uint32
	col_idx     uint32
	row_cell    uint32
	col_cell    uint32
	id          uint32
	initialized bool
	last_dims   int
}

func (c *cell) RenderImage(out io.Writer, opts KittyImgOpts) error {
	// Build the Kitty header
	header := opts.ToHeader("a=T", "f=100", "t=f")

	// Write header
	if _, err := fmt.Fprint(out, header); err != nil {
		return err
	}

	// Encode the absolute path in base64 (required by Kitty)
	enc64 := base64.NewEncoder(base64.StdEncoding, out)
	if _, err := fmt.Fprint(enc64, c.filename); err != nil {
		return err
	}
	if err := enc64.Close(); err != nil {
		return err
	}

	// Write the terminal escape sequence to finish
	if _, err := fmt.Fprint(out, KITTY_IMG_FTR); err != nil {
		return err
	}
	return nil
}

func (c *cell) Hide(out io.Writer, opts KittyImgOpts) error {
	header := opts.ToHeader("a=d", "d=i", fmt.Sprintf("i=%d", opts.ImageId), fmt.Sprintf("p=%d", opts.PlacementId))

	if _, err := fmt.Fprint(out, header); err != nil {
		return err
	}
	if _, err := fmt.Fprint(out, KITTY_IMG_FTR); err != nil {
		return err
	}

	return nil
}

func (c *cell) Show(out io.Writer, opts KittyImgOpts) error {

	header := opts.ToHeader("a=p", fmt.Sprintf("i=%d", opts.ImageId), fmt.Sprintf("p=%d", opts.PlacementId))

	if _, err := fmt.Fprint(out, header); err != nil {
		return err
	}
	if _, err := fmt.Fprint(out, KITTY_IMG_FTR); err != nil {
		return err
	}

	return nil
}

func (c *cell) Update(term_width int, term_height int) bool {

	dims := (term_width << 16) | (term_height & 0xFFFF)

	changed := dims != c.last_dims

	if !changed {
		return changed
	}

	c.img_width = uint32((term_width / COLS) - 2)
	c.img_height = uint32((term_height / ROWS) - 2)
	c.row_cell = (c.row_idx * c.img_height) + TOP_SPACING + (ROWS_SPACING * c.row_idx)
	c.col_cell = (c.col_idx * c.img_width) + LEFT_SPACING + (COLS_SPACING * c.col_idx)

	c.last_dims = dims

	return true
}
