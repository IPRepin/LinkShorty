package main

import (
	"LinkShorty/internal/link"
	"LinkShorty/internal/stat"
	"LinkShorty/internal/user"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
		//DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err) // Лучше заменить на логирование и возврат ошибки
	}
	err = db.AutoMigrate(&link.Link{}, &user.User{}, &stat.Stat{})
	if err != nil {
		panic(err)
	}
}
