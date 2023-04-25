package database

import (
	sql "database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	path string
}

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		`
			CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			due_date TEXT,
			should_notify INTEGER,
			topic_id INTEGER,
			is_done INTEGER,
		    FOREIGN KEY (topic_id) REFERENCES topics (id) ON DELETE SET NULL
		)`,
	)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		`
			CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
    )`,
	)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		`
			CREATE TABLE IF NOT EXISTS task_tags (
			task_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			PRIMARY KEY (task_id, tag_id),
			FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
		)`,
	)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		`
		CREATE TABLE IF NOT EXISTS topics (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
