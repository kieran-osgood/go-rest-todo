package models

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type TodoService struct {
	Logger *zap.SugaredLogger
	Db     *sql.DB
}

type Todo struct {
	UpdateTimestamp   sql.NullTime `json:"update_timestamp"`
	CreationTimestamp sql.NullTime `json:"creation_timestamp"`
	ID                uuid.UUID    `json:"id"`
	Text              string       `json:"text"`
	IsDone            bool         `json:"is_done"`
}

func (t *TodoService) GetTodos() ([]Todo, error) {
	rows, err := sq.Select("*").From("todo").RunWith(t.Db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]Todo, 0)
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.CreationTimestamp, &todo.UpdateTimestamp, &todo.ID, &todo.Text, &todo.IsDone)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		t.Logger.Error(err)
	}
	return todos, nil
}
