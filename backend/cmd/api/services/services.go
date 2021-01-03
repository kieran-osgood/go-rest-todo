package services

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/kieran-osgood/go-rest-todo/cmd/api/errors"
	"go.uber.org/zap"
	"net/http"
)

// Service - Initializer for service
type Service struct {
	Logger *zap.SugaredLogger
	Db     *sql.DB
}

func (s *Service) AbortDbConnError(c *gin.Context, err error) {
	if err != nil {
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   errors.ServerError,
			"message": "Database connection failed.",
		})
	}
}

func (s *Service) AbortSerializationError(c *gin.Context, err error) {
	if err != nil {
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   errors.ServerError,
			"message": "Error serializing database rows.",
		})
	}
}

func (s *Service) AbortBadPostBody(c *gin.Context, err error) {
	if err != nil {
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   errors.BadRequest,
			"message": "Invalid post body.",
		})
	}
}
