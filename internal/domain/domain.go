package domain



type User struct {
	ID int
	Username string `json:"username" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}