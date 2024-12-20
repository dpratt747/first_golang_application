package server

import (
	"db_access/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	router := gin.Default()

	router.POST("/user", s.InsertNewUserHandler)

	router.GET("/users", s.GetAllUsersHandler)

	router.DELETE("/user/:userId", s.DeleteUserHandler)

	return router
}

func (s *Server) GetAllUsersHandler(c *gin.Context) {
	users, err := s.Db.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{})
		return
	} else {
		c.JSON(http.StatusOK, users)
		return
	}
}

func (s *Server) DeleteUserHandler(c *gin.Context) {
	userIdParam := c.Param("userId")

	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId format. Must be an integer."})
		return
	}

	err = s.Db.SoftDeleteUser(userId)
	switch err.(type) {
	case *domain.UniqueConstraintDatabaseError:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to delete this user as they have already been deleted"})
		return
	case *domain.UserNotFoundError:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to delete this user as they do not exist"})
		return
	default:
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot insert user as this email is already used"})
		return
	default:
		c.JSON(http.StatusCreated, gin.H{"userId": userId})
		return
	}
}
