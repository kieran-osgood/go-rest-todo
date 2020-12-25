package models

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"net/http"

	uuidv4 "github.com/gofrs/uuid"
	"go.uber.org/zap"
)

// TodoService - Initializer for service
type TodoService struct {
	Logger *zap.SugaredLogger
	Db     *sql.DB
}

var tableName = "todo"

// Todo - Full model representative of the database shape
type Todo struct {
	UpdateTimestamp   sql.NullTime
	CreationTimestamp sql.NullTime
	ID                uuidv4.UUID
	Text              string
	IsDone            bool
}

// ListTodos - retrieves all Todos in the database
func (t *TodoService) ListTodos(c *gin.Context) {
	rows, err := sq.
		Select("*").
		From(tableName).
		RunWith(t.Db).
		Limit(10).
		Query()

	if err != nil {
		t.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Database connection failed.",
		})
		return
	}
	defer rows.Close()

	todos := make([]Todo, 0)
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.CreationTimestamp, &todo.UpdateTimestamp, &todo.ID, &todo.Text, &todo.IsDone)
		if err != nil {
			t.Logger.Error(err)
			return
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		t.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error serializing database rows.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": todos,
	})
}

type CreateTodoResponseBody struct {
	id uuidv4.UUID
}

// PostTodos - retrieves all Todos in the database
func (t *TodoService) CreateTodo(c *gin.Context) {
	var todo Todo
	err := c.BindJSON(&todo)
	if err != nil {
		t.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid post body.",
		})
		return
	}

	if todo.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Text is a required field.",
		})
		return
	}

	uuid, err := uuidv4.NewV4()
	if err != nil {
		t.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to add todo.",
		})
		return
	}

	err = sq.
		Insert(tableName).
		Columns("id", "is_done", "text").
		Values(uuid, todo.IsDone, todo.Text).
		Suffix("RETURNING \"id\"").
		RunWith(t.Db).
		PlaceholderFormat(sq.Dollar).
		QueryRow().
		Scan(&uuid)

	if err != nil {
		t.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to add todo.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": CreateTodoResponseBody{
			id: uuid,
		},
	})
}

// DeleteTodo - retrieves all Todos in the database
func (t *TodoService) DeleteTodo(c *gin.Context) {
	id := c.Param("id")

	_, err := sq.
		Delete(tableName).
		Where(sq.Eq{"id": id}).
		RunWith(t.Db).
		PlaceholderFormat(sq.Dollar).
		Exec()

	if err != nil {
		t.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
	return
}
