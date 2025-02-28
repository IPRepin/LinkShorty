package auth_test

import (
	"LinkShorty/configs"
	"LinkShorty/internal/auth"
	"LinkShorty/internal/user"
	"LinkShorty/pkg/db"
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func bootstrap() (*auth.AuthHandler, sqlmock.Sqlmock, error) {
	base, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: base,
	}))
	if err != nil {
		return nil, nil, err
	}
	userRepo := user.NewUserRepository(&db.Db{
		DB: gormDB,
	})
	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: auth.NewAuthService(userRepo),
	}
	return &handler, mock, nil
}

func TestRegisterHandlerSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"email", "password"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	mock.ExpectCommit()

	if err != nil {
		t.Fatal(err)
		return
	}
	data, _ := json.Marshal(&auth.RegisterRequest{
		Email:    "test@example.com",
		Password: "password",
		Name:     "test",
	})
	reader := bytes.NewReader(data)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/auth/register", reader)
	if err != nil {
		t.Fatal(err)
		return
	}
	handler.Register()(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("w.Code = %d, want %d", w.Code, http.StatusCreated)
		return
	}

}

func TestLoginHandlerSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow("test@example.com", "$2a$10$mE5FwbU/E4of3HXGpFxVJ.sxCg3kz1C6ubaFfirlFYqq1vDuHiRHa")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	if err != nil {
		t.Fatal(err)
		return
	}
	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	})
	reader := bytes.NewReader(data)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/auth/login", reader)
	if err != nil {
		t.Fatal(err)
		return
	}
	handler.Login()(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("w.Code = %d, want %d", w.Code, http.StatusOK)
		return
	}
}
