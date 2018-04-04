package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

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
	cookie := &http.Cookie{
		Name:     "voting_app",
		Value:    user.hash,
		HttpOnly: true,
		MaxAge:   600,
	}
	http.SetCookie(w, cookie)
	return
}

type Session struct {
}

var tmp *template.Template

func main() {
	router := http.NewServeMux()

	files := http.FileServer(http.Dir("assets"))
	router.Handle("/assets/", http.StripPrefix("/assets/", files))

	router.HandleFunc("/", homePage)
	router.HandleFunc("/all", all)
	router.HandleFunc("/new-poll", newPoll)
	router.HandleFunc("/new-user", newUser)

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
