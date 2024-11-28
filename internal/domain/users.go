package domain

import "gorm.io/gorm"

// User is a struct representing a user
type User struct {
	gorm.Model
	Name     string
	Password string
}

// UserRepository interface
type UserRepository interface {
	Create(user *User) error
	FindAll() ([]User, error)
	FindByID(id int) (*User, error)
	Update(user *User) error
	Delete(id int) error
}

// UserUseCase interface
type UserUseCase interface {
	CreateUser(user *User) error
	GetAllUsers() ([]User, error)
	GetUserByID(id int) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int) error
}
