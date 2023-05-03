package ui

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"planner/task"
	"planner/ui/task_entry"
	"strconv"
)

type (
	Model struct {
		mode Mode

		taskList  list.Model
		taskEntry task_entry.Model
		helpText  string

		windowWidth  int
		windowHeight int
	}
	Mode struct {
		id    int
		label string
	}
	taskItem task.Task
)

func (i taskItem) Title() string       { return i.Name }
func (i taskItem) Description() string { return i.Desc }
func (i taskItem) FilterValue() string { return strconv.FormatInt(i.ID, 20) }
func (i taskItem) getTask() task.Task  { return task.Task(i) }

var (
	database *sql.DB

	normalModeStyle = lipgloss.NewStyle().Background(lipgloss.Color("#74c7ec")).Foreground(lipgloss.Color("#000")).MarginTop(1)
	addModeStyle    = lipgloss.NewStyle().Background(lipgloss.Color("#eba0ac")).Foreground(lipgloss.Color("#000")).MarginTop(1)
	editModeStyle   = lipgloss.NewStyle().Background(lipgloss.Color("#89b4fa")).Foreground(lipgloss.Color("#000")).MarginTop(1)
)

func setMode(m *Model, modeId int) {
	var modeLabels = map[int]string{
		0: "View",
		1: "Add",
		2: "Edit",
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
		helpText = normalModeStyle.Render(prefix) + " " + "(n)ew task / (e)dit task / ctrl+(q)uit"
	case 1:
		helpText = addModeStyle.Render(prefix) + " " + "(esc) cancel / (enter) submit / (tab) next input / (shift+tab) previous input / ctrl+(q)uit"
	case 2:
		helpText = editModeStyle.Render(prefix) + " " + "(esc) cancel / (enter) submit / (tab) next input / (shift+tab) previous input / ctrl+(q)uit"
	}

	return helpText
}

func getCurrentItem(m *Model) taskItem {
	return m.taskList.Items()[m.taskList.Cursor()].(taskItem)
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
		taskEntry:    task_entry.InitialModel(),
		windowWidth:  120,
		windowHeight: 20,
	}

	m.helpText = createHelpText(&m, 0)

	setMode(&m, 0)

	m.taskList.SetShowTitle(false)
	m.taskList.SetItems(taskItems)
	m.taskList.SetShowHelp(false)
	m.taskList.SetShowStatusBar(false)

	m.taskList.DisableQuitKeybindings()

	return m
}

func updateTaskEntryModel(m *Model, msg tea.Msg) (task_entry.Model, tea.Cmd) {
	addTaskModel, cmd := m.taskEntry.Update(msg)

	return addTaskModel.(task_entry.Model), cmd
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case task_entry.AddTaskMsg:
		_, err := task.AddTask(database, msg.Task)
		if err != nil {
			fmt.Printf("Error adding task: %v", err)
		} else {
			reloadTasks(&m, database)
			setMode(&m, 0)
		}
	case task_entry.EditTaskMsg:
		_, err := task.UpdateTask(database, msg.Task)
		if err != nil {
			fmt.Printf("Error adding task: %v", err)
		} else {
			reloadTasks(&m, database)
			setMode(&m, 0)
		}

	case tea.KeyMsg:
		if m.mode.id == 0 {
			switch msg.String() {
			case "d":
				_, err := task.DeleteTask(database, getCurrentItem(&m).ID)
				if err != nil {
					fmt.Printf("Error deleting task: %v", err)
				}

				reloadTasks(&m, database)
			case "n":
				setMode(&m, 1)
				task_entry.SetFocused(&m.taskEntry, true)
				m.helpText = createHelpText(&m, m.mode.id)

				return m, cmd
			case "e":
				setMode(&m, 2)
				task_entry.SetFocused(&m.taskEntry, true)
				task_entry.LoadTask(&m.taskEntry, getCurrentItem(&m).getTask())
				m.helpText = createHelpText(&m, m.mode.id)

				return m, cmd
			}
		}

		// global keybindings
		switch msg.String() {
		case "esc":
			setMode(&m, 0)
			task_entry.SetFocused(&m.taskEntry, false)

			m.taskEntry, _ = updateTaskEntryModel(&m, msg)
		case "ctrl+q":
			return m, tea.Quit
		}
	}

	m.helpText = createHelpText(&m, m.mode.id)

	if m.mode.id == 0 {
		m.taskList, cmd = m.taskList.Update(msg)
	}

	if m.mode.id == 1 || m.mode.id == 2 {
		m.taskEntry, cmd = updateTaskEntryModel(&m, msg)
	}

	return m, cmd
}

func (m Model) View() string {
	return m.taskList.View() +
		"\n" +
		m.taskEntry.View() +
		"\n" +
		m.helpText
}
