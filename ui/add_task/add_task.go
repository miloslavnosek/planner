package add_task

import (
	"database/sql"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"planner/task"
)

type Model struct {
	isFocused bool

	nameInput textinput.Model
	descInput textinput.Model
}

type AddTaskMsg struct {
	Task task.Task
}

var (
	database *sql.DB

	containerStyle = lipgloss.NewStyle()
)

func sendAddTaskMsg(msg AddTaskMsg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func clearForm(m *Model) {
	m.nameInput.Reset()
	m.descInput.Reset()
}

func SetFocused(m *Model, focus bool) {
	if focus && !m.isFocused {
		m.isFocused = true

		m.nameInput.Focus()
	} else if !focus && m.isFocused {
		m.isFocused = false

		m.nameInput.Blur()
		m.descInput.Blur()
	}
}

func InitialModel(db *sql.DB) Model {
	m := Model{
		nameInput: textinput.New(),
		descInput: textinput.New(),
	}

	m.nameInput.Placeholder = "name of the task"
	m.descInput.Placeholder = "description"

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd, inputCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			cmd = sendAddTaskMsg(AddTaskMsg{task.Task{Name: m.nameInput.Value(), Desc: m.descInput.Value()}})
			clearForm(&m)
			SetFocused(&m, false)
		case "tab":
			m.nameInput.Blur()
			m.descInput.Focus()
		case "shift+tab":
			m.nameInput.Focus()
			m.descInput.Blur()
		case "esc":
			m.descInput.Blur()
			m.nameInput.Blur()
		}
	}

	if m.nameInput.Focused() {
		m.nameInput, inputCmd = m.nameInput.Update(msg)
	} else if m.descInput.Focused() {
		m.descInput, inputCmd = m.descInput.Update(msg)
	}

	return m, tea.Batch(cmd, inputCmd)
}

func (m Model) View() string {
	return containerStyle.Render(
		m.nameInput.View() +
			" " +
			m.descInput.View(),
	)
}
