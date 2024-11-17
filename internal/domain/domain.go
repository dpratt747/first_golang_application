package domain

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserDeletion struct {
	UserId int `json:"userid" binding:"required"`
}
