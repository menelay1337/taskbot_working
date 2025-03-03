package mysql

import (
	"context"
	"database/sql"
	"fmt"


	_ "github.com/go-sql-driver/mysql"

	"taskbot1/storage"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Storage{db: db}, nil

}

func (s *Storage) Save(ctx context.Context, content string, chatID int) error {
	stmt := "INSERT INTO tasks (content, created, chatid) VALUES (?, CONVERT_TZ(UTC_TIMESTAMP(), 'UTC', 'Asia/Almaty'), ?)"

	if _, err := s.db.ExecContext(ctx, stmt, content, chatID); err != nil {
		return fmt.Errorf("Can't save page: %w", err)
	}

	return nil
}

func (s *Storage) Tasks(ctx context.Context, chatID int) ([]*storage.Task, error) {
	stmt := `SELECT id, content, created, deadline, completed FROM tasks WHERE chatid = ?`

	rows, err := s.db.QueryContext(ctx, stmt, chatID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []*storage.Task

	for rows.Next() {
		t := &storage.Task{}

		err = rows.Scan(&t.ID, &t.Content, &t.Created, &t.Deadline, &t.Completed)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, storage.ErrNoSavedTasks
	}

	return tasks, nil
}

//func (s *Storage) PastTasks() ( []*storage.Task, error ) {
//	stmt := `SELECT content, deadline, created FROM tasks
//	WHERE deadline < UTC_TIMESTAMP()`
//
//	rows, err := s.db.Query(stmt)
//	if err != nil {
//		return nil, err
//	}
//
//	defer rows.Close()
//
//	var tasks []*storage.Task
//
//	for rows.Next() {
//		var t task
//
//		err = rows.Scan(&t.Header, &t.Content, &t,Deadline, &t.Created)
//		if err != nil {
//			return nil, err
//		}
//
//		tasks = append(tasks, task)
//	}
//
//	if err = rows.Err(); err != nil {
//		return nil, err
//	}
//
//	if len(tasks) == 0 {
//		return nil, storage.ErrNoPastTasks
//	}
//
//	return tasks, nil
//}

func (s *Storage) Remove(ctx context.Context, id int) error {
	stmt := "DELETE FROM tasks WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, stmt, id); err != nil {
		return fmt.Errorf("Can't remove task: %w", err)
	}

	return nil
}

func (s *Storage) Complete(ctx context.Context, id int) error {
	stmt := "UPDATE tasks SET completed = 1 WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, stmt, id); err != nil {
		return fmt.Errorf("Can't complete task: %w", err)
	}

	return nil
}

func (s *Storage) Deadline(ctx context.Context, id int, days int) error {
	stmt := "UPDATE tasks SET deadline = DATE_ADD(created, INTERVAL ? DAY) WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, stmt, days, id); err != nil {
		return fmt.Errorf("Can't complete task: %w", err)
	}

	return nil
}

//func (s *Storage) Clear() error {
//	stmt := "DELETE FROM tasks WHERE deadline < UTC_TIMESTAMP()"
//
//	result, err := s.db.Exec(stmt)
//
//	if err != nil {
//		return err
//	}
//
//	rowsAffected, err := result.RowsAffected()
//	if err != nil {
//		return err
//	}
//	if rowsAffected == 0 {
//		return fmt.Errorf("Tasks wasn't removed")
//	}
//
//	return nil
//}

func (s *Storage) IsExists(ctx context.Context, content string) (bool, error) {
	stmt := "SELECT COUNT(*) FROM tasks where content = ?"

	var count int

	if err := s.db.QueryRowContext(ctx, stmt, content).Scan(&count); err != nil {
		return false, fmt.Errorf("Can't check if task exists: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) IsExistsID(ctx context.Context, id int) (bool, error) {
	stmt := "SELECT COUNT(*) FROM tasks where id = ?"

	var count int

	if err := s.db.QueryRowContext(ctx, stmt, id).Scan(&count); err != nil {
		return false, fmt.Errorf("Can't check if task exists: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Init() error {
	stmt := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INT AUTO_INCREMENT PRIMARY KEY,
		chatid INT NOT NULL,
		content VARCHAR(255) UNIQUE NOT NULL,
		created TIMESTAMP DEFAULT CONVERT_TZ(UTC_TIMESTAMP(), 'UTC', 'Asia/Almaty'),
		deadline TIMESTAMP DEFAULT "0000-00-00 00:00:00",
		completed BOOLEAN DEFAULT FALSE
	);
	`
	_, err := s.db.Exec(stmt)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Can't create table: %w", err)
	}

	return nil
}
