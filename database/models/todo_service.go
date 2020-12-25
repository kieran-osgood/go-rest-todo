package models

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

var tableName = "todo"

// TodoService - Initialiser for service
type TodoService struct {
	Logger *zap.SugaredLogger
	Db     *sql.DB
}

// TodoDatabaseEntity - Base type
type TodoDatabaseEntity struct {
	UpdateTimestamp   sql.NullTime
	CreationTimestamp sql.NullTime
	ID                uuid.UUID
	Text              string
	IsDone            bool
}

// GetTodos - retrieves all Todos in the database
func (t *TodoService) GetTodos() ([]TodoDatabaseEntity, error) {
	rows, err := sq.Select("*").From(tableName).RunWith(t.Db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]TodoDatabaseEntity, 0)
	for rows.Next() {
		var todo TodoDatabaseEntity
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

// TodoPost - Base type
type TodoPost struct {
	Text   string
	IsDone bool
}

// PostTodos - retrieves all Todos in the database
func (t *TodoService) PostTodos(c *gin.Context) (*uuid.UUID, error) {
	var todopostmap TodoPost

	err := c.BindJSON(&todopostmap)
	if err != nil {
		return nil, err
	}
	/**
	 * do some validation that the todopostmap has all the fields required
	 * figure out how to throw 400 BAD REQUEST
	 */
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	err = sq.
		Insert(tableName).
		Columns("id", "is_done", "text").
		Values(uuid, todopostmap.IsDone, todopostmap.Text).
		Suffix("RETURNING \"id\"").
		RunWith(t.Db).
		PlaceholderFormat(sq.Dollar).
		QueryRow().
		Scan(&uuid)

	if err != nil {
		return nil, err
	}
	return &uuid, nil
}
