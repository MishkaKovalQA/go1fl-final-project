package db

import (
	"database/sql"
	"errors"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func Tasks(search string, limit int) ([]*Task, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if search == "" {
		rows, err = db.Query(
			`SELECT id, date, title, comment, repeat
			 FROM scheduler
			 ORDER BY date
			 LIMIT ?`,
			limit,
		)
	} else {
		searchDate, parseErr := time.Parse("02.01.2006", search)

		if parseErr == nil {
			rows, err = db.Query(
				`SELECT id, date, title, comment, repeat
				 FROM scheduler
				 WHERE date = ?
				 ORDER BY date
				 LIMIT ?`,
				searchDate.Format("20060102"),
				limit,
			)
		} else {
			pattern := "%" + search + "%"

			rows, err = db.Query(
				`SELECT id, date, title, comment, repeat
				 FROM scheduler
				 WHERE title LIKE ? OR comment LIKE ?
				 ORDER BY date
				 LIMIT ?`,
				pattern,
				pattern,
				limit,
			)
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {
		task := &Task{}

		err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	task := &Task{}

	err := db.QueryRow(
		`SELECT id, date, title, comment, repeat
		 FROM scheduler
		 WHERE id = ?`,
		id,
	).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func UpdateTask(task *Task) error {
	result, err := db.Exec(
		`UPDATE scheduler
		 SET date = ?, title = ?, comment = ?, repeat = ?
		 WHERE id = ?`,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("task not found")
	}

	return nil
}

func DeleteTask(id string) error {
	result, err := db.Exec(
		`DELETE FROM scheduler WHERE id = ?`,
		id,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("task not found")
	}

	return nil
}

func UpdateDate(next, id string) error {
	result, err := db.Exec(
		`UPDATE scheduler SET date = ? WHERE id = ?`,
		next,
		id,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("task not found")
	}

	return nil
}
