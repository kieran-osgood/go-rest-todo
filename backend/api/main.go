package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
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

	todosGroup.GET("/", todoService.ListTodos)
	todosGroup.POST("/", todoService.CreateTodo)
	todosGroup.DELETE("/:id", todoService.DeleteTodo)
}
