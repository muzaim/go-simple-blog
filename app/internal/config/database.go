package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "abc123")
	dbHost := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "appdb")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)

	var db *gorm.DB
	var err error

	for i := 1; i <= 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Println("Connected to database successfully!")
			return db
		}
		log.Printf("Waiting for database to be ready (attempt %d/10)... error: %v", i, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("Failed to connect to database after 10 attempts: %v", err)
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

