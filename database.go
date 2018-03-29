package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=voting_app_db ssl-mode=disable")
	if err != nil {
		log.Println("Error when openning database...", err)
		return
	}
}

func createUser(name, email, password string) error {
	stmt := "INSERT INTO USERS (NAME, EMAIL, PASSWORD) VALUES ($1, $2, $3);"
	_, err := db.Exec(stmt, name, email, password)
	if err != nil {
		log.Println("Error adding new user to users table...", err)
		return err
	}
	log.Println("new user added to the database...")
	return nil
}
