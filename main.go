package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var tmp *template.Template

func main() {
	router := mux.NewRouter()

	files := http.FileServer(http.Dir("assets"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", files))

	router.HandleFunc("/", homePage)
	router.HandleFunc("/all", all)
	router.HandleFunc("/new-poll", newPoll)
	router.HandleFunc("/new-user", newUser)

	router.HandleFunc("/current/{id}", current)
	router.HandleFunc("/current/{id}/api/", sendJSON)
	router.HandleFunc("/create-pool", createPoll)

	router.HandleFunc("/current/{id}/submit_vote", submitVote)

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
