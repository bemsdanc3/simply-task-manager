package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(postgresDSN string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id BIGSERIAL PRIMARY KEY,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT now()
	);`
	if _, err = db.Exec(query); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateTask(contentOfTask string) (int64, error) {
	const op = "storage.postgres.CreateTask"

	stmt, err := s.db.Prepare(`INSERT INTO tasks (content) VALUES ($1) RETURNING id`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64

	err = stmt.QueryRow(contentOfTask).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}

func (s *Storage) UpdateTask(contentToUpdate string, id int) error {
	const op = "storage.postgres.UpdateTask"

	stmt, err := s.db.Prepare(`UPDATE tasks SET content = $1 WHERE id = $2`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(contentToUpdate, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return err
}

func (s *Storage) DeleteTask(id int) error {
	const op = "storage.postgres.DeleteTask"

	stmt, err := s.db.Prepare("DELETE FROM tasks WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return err
}
