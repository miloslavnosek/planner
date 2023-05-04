package ui

import "github.com/charmbracelet/lipgloss"

var (
	normalModeStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#74c7ec")).
			Foreground(lipgloss.Color("#000")).
			MarginTop(1)
	addModeStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#eba0ac")).
			Foreground(lipgloss.Color("#000")).
			MarginTop(1)
	editModeStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#89b4fa")).
			Foreground(lipgloss.Color("#000")).
			MarginTop(1)
)
