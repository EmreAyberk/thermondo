package domain

import (
	"gorm.io/gorm"
	"movie-rating-service/internal/application/models/response"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	IsAdmin  bool   `json:"is_admin"`
}

func (u *User) GetUserResponse() *response.GetUser {
	return &response.GetUser{
		ID:       u.ID,
		Username: u.Username,
		Name:     u.Name,
		Surname:  u.Surname,
		Email:    u.Email,
		Phone:    u.Phone,
		Address:  u.Address,
		IsAdmin:  u.IsAdmin,
	}
}

func (u *User) CreateUserResponse() *response.CreateUser {
	return &response.CreateUser{
		ID: u.ID,
	}
}
