package main

import (
	"LinkShorty/internal/auth"
	"LinkShorty/internal/user"
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func initDB() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err) // Лучше заменить на логирование и возврат ошибки
	}
	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "test@example.com",
		Password: "$2a$10$mE5FwbU/E4of3HXGpFxVJ.sxCg3kz1C6ubaFfirlFYqq1vDuHiRHa",
		Name:     "test",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().Where("email = ?", "test@example.com").Delete(&user.User{})
}

func TestLoginSuccess(t *testing.T) {
	db := initDB()
	initData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	})
	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Expected 200 got %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var resp auth.LoginResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Token == "" {
		t.Fatalf("Expected non-empty token")
	}
	removeData(db)
}

func TestLoginFail(t *testing.T) {
	db := initDB()
	initData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@example.com",
		Password: "passwordFail",
	})
	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 401 {
		t.Fatalf("Expected 401 got %d", res.StatusCode)
	}
	removeData(db)
}
