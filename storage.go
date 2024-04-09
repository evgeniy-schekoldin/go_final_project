package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type TaskStore struct {
	db *sql.DB
}

func NewTaskStore(db *sql.DB) TaskStore {
	return TaskStore{db: db}
}

func (s TaskStore) getTask(id int) (Task, error) {
	var task Task
	row := s.db.QueryRow("SELECT * FROM scheduler WHERE id = :id LIMIT 50", sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		log.Print("Not found")
		return task, err
	}
	return task, nil
}

func (s TaskStore) getTasks(search string) ([]Task, error) {
	var tasks []Task
	var query string

	switch {
	case search != "":
		if date, err := time.Parse("02.01.2006", search); err == nil {
			search = date.Format(dateLayout)
			query = "SELECT * FROM scheduler WHERE date = :search LIMIT 50"
		} else {
			query = "SELECT * FROM scheduler WHERE title LIKE concat('%', :search, '%') OR comment LIKE concat('%', :search, '%') ORDER BY date LIMIT 50"
		}

	default:
		query = "SELECT * FROM scheduler ORDER BY date LIMIT 50"

	}

	rows, err := s.db.Query(query, sql.Named("search", search))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		tasks = append(tasks, task)

		fmt.Println(tasks)

	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	if tasks == nil {
		tasks = make([]Task, 0)
	}

	return tasks, nil
}

func (s TaskStore) update(t Task) error {
	_, err := s.db.Exec("UPDATE scheduler SET date=:date , title=:title, comment=:comment, repeat=:repeat WHERE id=:id",
		sql.Named("id", t.ID),
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))

	if err != nil {
		return err
	}

	return nil

}

func (s TaskStore) delete(id int) error {
	_, err := s.db.Exec("DELETE FROM scheduler WHERE id=:id",
		sql.Named("id", id))

	if err != nil {
		return err
	}

	return nil

}

func (s TaskStore) add(t Task) (int, error) {
	result, err := s.db.Exec("INSERT INTO scheduler (date , title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	return int(id), err
}
