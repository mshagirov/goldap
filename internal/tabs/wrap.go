package tabs

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

const (
	boxMargin       = 2  // free columns kept on each side of a box
	minBoxTextWidth = 20 // boxes never wrap text narrower than this
)

// wrapText wraps s to width cells, breaking at spaces or, for long words
// such as DNs, at commas; words that still don't fit are broken anywhere.
func wrapText(s string, width int) string {
	if width <= 0 {
		return s
	}
	return ansi.Wrap(s, width, ",")
}

// boxTextWidth returns the widest text line that fits in a box drawn with
// style on a screen screenWidth cells wide.
func boxTextWidth(style lipgloss.Style, screenWidth int) int {
	return max(screenWidth-style.GetHorizontalFrameSize()-2*boxMargin, minBoxTextWidth)
}
