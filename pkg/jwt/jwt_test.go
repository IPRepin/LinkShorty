package jwt_test

import (
	"LinkShorty/pkg/jwt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func init() {
	// Загружаем переменные из .env файла
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func TestJWTCreate(t *testing.T) {
	const email = "test@example.com"
	jwtService := jwt.NewJWT(os.Getenv("KEY"))
	token, err := jwtService.CreateToken(jwt.JWTData{
		Email: email,
	})
	if err != nil {
		t.Fatal(err)
	}
	isValid, data := jwtService.ParseToken(token)
	if !isValid {
		t.Fatal("Token is not valid")
	}
	if data.Email != email {
		t.Fatalf("Email %s != %s", data.Email, email)
	}
}
