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

	r.GET("/users", s.GetAllUsersHandler)

	return r
}

func (s *Server) GetAllUsersHandler(c *gin.Context) {
	users, err :=s.Db.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	
	c.JSON(http.StatusOK, users)
}

func (s *Server) InsertNewUserHandler(c *gin.Context) {
	
	var newUser domain.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	userId, err := s.Db.InsertNewUser(newUser)
	switch err.(type) {
		case *domain.UniqueConstraintDatabaseError:
			c.JSON(http.StatusInternalServerError, err)
			return
		default:
			c.JSON(http.StatusCreated, userId)
			return
	}
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.Db.Health())
}
