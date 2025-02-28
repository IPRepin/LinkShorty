package auth_test

import (
	"LinkShorty/internal/auth"
	"LinkShorty/internal/user"
	"testing"
)

type MockIUserRepository struct {
}

func (repo *MockIUserRepository) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: "test@test.com",
	}, nil
}

func (repo *MockIUserRepository) FindByEmail(email string) (*user.User, error) {
	return nil, nil
}

func TestRegisterSuccess(t *testing.T) {
	const initialEmail = "test@test.com"
	authService := auth.NewAuthService(&MockIUserRepository{})
	email, err := authService.Register(initialEmail, "test", "test")
	if err != nil {
		t.Fatalf("Register error = %v", err)
	}
	if email != initialEmail {
		t.Fatalf("email = %s, want %s", email, initialEmail)
	}
}
