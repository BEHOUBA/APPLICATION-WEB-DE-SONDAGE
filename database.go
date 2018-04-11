package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Poll struct {
	ID      int
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

func getUserDataFromDB(cookie string, r *http.Request) (User, error) {
	cook, err := r.Cookie(cookie)
	if err != nil {
		return User{}, err
	}
	var u User
	stmt := "SELECT * FROM USERS WHERE PASSWORD_HASH=$1;"
	row := db.QueryRow(stmt, cook.Value)
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

func getAllPollTitle(limit int) (idAndTitle map[int]string, err error) {
	idAndTitle = map[int]string{}
	stmt := "SELECT POLL_ID, POLL_NAME FROM POLLS LIMIT $1;"
	rows, err := db.Query(stmt, limit)
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var id int
		var title string
		rows.Scan(&id, &title)
		idAndTitle[id] = title
	}
	return
}

func (p *Poll) getPollOptions() (err error) {
	query := "SELECT OPTION_NAME FROM POLL_OPTIONS WHERE POLL=$1;"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println(err)
		return
	}
	rows, err := stmt.Query(p.ID)
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var option string
		rows.Scan(&option)
		p.Options = append(p.Options, option)
	}
	return
}

func (p *Poll) submitVote(opt string) (err error) {
	if p.asAlreadyVoted(currentUser) {
		return errors.New("USER AS ALREADY PARTICIPATED!")
	}
	var optID int
	row := db.QueryRow("SELECT OPTION_ID FROM POLL_OPTIONS WHERE OPTION_NAME=$1 AND POLL=$2", opt, p.ID)
	if err != nil {
		log.Println(err)
		return
	}
	row.Scan(&optID)

	query := "INSERT INTO VOTES (OPTIONS, USER_ID, POLL_ID) VALUES ($1, $2, $3);"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = stmt.Exec(optID, currentUser.id, p.ID)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (p *Poll) asAlreadyVoted(user User) bool {
	fmt.Println(user.id, p.ID)
	res, err := db.Exec("SELECT VOTE_ID FROM VOTES WHERE USER_ID=$1 AND POLL_ID=$2;", user.id, p.ID)
	if err != nil {
		log.Println(err)
		return true
	}
	VID, _ := res.RowsAffected()
	fmt.Println(VID)
	if VID != 0 {
		return true
	}
	return false

}
