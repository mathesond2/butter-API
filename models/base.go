package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func GetDB() *gorm.DB {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	var sslMode string
	if os.Getenv("db_host") == "localhost" {
		sslMode = "disable"
	} else {
		sslMode = "require"
	}

	dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", dbHost, dbPort, username, dbName, sslMode, password)
	db, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db.Debug().AutoMigrate(&Account{}, &Invoice{}, &Webhook{}, &Address{})

	return db
}
