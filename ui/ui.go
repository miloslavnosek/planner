package ui

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"planner/task"
	"strconv"
)

type Model struct {
	taskList  list.Model
	taskInput textinput.Model

	windowWidth  int
	windowHeight int
}

type taskItem task.Task

func (i taskItem) Title() string       { return i.Name }
func (i taskItem) Description() string { return i.Desc }
func (i taskItem) FilterValue() string { return strconv.FormatInt(i.ID, 20) }

var (
	database *sql.DB

	inputContainerStyle lipgloss.Style
	listContainerStyle  lipgloss.Style
)

func createContainerStyle(focused bool, width, height int) lipgloss.Style {
	borderColor := lipgloss.Color("#CCC")
	if focused {
		borderColor = lipgloss.Color("#FF5733")
	}

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Width(width).
		Height(height)
}

func ConvertTasksToItems(tasks []task.Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, t := range tasks {
		items[i] = taskItem(t)
	}
	return items
}

func updateContainerStyle(m *Model) {
	inputContainerStyle = createContainerStyle(m.taskInput.Focused(), 80, 0)
	listContainerStyle = createContainerStyle(!m.taskInput.Focused(), 80, 3)
}

func reloadTasks(m *Model, database *sql.DB) {
	m.taskInput.Reset()
	storedTasks, _ := task.GetTasks(database)
	taskItems := ConvertTasksToItems(storedTasks)
	m.taskList.SetItems(taskItems)
}

func InitialModel(db *sql.DB) Model {
	database = db

	storedTasks, _ := task.GetTasks(db)
	taskItems := ConvertTasksToItems(storedTasks)

	m := Model{
		taskList:     list.New([]list.Item{}, list.NewDefaultDelegate(), 80, 20),
		taskInput:    textinput.New(),
		windowWidth:  120,
		windowHeight: 20,
	}

	m.taskList.SetShowTitle(false)
	m.taskList.SetItems(taskItems)
	m.taskList.SetShowHelp(false)
	m.taskList.SetShowStatusBar(false)
	m.taskList.KeyMap.CursorUp.SetKeys("up")
	m.taskList.KeyMap.CursorDown.SetKeys("down")
	m.taskList.KeyMap.PrevPage.SetKeys()
	m.taskList.KeyMap.NextPage.SetKeys()

	m.taskInput.Focus()
	m.taskList.DisableQuitKeybindings()

	updateContainerStyle(&m)

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.taskInput.Focused() {
			switch msg.String() {
			case "ctrl+enter":
				_, err := task.AddTask(database, task.Task{Name: m.taskInput.Value()})
				if err != nil {
					fmt.Sprintf("Error adding task: %v", err)
				} else {
					reloadTasks(&m, database)
				}
			case "up", "esc":
				m.taskInput.Blur()
			}
		} else if !m.taskInput.Focused() {
			switch msg.String() {
			case "down":
				if m.taskList.Cursor() == len(m.taskList.Items())-1 {
					m.taskInput.Focus()
				}
			}
		}

		switch msg.String() {
		case "ctrl+q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		return m, nil
	}

	if m.taskInput.Focused() {
		m.taskInput, cmd = m.taskInput.Update(msg)
	} else {
		m.taskList, cmd = m.taskList.Update(msg)
	}

	updateContainerStyle(&m)

	return m, cmd
}

func (m Model) View() string {
	return listContainerStyle.Render(m.taskList.View()) + "\n" + inputContainerStyle.Render(m.taskInput.View())
}
