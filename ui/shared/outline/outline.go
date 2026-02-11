package outline

import (
	"strings"
)

const (
	Cols = 8
	Rows = 8
)

const (
	FirstCol = 0
	FirstRow = 0
	LastCol  = Cols - 1
	LastRow  = Rows - 1
)
const (
	// cellHeight represents how many rows are in a cell
	cellHeight = 2
	// cellWidth represents how many columns are in a cell
	cellWidth = 4

	// marginLeft and marginTop represent the offset of the chess
	// board from the top left of the terminal window. This is to
	// account for padding and rank labels
	marginLeft = 3
	marginTop  = 1

	Vertical   = "│"
	Horizontal = "─"
)

// Build returns a string with a border for a given row (top, middle, bottom)
func Build(left, middle, right string) string {
	border := left + Horizontal + strings.Repeat(Horizontal+Horizontal+middle+Horizontal, LastRow)
	border += Horizontal + Horizontal + right + "\n"
	return withMarginLeft(border)
}

// withMarginLeft returns a string with a prepended left margin
func withMarginLeft(s string) string {
	return strings.Repeat(" ", marginLeft) + s
}

// Top returns a built border with the top row
func Top() string {
	return Build("┌", "┬", "┐")
}

// Middle returns a built border with the middle row
func Middle() string {
	return Build("├", "┼", "┤")
}

// Bottom returns a built border with the bottom row
func Bottom() string {
	return Build("└", "┴", "┘")
}
