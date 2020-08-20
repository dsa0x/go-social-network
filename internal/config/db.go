package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB returns the database object
func DB() *gorm.DB {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)

	db, err := gorm.Open("postgres", psqlInfo)

	db = db.BlockGlobalUpdate(true)

	if err != nil {
		log.Fatal("Database error: ", err)
		panic(err)
	}
	log.Println("Database Connected")

	return db

}
