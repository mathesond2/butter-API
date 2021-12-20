package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func GetDB() *gorm.DB {
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=require password=%s", dbHost, dbPort, username, dbName, password)
	fmt.Println(dbUri, "dave")

	db, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db.Debug().AutoMigrate(&Account{}, &Invoice{})

	return db
}
