package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var ErrNoMatch = fmt.Errorf("no matching record")

type Database struct {
	Conn *sql.DB
}

func InitDB() Database {
	database, err := Initialize(os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}

	return database
}

func Initialize(DbUser, DbPassword, DbPort, DbHost, DbName string) (Database, error) {
	db := Database{}
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		DbHost, DbPort, DbUser, DbName, DbPassword)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, err
	}

	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	log.Println("Database connection established")
	return db, nil
}

func DbInitTest() Database {
	db, err := Initialize(os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_NAME"))
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}

	return db
}