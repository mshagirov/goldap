package tabs

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"golang.org/x/term"

	"github.com/mshagirov/goldap/ldapapi"
)

func NewTable(ti ldapapi.TableInfo) table.Model {
	if len(ti.Rows) == 0 {
		ti.Rows = append(ti.Rows, make([]string, len(ti.Cols)))
	}

	w, h := GetTableDimensions()

	t := table.New(
		table.WithColumns(ti.Cols),
		table.WithRows(ti.Rows),
		table.WithFocused(true),
		table.WithHeight(h),
		table.WithWidth(w),
		table.WithStyles(GetTableStyle()),
	)
	return t
}

func newTableWithFilter(ti ldapapi.TableInfo, filter string) table.Model {
	if len(filter) == 0 {
		return NewTable(ti)
	}
	filter = strings.ToLower(filter)

	newRows := []table.Row{}

	var contains bool

	for _, row := range ti.Rows {
		contains = false
		for _, col := range row {
			if strings.Contains(strings.ToLower(col), filter) {
				contains = true
				break
			}
		}
		if contains {
			newRows = append(newRows, row)
		}
	}

	if len(newRows) > 0 {
		ti.Rows = newRows
	}
	return NewTable(ti)
}

func GetTableStyle() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.Foreground(tableHighlightColor)
	s.Selected = s.Selected.Foreground(tableForegroundColor).Background(tableHighlightColor)
	return s
}

func GetTableDimensions() (int, int) {
	termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth, termHeight = 20, 20
	}
	w, h := windowStyle.GetHorizontalFrameSize(), windowStyle.GetVerticalFrameSize()
	return (termWidth - w), (termHeight - 6*h)
}
