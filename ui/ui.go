package ui

import (
	"database/sql"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"planner/task"
	"planner/ui/task_entry"
	"planner/ui/task_list"
)

type (
	Model struct {
		mode Mode

		taskList  task_list.Model
		taskEntry task_entry.Model
		helpText  string

		windowWidth  int
		windowHeight int
	}
	Mode struct {
		id    int
		label string
	}
)

var database *sql.DB

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

func InitialModel(db *sql.DB) Model {
	database = db

	m := Model{
		taskList:     task_list.InitialModel(db),
		taskEntry:    task_entry.InitialModel(),
		windowWidth:  120,
		windowHeight: 20,
	}

	m.helpText = createHelpText(&m, 0)

	setMode(&m, 0)

	storedTasks, err := task.GetTasks(db)
	if err != nil {
		fmt.Printf("Error getting tasks: %v\n", err)
		return m
	}

	task_list.SetItems(&m.taskList, storedTasks)
	m.taskList.List.SetShowTitle(false)
	m.taskList.List.SetShowHelp(false)
	m.taskList.List.SetShowStatusBar(false)
	m.taskList.List.DisableQuitKeybindings()

	return m
}

func updateTaskEntryModel(m *Model, msg tea.Msg) (task_entry.Model, tea.Cmd) {
	addTaskModel, cmd := m.taskEntry.Update(msg)

	return addTaskModel.(task_entry.Model), cmd
}

func updateTaskListModel(m *Model, msg tea.Msg) (task_list.Model, tea.Cmd) {
	taskListModel, cmd := m.taskList.Update(msg)

	return taskListModel.(task_list.Model), cmd
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
			task_list.ReloadTasks(&m.taskList, database)
			setMode(&m, 0)
		}
	case task_entry.EditTaskMsg:
		_, err := task.UpdateTask(database, msg.Task)
		if err != nil {
			fmt.Printf("Error adding task: %v", err)
		} else {
			task_list.ReloadTasks(&m.taskList, database)
			setMode(&m, 0)
		}

	case tea.KeyMsg:
		if m.mode.id == 0 {
			switch msg.String() {
			case "d":
				_, err := task.DeleteTask(database, task_list.GetCurrentTask(&m.taskList).ID)
				if err != nil {
					fmt.Printf("Error deleting task: %v", err)
				}

				task_list.ReloadTasks(&m.taskList, database)
			case "n":
				setMode(&m, 1)
				task_entry.SetFocused(&m.taskEntry, true)
				m.helpText = createHelpText(&m, m.mode.id)

				return m, cmd
			case "e":
				setMode(&m, 2)
				task_entry.SetFocused(&m.taskEntry, true)
				task_entry.LoadTask(&m.taskEntry, task_list.GetCurrentTask(&m.taskList))
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
		m.taskList, cmd = updateTaskListModel(&m, msg)
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
