package tag

import (
	"database/sql"
)

type Tag struct {
	ID   int64
	Name string
}

func AddTag(db *sql.DB, t Tag) (int64, error) {
	res, err := db.Exec(
		"INSERT INTO tags(name) VALUES(?)",
		t.Name,
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

func UpdateTag(db *sql.DB, t Tag) (int64, error) {
	query := `
		UPDATE tags
		SET name = ?
		WHERE id = ?;
	`

	_, err := db.Exec(query, t.Name, t.ID)
	if err != nil {
		return 0, err
	}

	return t.ID, err
}

func DeleteTag(db *sql.DB, id int64) (int64, error) {
	query := `
		DELETE FROM tags
		WHERE id = ?;
	`

	_, err := db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
