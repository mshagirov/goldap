package tabs

import "github.com/charmbracelet/lipgloss"

var (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")

	// tabs
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(0, 0, 0, 0)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#DAA520", Dark: "#FFD700"}
	blurredColor      = lipgloss.Color("241")
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
	infoBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#3B3B3B", Dark: "#ADADAD"}).
			Align(lipgloss.Right)
	searchBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}).
			Align(lipgloss.Left)

	// tables
	tableHighlightColor  = lipgloss.AdaptiveColor{Light: "#0014a8", Dark: "#265ef7"}
	tableForegroundColor = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}

	// forms
	formTitleStyle = lipgloss.NewStyle().Foreground(highlightColor)
	inputStyle     = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle  = lipgloss.NewStyle().Foreground(darkGray)
)
