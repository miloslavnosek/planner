package task_list

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"planner/task"
)

type (
	Model struct {
		List list.Model
	}
	TaskItem task.Task
)

var database *sql.DB

func (i TaskItem) Title() string       { return i.Name }
func (i TaskItem) Description() string { return i.Desc }
func (i TaskItem) FilterValue() string { return i.Name }
func (i TaskItem) getTask() task.Task  { return task.Task(i) }

func newTaskListDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = selectedTitleStyle
	d.Styles.SelectedDesc = selectedDescStyle

	d.Styles.FilterMatch = filterMatchStyle

	d.Styles.NormalTitle = normalTitleStyle
	d.Styles.NormalDesc = normalDescStyle

	return d
}

func SetItems(m *Model, tasks []task.Task) {
	taskItems := ConvertTasksToItems(tasks)

	m.List.SetItems(taskItems)
}

func GetCurrentTask(m *Model) task.Task {
	currentItem := m.List.Items()[m.List.Cursor()].(TaskItem)

	return currentItem.getTask()
}

func ConvertTasksToItems(tasks []task.Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, t := range tasks {
		items[i] = TaskItem(t)
	}
	return items
}

func ReloadTasks(m *Model, database *sql.DB) {
	storedTasks, err := task.GetTasks(database)
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}
	taskItems := ConvertTasksToItems(storedTasks)
	m.List.SetItems(taskItems)
}

func ToggleTaskCompleted(m *Model, database *sql.DB, selectedTask task.Task) {
	selectedTask.IsDone = !selectedTask.IsDone
	_, err := task.UpdateTask(database, selectedTask)
	if err != nil {
		fmt.Printf("Error toggling task completed: %v\n", err)
		return
	}
	ReloadTasks(m, database)
}

func InitialModel(db *sql.DB) Model {
	database = db

	l := list.New([]list.Item{}, newTaskListDelegate(), 80, 20)

	m := Model{
		List: l,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.List, cmd = m.List.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return m.List.View()
}
