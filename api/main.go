package api

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Todo struct {
	UpdateTimestamp   sql.NullTime `json:"update_timestamp"`
	CreationTimestamp sql.NullTime `json:"creation_timestamp"`
	ID                uuid.UUID    `json:"id"`
	Text              string       `json:"text"`
	IsDone            bool         `json:"is_done"`
}

// Init function for API service
func Init(logger *zap.SugaredLogger, db *sql.DB) error {
	router := gin.Default()
	apiGroup := router.Group("/api")

	declareAPIRoutes(logger, apiGroup, db)

	err := router.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("api service initialized")
	return nil
}

func declareAPIRoutes(logger *zap.SugaredLogger, apiGroup *gin.RouterGroup, db *sql.DB) {
	todoRoutes(logger, apiGroup, db)
}

func todoRoutes(logger *zap.SugaredLogger, apiGroup *gin.RouterGroup, db *sql.DB) {
	todosGroup := apiGroup.Group("/todos")
	/**
	 * Create a todo_service that will define the struct and has the
	 */
	todosGroup.GET("/", func(c *gin.Context) {
		todos, err := getTodos(logger, apiGroup, db)
		if err != nil {
			logger.Error(err)
		}
		c.JSON(200, gin.H{
			"data": todos,
		})
	})
}

func getTodos(logger *zap.SugaredLogger, apiGroup *gin.RouterGroup, db *sql.DB) ([]Todo, error) {
	rows, err := sq.Select("*").From("todo").RunWith(db).Query()
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
		logger.Error(err)
	}
	return todos, nil
}
