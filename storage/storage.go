package storage

import (
	"context"
	"errors"
	"time"
)

type Storage interface {
	Save(ctx context.Context, content string) error
	Tasks(ctx context.Context) ([]*Task, error)
	Complete(ctx context.Context, id int) error
	Remove(ctx context.Context, id int) error
	IsExists(ctx context.Context, content string) (bool, error)
	IsExistsID(ctx context.Context, id int) (bool, error)
	SaveUser(u *User) error
	RetrieveUser(u *User) (*User, error)
}

var ErrNoSavedTasks = errors.New("There are no saved tasks.")

type Task struct {
	ID      int
	Content string
	//Deadline time.Time
	Created   time.Time
	Completed uint8
}
type User struct {
	Username string
	Chatid   int
}
