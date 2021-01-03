package api

import (
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kieran-osgood/go-rest-todo/cmd/api/services"
	"go.uber.org/zap"
	"time"
)

// Init function for API service
func Init(logger *zap.SugaredLogger, db *sql.DB) error {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:3000", "*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Content-Length"},
		//AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "https://github.com"
		//},
		MaxAge: 12 * time.Hour,
	}))

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
	todoService := &services.TodoService{
		Db:     db,
		Logger: logger,
	}

	todosGroup.GET("/search", todoService.SearchTodos)
	todosGroup.GET("", todoService.ListTodos)
	todosGroup.POST("", todoService.CreateTodo)
	todosGroup.DELETE("/:id", todoService.DeleteTodo)
}
