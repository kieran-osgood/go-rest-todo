package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Init function for API service
func Init(logger *zap.SugaredLogger) error {
	r := gin.Default()
	declareRoutes(r)

	err := r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("api service initialized")
	return nil
}

func declareRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
