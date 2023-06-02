package task

import (
	"database/sql"
)

type Task struct {
	ID           int64
	Name         string
	Desc         string
	DueDate      string
	TopicId      int64
	ShouldNotify bool
	IsDone       bool
}

func GetTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT id, name, description, due_date, should_notify, is_done FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		var dueDate sql.NullString
		err := rows.Scan(&t.ID, &t.Name, &t.Desc, &dueDate, &t.ShouldNotify, &t.IsDone)
		if err != nil {
			return nil, err
		}

		if dueDate.Valid {
			t.DueDate = dueDate.String
		}

		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func AddTask(db *sql.DB, t Task) (int64, error) {
	res, err := db.Exec(
		"INSERT INTO tasks(name, description, due_date, topic_id, should_notify) VALUES(?, ?, ? ,?, ?)",
		t.Name,
		t.Desc,
		t.DueDate,
		t.TopicId,
		t.ShouldNotify,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateTask(db *sql.DB, t Task) (int64, error) {
	query := `
		UPDATE tasks
		SET name = ?, description = ?, due_date = ?, topic_id = ?, should_notify = ?, is_done = ?
		WHERE id = ?;
	`

	_, err := db.Exec(query, t.Name, t.Desc, t.DueDate, t.TopicId, t.ShouldNotify, t.IsDone, t.ID)
	if err != nil {
		return 0, err
	}

	return t.ID, err
}

func DeleteTask(db *sql.DB, id int64) (int64, error) {
	query := `
		DELETE FROM tasks
		WHERE id = ?;
	`

	_, err := db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
