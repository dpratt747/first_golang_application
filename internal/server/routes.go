package server

import (
	"db_access/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.POST("/user", s.InsertNewUserHandler)

	return r
}

func (s *Server) InsertNewUserHandler(c *gin.Context) {
	
	var newUser domain.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	userId := s.Db.InsertNewUser(newUser)

	c.JSON(http.StatusOK, userId)
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.Db.Health())
}
