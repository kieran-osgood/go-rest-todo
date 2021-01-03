package services

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// Service - Initializer for service
type Service struct {
	Logger *zap.SugaredLogger
	Db     *sql.DB
}

func ErrorBody(success bool, err error, message string) gin.H {
	return gin.H{
		"success": success,
		"error":   err,
		"message": message,
	}
}

// Customisable default Handlers
func (s *Service) AbortBadRequest(c *gin.Context, err error, message string) {
	if err != nil {
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorBody(false, err, message))
	}
}

func (s *Service) AbortServerError(c *gin.Context, err error, message string) {
	if err != nil {
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorBody(false, err, message))
	}
}

// Reusable scenario specific handlers
func (s *Service) AbortDbConnError(c *gin.Context, err error) {
	s.AbortServerError(c, err, "Database connection failed.")
}

func (s *Service) AbortSerializationError(c *gin.Context, err error) {
	s.AbortServerError(c, err, "Error serializing database rows.")
}

func (s *Service) AbortBadPostBody(c *gin.Context, err error) {
	s.AbortBadRequest(c, err, "Invalid post body.")
}
