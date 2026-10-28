package mouse

import "testing"

func TestHit(t *testing.T) {
	view := "╭──────────╮\n│ \x1b[1m[C]ancel\x1b[0m  [S]ave │\n╰──────────╯"
	cases := []struct {
		name string
		text string
		x, y int
		want bool
	}{
		{name: "first cell", text: "[C]ancel", x: 2, y: 1, want: true},
		{name: "last cell", text: "[C]ancel", x: 9, y: 1, want: true},
		{name: "before text", text: "[C]ancel", x: 1, y: 1, want: false},
		{name: "after text", text: "[C]ancel", x: 10, y: 1, want: false},
		{name: "wrong line", text: "[C]ancel", x: 2, y: 0, want: false},
		{name: "second button", text: "[S]ave", x: 12, y: 1, want: true},
		{name: "missing text", text: "[ OK ]", x: 2, y: 1, want: false},
		{name: "empty text", text: "", x: 0, y: 0, want: false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Hit(view, c.text, c.x, c.y); got != c.want {
				t.Errorf("Hit(%q, %d, %d) = %v, want %v", c.text, c.x, c.y, got, c.want)
			}
		})
	}
}
