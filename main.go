package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var tmp *template.Template

type User struct {
	id       int
	Email    string
	Name     string
	Password string
	hash     string
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
	row.Scan(&user.id, &user.Name, &user.Password, &user.hash)
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

func main() {
	router := mux.NewRouter()

	files := http.FileServer(http.Dir("assets"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", files))

	router.HandleFunc("/", homePage)
	router.HandleFunc("/all", all)
	router.HandleFunc("/new-poll", newPoll)
	router.HandleFunc("/new-user/", newUser)

	router.HandleFunc("/current/{id}", current)
	router.HandleFunc("/current/{id}/api/", sendJSON)
	router.HandleFunc("/create-pool", createPoll)

	router.HandleFunc("/submit-vote", submitVote)

	router.HandleFunc("/login", login)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/signupPage", signupPage)
	router.HandleFunc("/loginPage", loginPage)

	server := http.Server{
		Addr:    getPort(),
		Handler: router,
	}
	server.ListenAndServe()
}

func getPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		log.Println("Running on localhost:4200...")
		return ":4200"
	}
	log.Println("Running on port: " + port)
	return ":" + port
}
