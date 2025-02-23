package main

import (
	"os"
	"restapi/internal/link"
	"restapi/internal/stat"
	"restapi/internal/user"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
		//DisableForeignKeyConstraintWhenMigrating: true, //временно игнорировать миграции в первый раз а потом их добавить
	})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&link.Link{}, &user.User{}, &stat.Stat{})
}
