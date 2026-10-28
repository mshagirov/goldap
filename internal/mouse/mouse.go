// Package mouse hit-tests mouse positions against rendered views.
package mouse

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// Region is a span of cells on a single screen line.
type Region struct {
	X, Y, Width int
}

// Contains reports whether the cell at x, y is inside the region.
func (r Region) Contains(x, y int) bool {
	return y == r.Y && x >= r.X && x < r.X+r.Width
}

// Find returns the region of the first occurrence of text in a rendered view.
func Find(view, text string) (Region, bool) {
	if text == "" {
		return Region{}, false
	}
	for y, line := range strings.Split(view, "\n") {
		line = ansi.Strip(line)
		if i := strings.Index(line, text); i > -1 {
			return Region{X: ansi.StringWidth(line[:i]), Y: y, Width: ansi.StringWidth(text)}, true
		}
	}
	return Region{}, false
}

// Hit reports whether the cell at x, y is on the first occurrence of text in view.
func Hit(view, text string, x, y int) bool {
	r, ok := Find(view, text)
	return ok && r.Contains(x, y)
}
