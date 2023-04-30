package ui

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"planner/task"
	"planner/ui/add_task"
	"strconv"
)

type Model struct {
	mode Mode

	taskList list.Model
	addTask  add_task.Model
	helpText string

	windowWidth  int
	windowHeight int
}

type Mode struct {
	id    int
	label string
}

type taskItem task.Task

func (i taskItem) Title() string       { return i.Name }
func (i taskItem) Description() string { return i.Desc }
func (i taskItem) FilterValue() string { return strconv.FormatInt(i.ID, 20) }

var (
	database *sql.DB

	normalModeStyle = lipgloss.NewStyle().Background(lipgloss.Color("#74c7ec")).Foreground(lipgloss.Color("#000")).MarginTop(1)
	addModeStyle    = lipgloss.NewStyle().Background(lipgloss.Color("#eba0ac")).Foreground(lipgloss.Color("#000")).MarginTop(1)
)

func (m *Model) setMode(modeId int) {
	var modeLabels = map[int]string{
		0: "Normal",
		1: "Add Task",
	}

	m.mode = Mode{
		id:    modeId,
		label: modeLabels[modeId],
	}
}

func createHelpText(m *Model, modeId int) string {
	var helpText string
	prefix := "[" + m.mode.label + "]"

	switch modeId {
	case 0:
		helpText = normalModeStyle.Render(prefix) + " " + "(n)ew task / ctrl+(q)uit"
	case 1:
		helpText = addModeStyle.Render(prefix) + " " + "(tab) next input / (shift+tab) previous input / ctrl+(q)uit"
	}

	return helpText
}

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

func reloadTasks(m *Model, database *sql.DB) {
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
		addTask:      add_task.InitialModel(db),
		windowWidth:  120,
		windowHeight: 20,
	}

	m.helpText = createHelpText(&m, 0)

	m.setMode(0)

	m.taskList.SetShowTitle(false)
	m.taskList.SetItems(taskItems)
	m.taskList.SetShowHelp(false)
	m.taskList.SetShowStatusBar(false)

	m.taskList.DisableQuitKeybindings()

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case add_task.AddTaskMsg:
		_, err := task.AddTask(database, msg.Task)
		if err != nil {
			fmt.Sprintf("Error adding task: %v", err)
		} else {
			reloadTasks(&m, database)
		}

	case tea.KeyMsg:
		if m.mode.id == 0 {
			switch msg.String() {
			case "down":
				if m.taskList.Cursor() == len(m.taskList.Items())-1 {
					add_task.SetFocused(&m.addTask, true)
				}
			case "d":
				currentItem := m.taskList.Items()[m.taskList.Cursor()].(taskItem)

				_, err := task.DeleteTask(database, currentItem.ID)
				if err != nil {
					fmt.Sprintf("Error deleting task: %v", err)
				}

				reloadTasks(&m, database)
			}
		}

		if m.mode.id != 1 {
			switch msg.String() {
			case "n":
				m.setMode(1)
				add_task.SetFocused(&m.addTask, true)
				m.helpText = createHelpText(&m, m.mode.id)

				return m, cmd
			}
		}

		// global keybindings
		switch msg.String() {
		case "esc":
			m.setMode(0)
			add_task.SetFocused(&m.addTask, false)
		case "ctrl+q":
			return m, tea.Quit
		}
	}

	m.helpText = createHelpText(&m, m.mode.id)

	if m.mode.id == 0 {
		m.taskList, cmd = m.taskList.Update(msg)
	}

	if m.mode.id == 1 {
		var addTaskModel tea.Model

		addTaskModel, cmd = m.addTask.Update(msg)
		m.addTask = addTaskModel.(add_task.Model)
	}

	return m, cmd
}

func (m Model) View() string {
	return m.taskList.View() +
		"\n" +
		m.addTask.View() +
		"\n" +
		m.helpText
}
