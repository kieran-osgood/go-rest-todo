package api

import (
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kieran-osgood/go-rest-todo/cmd/api/services"
	"go.uber.org/zap"
	"time"
)

// New function for API service
func New(logger *zap.SugaredLogger, db *sql.DB) error {
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
	authRoutes(logger, apiGroup, db)
}

func todoRoutes(logger *zap.SugaredLogger, apiGroup *gin.RouterGroup, db *sql.DB) {
	todosGroup := apiGroup.Group("/todos")
	service := &services.Service{
		Db:     db,
		Logger: logger,
	}

	todosGroup.GET("/search", service.SearchTodos)
	todosGroup.GET("", service.ListTodos)
	todosGroup.POST("", service.CreateTodo)
	todosGroup.DELETE("/:id", service.DeleteTodo)
}

func authRoutes(logger *zap.SugaredLogger, apiGroup *gin.RouterGroup, db *sql.DB) {
	authGroup := apiGroup.Group("/login")
	service := &services.Service{
		Db:     db,
		Logger: logger,
	}

	authGroup.POST("", service.Login)
}
