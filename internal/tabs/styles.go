package tabs

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
)

var (
	// tabs
	highlightColor    = compat.AdaptiveColor{Light: lipgloss.Color("#DAA520"), Dark: lipgloss.Color("#FFD700")}
	goldbar           = lipgloss.Color("#DAA520")
	blurredColor      = lipgloss.Color("241")
	highContrastColor = compat.AdaptiveColor{Light: lipgloss.Color("#000000"), Dark: lipgloss.Color("#ffffff")}
	grays             = compat.AdaptiveColor{Light: lipgloss.Color("#3B3B3B"), Dark: lipgloss.Color("#ADADAD")}

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
	tableHighlightColor  = highlightColor // compat.AdaptiveColor{Light: lipgloss.Color("#0014a8"), Dark: lipgloss.Color("#265ef7")}
	tableForegroundColor = compat.AdaptiveColor{Light: highContrastColor.Dark, Dark: highContrastColor.Light}

	// forms
	formHeaderStyle    = lipgloss.NewStyle().Foreground(highlightColor).Padding(0, 1).Border(lipgloss.RoundedBorder()).BorderForeground(blurredColor)
	formFieldNameStyle = lipgloss.NewStyle().Foreground(goldbar).Padding(0, 1)
	formInputPadding   = lipgloss.NewStyle().Padding(0, 2)
	formBlurredStyle   = lipgloss.NewStyle().Foreground(grays)
	formActiveStyle    = lipgloss.NewStyle().Foreground(tableForegroundColor).Background(blurredColor)
	formModifiedStyle  = lipgloss.NewStyle().Foreground(tableHighlightColor)
	formFooterStyle    = lipgloss.NewStyle().Foreground(grays)

	msgBoxStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Padding(2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(blurredColor)

	msgStyle           = lipgloss.NewStyle()
	msgBlurredBtnStyle = msgStyle.Foreground(blurredColor).Padding(0, 1)
	msgFocusedBtnStyle = msgStyle.Foreground(highlightColor).Padding(0, 1)

	msgTextStyle  = msgStyle.Foreground(blurredColor).Align(lipgloss.Left)
	msgTitleStyle = msgStyle.Foreground(highlightColor).Bold(true) //.PaddingLeft(2).PaddingRight(2)
)
