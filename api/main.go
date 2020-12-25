package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/kieran-osgood/go-rest-todo/database/models"
	"go.uber.org/zap"
)

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

	todoService := &models.TodoService{
		Db:     db,
		Logger: logger,
	}

	todosGroup.GET("/", func(c *gin.Context) {
		todos, err := todoService.GetTodos()
		if err != nil {
			logger.Error(err)

		}
		c.JSON(200, gin.H{
			"data": todos,
		})
	})

	type PostRes struct {
		id uuid.UUID
	}

	todosGroup.POST("/add", func(c *gin.Context) {
		id, err := todoService.PostTodos(c)
		if err != nil {
			logger.Error(err)
			c.JSON(500, gin.H{
				"success": false,
				"error":   "Failed to add todo.",
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"data": PostRes{
				id: *id,
			},
		})
	})

	todosGroup.DELETE("/:id", todoService.DeleteTodo)
}
