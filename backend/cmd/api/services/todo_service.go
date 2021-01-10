package services

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	apiErrors "github.com/kieran-osgood/go-rest-todo/cmd/api/errors"
	errorHandler "github.com/kieran-osgood/go-rest-todo/cmd/error"
	"net/http"
	"strings"

	uuidv4 "github.com/gofrs/uuid"
)

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
func (s *Service) ListTodos(c *gin.Context) {
	rows, err := sq.
		Select("*").
		From(tableName).
		RunWith(s.Db).
		Limit(10).
		Query()

	if err != nil {
		s.AbortDbConnError(c, err)
	}

	defer errorHandler.CleanUpAndHandleError(rows.Close, s.Logger)

	todos := make([]Todo, 0)
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.CreationTimestamp, &todo.UpdateTimestamp, &todo.ID, &todo.Text, &todo.IsDone)
		if err != nil {
			s.Logger.Error(err)
			return
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		s.AbortSerializationError(c, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": todos,
	})
}

type CreateTodoResponseBody struct {
	id uuidv4.UUID
}

// CreateTodo - retrieves all Todos in the database
func (s *Service) CreateTodo(c *gin.Context) {
	var todo Todo
	err := c.BindJSON(&todo)
	if err != nil {
		s.AbortBadPostBody(c, err)
		return
	}

	if todo.Text == "" {
		s.AbortBadRequest(c, errors.New(apiErrors.BadRequest), "Text is a required field.")
		return
	}

	uuid, err := uuidv4.NewV4()
	if err != nil {
		s.AbortServerError(c, errors.New(apiErrors.ServerError), "Couldn't create ID.")
		return
	}

	err = sq.
		Insert(tableName).
		Columns("id", "is_done", "text").
		Values(uuid, todo.IsDone, todo.Text).
		Suffix("RETURNING \"id\"").
		RunWith(s.Db).
		PlaceholderFormat(sq.Dollar).
		QueryRow().
		Scan(&uuid)

	if err != nil {
		s.AbortDbConnError(c, errors.New(apiErrors.ServerError))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": &CreateTodoResponseBody{
			id: uuid, // TODO Still not returning this
		},
	})
}

// DeleteTodo - delete a todo via the ID
func (s *Service) DeleteTodo(c *gin.Context) {
	id := c.Param("id")

	_, err := sq.
		Delete(tableName).
		Where(sq.Eq{"id": id}).
		RunWith(s.Db).
		PlaceholderFormat(sq.Dollar).
		Exec()

	if err != nil {
		s.AbortDbConnError(c, errors.New(apiErrors.ServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// SearchTodos searches `todo` table using postgresql `trigram`  search algorithm
func (s *Service) SearchTodos(c *gin.Context) {
	values := c.Request.URL.Query()
	searchText := values["search"]

	if searchText == nil {
		s.AbortBadRequest(c, errors.New(apiErrors.BadRequest), "Search query parameter cannot be empty.")
		return
	}

	rows, err := sq.
		Select("*").
		From(tableName).
		OrderByClause("SIMILARITY(text ,?) DESC", strings.Join(searchText, "")).
		Limit(5).
		PlaceholderFormat(sq.Dollar).
		RunWith(s.Db).
		Query()

	if err != nil {
		s.AbortDbConnError(c, errors.New(apiErrors.ServerError))
		return
	}

	defer errorHandler.CleanUpAndHandleError(rows.Close, s.Logger)

	todos := make([]Todo, 0)
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.CreationTimestamp, &todo.UpdateTimestamp, &todo.ID, &todo.Text, &todo.IsDone)
		if err != nil {
			s.Logger.Error(err)
			return
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		s.AbortServerError(c, errors.New(apiErrors.ServerError), "Error serializing database rows.")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": todos,
	})
}
