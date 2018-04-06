package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Poll struct {
	Title   string
	Options []string
}

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
	stmt := "SELECT * FROM USERS WHERE PASSWORD=$1 AND EMAIL=$2;"
	row := db.QueryRow(stmt, user.Password, user.Email)
	row.Scan(&user.id, &user.Name, &user.Email, &user.Password, &user.hash)
	if user.Email != "" {
		fmt.Println("user authentificated !")
		return true
	}
	return false
}

func getUserDataFromDB(cookie string) (User, error) {
	var u User
	stmt := "SELECT * FROM USERS WHERE PASSWORD_HASH=$1;"
	row := db.QueryRow(stmt, cookie)
	row.Scan(&u.id, &u.Name, &u.Email, &u.Password, &u.hash)
	if u.Name != "" {
		return u, nil
	}
	return u, errors.New("user not found in database...")
}

func (p *Poll) addPollToDB() (err error) {
	var pollID int
	stmt := "INSERT INTO POLLS (POLL_NAME, OWNER_ID) VALUES ($1, $2) RETURNING poll_id;"
	err = db.QueryRow(stmt, p.Title, currentUser.id).Scan(&pollID)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(pollID)
	for _, val := range p.Options {
		stmt2 := "INSERT INTO POLL_OPTIONS (OPTION_NAME, POLL) VALUES ($1, $2);"
		_, err := db.Exec(stmt2, val, pollID)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	fmt.Println("POLL ADD TO DATABASE SUCCESFULLY!!!")
	return
}

func getAllPollTitle() (titles []string, err error) {
	titles = []string{}
	stmt := "SELECT POLL_NAME FROM POLLS LIMIT 10;"
	rows, err := db.Query(stmt)
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var title string
		rows.Scan(&title)
		titles = append(titles, title)
	}
	fmt.Println(titles)
	return
}
