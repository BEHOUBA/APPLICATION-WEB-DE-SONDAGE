package main

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Email    string
	Name     string
	Password string
	hash     string
	CanVote  bool
	OwnPolls map[int]string
}

func (user *User) hashPassword() string {
	hashByte, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		panic("Failled to generate hash value from password!")
	}
	user.hash = string(hashByte)
	return user.hash
}

func (user *User) createSession(w http.ResponseWriter) {
	row := db.QueryRow("UPDATE USERS SET PASSWORD_HASH=$1 WHERE EMAIL=$2", user.hashPassword(), user.Email)
	row.Scan(&user.ID, &user.Name, &user.Password, &user.hash)
	cookie := &http.Cookie{
		Name:     "voting_app",
		Value:    user.hash,
		HttpOnly: true,
		MaxAge:   200000,
	}
	fmt.Println("HASH UPDATED...", user.hash)
	http.SetCookie(w, cookie)
	return
}

func newUser(w http.ResponseWriter, r *http.Request) {
	u := getUserData(r)
	err := u.createUser()
	if err != nil {
		data.Error = err
		generateHTML(w, r, data, "home", "error", "footer")
		return
	}
	u.createSession(w)
	http.Redirect(w, r, "/", 302)
	fmt.Println(u)
}

func getUserData(r *http.Request) (u User) {
	r.ParseForm()
	u.Name = r.FormValue("name")
	u.Email = r.FormValue("email")
	u.Password = r.FormValue("password")
	return u
}
