package Db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectToDb() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load env")
	}

	dbhost := os.Getenv("dbHost")
	dbuser := os.Getenv("dbUser")
	dbpassword := os.Getenv("dbPassword")
	dbname := os.Getenv("dbName")
	dbport := os.Getenv("dbPort")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable port=%s", dbhost, dbuser, dbpassword, dbname, dbport)
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Db not connected")
	}

}
