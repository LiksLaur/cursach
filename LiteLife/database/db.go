package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := "user=postgres dbname=litelifedb password=Pgadmin port=5432 sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(100) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS repair_requests (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			username VARCHAR(50) NOT NULL,
			apartment VARCHAR(20) NOT NULL,
			repair_type VARCHAR(50) NOT NULL,
			comment TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			is_approved BOOLEAN DEFAULT FALSE
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS chat_messages (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL,
			message TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS room_bookings (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			phone VARCHAR(20) NOT NULL,
			room_number INTEGER NOT NULL,
			booking_date DATE NOT NULL,
			is_approved BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}
