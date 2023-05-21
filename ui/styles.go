package ui

import lg "github.com/charmbracelet/lipgloss"

var (
	normalModeStyle = lg.NewStyle().
			Background(lg.Color("#a6e3a1")).
			Foreground(lg.Color("#000")).
			MarginTop(1)
	addModeStyle = lg.NewStyle().
			Background(lg.Color("#74c7ec")).
			Foreground(lg.Color("#000")).
			MarginTop(1)
	editModeStyle = lg.NewStyle().
			Background(lg.Color("#fab387")).
			Foreground(lg.Color("#000")).
			MarginTop(1)
	statusBarStyle = lg.NewStyle().Foreground(lg.Color("#585b70"))

	windowStyle = lg.NewStyle().Padding(1)
)
