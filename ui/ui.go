package ui

import (
	"database/sql"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"planner/task"
	"strconv"
)

type Model struct {
	taskList list.Model
}

type taskItem task.Task

func (i taskItem) Title() string       { return i.Name }
func (i taskItem) Description() string { return i.Desc }
func (i taskItem) FilterValue() string { return strconv.FormatInt(i.ID, 10) }

func ConvertTasksToItems(tasks []task.Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, t := range tasks {
		items[i] = taskItem(t)
	}
	return items
}

func InitialModel(db *sql.DB) Model {
	storedTasks, _ := task.GetTasks(db)
	taskItems := ConvertTasksToItems(storedTasks)

	m := Model{
		taskList: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
	}
	m.taskList.SetShowTitle(false)
	m.taskList.SetItems(taskItems)

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.taskList.SetWidth(msg.Width)
		return m, nil
	}

	var cmd tea.Cmd
	m.taskList, cmd = m.taskList.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.taskList.View()
}
