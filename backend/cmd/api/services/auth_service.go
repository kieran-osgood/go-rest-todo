package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	apiErrors "github.com/kieran-osgood/go-rest-todo/cmd/api/errors"
	"github.com/kieran-osgood/go-rest-todo/cmd/config"
	"net/http"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

//var tableName = "todo"
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ListTodos - retrieves all Todos in the database
func (s *Service) Login(c *gin.Context) {
	var user User
	err := c.BindJSON(&user)
	if err != nil {
		s.AbortBadPostBody(c, err)
		return
	}
	expectedPassword, ok := users[user.Username]
	if !ok || expectedPassword != user.Password {
		s.AbortUnauthorized(c, errors.New(apiErrors.Unauthorized), "Incorrect Username/Password")
		return
	}

	expirationTime := time.Now().Add(2 * time.Minute).Unix()
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime,
		},
	}
	cfg := config.New()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JwtSecret))

	if err != nil {
		s.AbortServerError(c, err, "Error creating JWT")
		return
	}

	c.SetCookie(
		"access_token",
		tokenString,
		int(expirationTime),
		"/", "localhost",
		false, true,
	)
	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenString,
		"expires":      expirationTime,
	})
}
