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
func (t *TodoService) ListTodos(c *gin.Context) {
	rows, err := sq.
		Select("*").
		From(tableName).
		RunWith(t.Db).
		Query()

	if err != nil {
		t.Logger.Error(err)
		c.JSON(500, gin.H{
			"success": false,
			"error":   "Database connection failed.",
		})
		return
	}
	defer rows.Close()

	todos := make([]TodoDatabaseEntity, 0)
	for rows.Next() {
		var todo TodoDatabaseEntity
		err := rows.Scan(&todo.CreationTimestamp, &todo.UpdateTimestamp, &todo.ID, &todo.Text, &todo.IsDone)
		if err != nil {
			t.Logger.Error(err)
			return
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		t.Logger.Error(err)
		c.JSON(500, gin.H{
			"success": false,
			"error":   "Error serializing database rows.",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": todos,
	})
}

// TodoPost - Base type
type TodoPost struct {
	Text   string
	IsDone bool
}

type PostRes struct {
	id uuid.UUID
}

// PostTodos - retrieves all Todos in the database
func (t *TodoService) CreateTodo(c *gin.Context) {
	var todopostmap TodoPost
	err := c.BindJSON(&todopostmap)
	if err != nil {
		t.Logger.Error(err)
		return
	}

	/**
	 * do some validation that the todopostmap has all the fields required
	 * figure out how to throw 400 BAD REQUEST
	 */
	uuidv, err := uuid.NewV4()
	if err != nil {
		t.Logger.Error(err)
		c.JSON(500, gin.H{
			"success": false,
			"error":   "Failed to add todo.",
		})
		return
	}

	err = sq.
		Insert(tableName).
		Columns("id", "is_done", "text").
		Values(uuidv, todopostmap.IsDone, todopostmap.Text).
		Suffix("RETURNING \"id\"").
		RunWith(t.Db).
		PlaceholderFormat(sq.Dollar).
		QueryRow().
		Scan(&uuidv)

	if err != nil {
		t.Logger.Error(err)
		c.JSON(500, gin.H{
			"success": false,
			"error":   "Failed to add todo.",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data": PostRes{
			id: uuidv,
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
		c.JSON(500, gin.H{
			"success": false,
			"error":   "Failed to delete.",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
	return
}
