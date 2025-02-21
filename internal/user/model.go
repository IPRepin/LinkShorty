package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"index"`
	Password string
	Name     string
}

func NewUser(email string, password string, name string) *User {
	return &User{
		Email:    email,
		Password: password,
		Name:     name,
	}
}
