package task_list

import lg "github.com/charmbracelet/lipgloss"

var (
	selectedTitleStyle = lg.NewStyle().Foreground(lg.Color("#cdd6f4")).Bold(true)
	selectedDescStyle  = lg.NewStyle().Foreground(lg.Color("#bac2de"))
	filterMatchStyle   = lg.NewStyle().Background(lg.Color("#cdd6f4")).Foreground(lg.Color("#000"))
	normalTitleStyle   = lg.NewStyle().Foreground(lg.Color("#585b70")).Bold(true)
	normalDescStyle    = lg.NewStyle().Foreground(lg.Color("#45475a"))
)
