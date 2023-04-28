package ui

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"planner/task"
	"strconv"
)

type Model struct {
	taskList  list.Model
	taskInput textinput.Model
}

type taskItem task.Task

func (i taskItem) Title() string       { return i.Name }
func (i taskItem) Description() string { return i.Desc }
func (i taskItem) FilterValue() string { return strconv.FormatInt(i.ID, 10) }

var database *sql.DB

func ConvertTasksToItems(tasks []task.Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, t := range tasks {
		items[i] = taskItem(t)
	}
	return items
}

func reloadTasks(m *Model, database *sql.DB) {
	m.taskInput.Reset()
	storedTasks, _ := task.GetTasks(database)
	taskItems := ConvertTasksToItems(storedTasks)
	m.taskList.SetItems(taskItems)
}

func InitialModel(db *sql.DB) Model {
	storedTasks, _ := task.GetTasks(db)
	taskItems := ConvertTasksToItems(storedTasks)
	database = db

	m := Model{
		taskList:  list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 20),
		taskInput: textinput.New(),
	}

	m.taskList.SetShowTitle(false)
	m.taskList.SetItems(taskItems)
	m.taskList.SetShowHelp(false)
	m.taskList.KeyMap.CursorUp.SetKeys("up")
	m.taskList.KeyMap.CursorDown.SetKeys("down")

	m.taskInput.Focus()
	m.taskList.DisableQuitKeybindings()

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
			case "enter":
				_, err := task.AddTask(database, task.Task{Name: m.taskInput.Value()})
				if err != nil {
					fmt.Sprintf("Error adding task: %v", err)
				} else {
					reloadTasks(&m, database)
				}
			}
		} else {
			switch msg.String() {
			case "w":
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.taskList.SetWidth(msg.Width)
		return m, nil
	}

	m.taskList, cmd = m.taskList.Update(msg)
	m.taskInput, _ = m.taskInput.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return m.taskList.View() + "\n" + m.taskInput.View()
}
