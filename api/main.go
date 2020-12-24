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
	declareRoutes(logger, router, db)

	err := router.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("api service initialized")
	return nil
}

func declareRoutes(logger *zap.SugaredLogger, router *gin.Engine, db *sql.DB) {
	router.GET("/api/todos", func(c *gin.Context) {
		rows, err := sq.Select("*").From("todo").RunWith(db).Query()
		if err != nil {
			logger.Error(err)
		}
		defer rows.Close()

		todos := make([]Todo, 0)
		for rows.Next() {
			var todo Todo
			err := rows.Scan(&todo.CreationTimestamp, &todo.UpdateTimestamp, &todo.ID, &todo.Text, &todo.IsDone)
			if err != nil {
				logger.Error(err)
			}
			logger.Debug(todo.ID)
			todos = append(todos, todo)
		}

		if err = rows.Err(); err != nil {
			logger.Fatal(err)
		}
		// m, _ := json.Marshal(&todos)
		c.JSON(200, todos)
	})
}
