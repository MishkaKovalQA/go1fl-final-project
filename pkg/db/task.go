package db

import (
	"database/sql"
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

	return tasks, rows.Err()
}
