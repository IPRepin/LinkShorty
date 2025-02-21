package auth

import (
	"LinkShorty/internal/user"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *user.UserRepository
}

func NewAuthService(userRepo *user.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (as *AuthService) Register(email string, password string, name string) (string, error) {
	existedUser, _ := as.UserRepo.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	newUser := &user.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	_, err = as.UserRepo.Create(newUser)
	if err != nil {
		return "", err
	}
	return newUser.Email, nil
}

func (as *AuthService) Login(email string, password string) (string, error) {
	existedUser, _ := as.UserRepo.FindByEmail(email)
	if existedUser == nil {
		return "", errors.New(ErrWrongCredentials)
	}
	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}
	return existedUser.Email, nil
}
