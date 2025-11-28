package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	tabsModel "github.com/mshagirov/goldap/internal/tabs"
)

func runTabs(tabs []string, tabContent []string) {

	m := tabsModel.Model{Tabs: tabs, TabContent: tabContent}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
