package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=voting_app_db sslmode=disable port=5000")
	if err != nil {
		log.Println("Error when openning database...", err)
		return
	}
}

func (user *User) createUser() error {
	if user.isEmailInDatabase() {
		return errors.New("This email address is already in use.")
	}
	stmt := "INSERT INTO USERS (NAME, EMAIL, PASSWORD, PASSWORD_HASH) VALUES ($1, $2, $3, $4);"
	_, err := db.Exec(stmt, user.Name, user.Email, user.Password, user.hashPassword())
	if err != nil {
		log.Println("Error adding new user to users table...", err)
		return err
	}
	log.Println("new user added to the database...")
	return nil
}

func (user *User) isEmailInDatabase() bool {
	var u = User{}
	stmt := "SELECT * FROM USERS WHERE EMAIL=$1;"
	row := db.QueryRow(stmt, user.Email)
	row.Scan(&user.id, &user.Name, &u.Email, &u.Password, &user.hash)
	if u.Email == user.Email {
		fmt.Println("user email is in the database", u)
		return true
	}
	fmt.Println("email does not exist in the database")
	return false
}

func (user *User) authentification() bool {
	var u = User{}
	stmt := "SELECT * FROM USERS WHERE PASSWORD=$1;"
	row := db.QueryRow(stmt, user.Password)
	row.Scan(&user.id, &user.Name, &u.Email, &u.Password, &user.hash)
	if user.Password == u.Password && u.Email == user.Email {
		fmt.Println("user authentificated !")
		return true
	}
	return false
}
