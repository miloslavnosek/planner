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
	mode string

	taskList list.Model
	addTask  add_task.Model

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
		mode:         "normal", // "normal" or "add-task
		taskList:     list.New([]list.Item{}, list.NewDefaultDelegate(), 80, 20),
		addTask:      add_task.InitialModel(db),
		windowWidth:  120,
		windowHeight: 20,
	}

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
		if m.mode == "normal" {
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

		// global keybindings
		switch msg.String() {
		case "n":
			m.mode = "add-task"
			add_task.SetFocused(&m.addTask, true)
		case "esc":
			m.mode = "normal"
			add_task.SetFocused(&m.addTask, false)
		case "ctrl+q":
			return m, tea.Quit
		}
	}

	m.taskList, cmd = m.taskList.Update(msg)

	if m.mode == "add-task" {
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
		"[" + m.mode + "]" + " " + "(d)elete selected task / ctrl+(q)uit"
}
