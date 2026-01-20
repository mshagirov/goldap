package tabs

import "github.com/charmbracelet/lipgloss"

var (
	// tabs
	highlightColor    = lipgloss.AdaptiveColor{Light: "#DAA520", Dark: "#FFD700"}
	goldbar           = lipgloss.Color("#DAA520")
	blurredColor      = lipgloss.Color("241")
	highContrastColor = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}
	grays             = lipgloss.AdaptiveColor{Light: "#3B3B3B", Dark: "#ADADAD"}

	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(0, 0, 0, 0)
	inactiveTabStyle  = lipgloss.NewStyle().Foreground(blurredColor).Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Foreground(highlightColor).BorderForeground(highlightColor).Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().
				BorderForeground(highlightColor).
				Align(lipgloss.Left).
				BorderStyle(lipgloss.NormalBorder()).
				Border(lipgloss.NormalBorder()).
				UnsetBorderTop()
	fillerBorderStyle = lipgloss.NewStyle().Border(
		lipgloss.Border{Bottom: "─", BottomRight: "┐"}, false, true, true, false).
		BorderForeground(highlightColor)
	infoBarStyle   = lipgloss.NewStyle().Foreground(grays).Align(lipgloss.Right)
	searchBarStyle = lipgloss.NewStyle().Foreground(highContrastColor).Align(lipgloss.Left)

	// tables
	tableHighlightColor  = lipgloss.AdaptiveColor{Light: "#0014a8", Dark: "#265ef7"}
	tableForegroundColor = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}

	// forms
	formTitleStyle     = lipgloss.NewStyle().Foreground(highlightColor).Padding(0, 1).Border(lipgloss.RoundedBorder()).BorderForeground(blurredColor)
	formFieldNameStyle = lipgloss.NewStyle().Foreground(goldbar).Padding(0, 1)
	formInputPadding   = lipgloss.NewStyle().Padding(0, 2)
	formBlurredStyle   = lipgloss.NewStyle().Foreground(grays)
	formActiveStyle    = lipgloss.NewStyle().Foreground(tableForegroundColor).Background(blurredColor)
	formModifiedStyle  = lipgloss.NewStyle().Foreground(tableHighlightColor)
	formHelpStyle      = lipgloss.NewStyle().Foreground(grays)

	msgBoxStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(blurredColor)

	msgStyle           = lipgloss.NewStyle()
	msgBlurredBtnStyle = msgStyle.Foreground(blurredColor)
	msgFocusedBtnStyle = msgStyle.Foreground(highlightColor)

	msgTextStyle  = msgStyle.Foreground(blurredColor).Padding(0, 2)
	msgTitleStyle = msgStyle.Foreground(blurredColor).Bold(true).PaddingLeft(2)
)
