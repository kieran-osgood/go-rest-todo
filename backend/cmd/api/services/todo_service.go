package services

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/kieran-osgood/go-rest-todo/cmd/api/errors"
	errorHandler "github.com/kieran-osgood/go-rest-todo/cmd/error"
	"net/http"
	"strings"

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
	UpdateTimestamp   sql.NullTime `json:"updateTimestamp"`
	CreationTimestamp sql.NullTime `json:"creationTimestamp"`
	ID                uuidv4.UUID  `json:"id"`
	Text              string       `json:"text"`
	IsDone            bool         `json:"isDone"`
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
			"error":   errors.ServerError,
			"message": "Database connection failed.",
		})
		return
	}

	defer errorHandler.CleanUpAndHandleError(rows.Close, t.Logger)

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
			"error":   errors.ServerError,
			"message": "Error serializing database rows.",
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

// CreateTodo - retrieves all Todos in the database
func (t *TodoService) CreateTodo(c *gin.Context) {
	var todo Todo
	err := c.BindJSON(&todo)
	if err != nil {
		t.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   errors.BadRequest,
			"message": "Invalid post body.",
		})
		return
	}

	if todo.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   errors.BadRequest,
			"message": "Text is a required field.",
		})
		return
	}

	uuid, err := uuidv4.NewV4()
	if err != nil {
		t.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   errors.ServerError,
			"message": "Couldn't create ID.",
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
			"error":   errors.ServerError,
			"message": "Failed to access database.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": &CreateTodoResponseBody{
			id: uuid,
		},
	})
}

// DeleteTodo - delete a todo via the ID
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
			"error":   errors.ServerError,
			"message": "Failed to access database.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
	return
}

// SearchTodos searches `todo` table using postgresql `trigram`  search algorithm
func (t *TodoService) SearchTodos(c *gin.Context) {
	values := c.Request.URL.Query()
	searchText := values["search"]

	if searchText == nil {
		t.Logger.Error("Search query parameter was empty.")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   errors.BadRequest,
			"message": "Search query parameter cannot be empty.",
		})
		return
	}

	rows, err := sq.
		Select("*").
		From(tableName).
		OrderByClause("SIMILARITY(text ,?) DESC", strings.Join(searchText, "")).
		Limit(5).
		PlaceholderFormat(sq.Dollar).
		RunWith(t.Db).
		Query()

	if err != nil {
		t.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   errors.ServerError,
			"message": "Database connection failed.",
		})
		return
	}

	defer errorHandler.CleanUpAndHandleError(rows.Close, t.Logger)

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
			"error":   errors.ServerError,
			"message": "Error serializing database rows.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": todos,
	})
}