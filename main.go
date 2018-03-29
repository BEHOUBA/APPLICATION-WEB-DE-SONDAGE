package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var tmp *template.Template

func main() {
	router := http.NewServeMux()

	files := http.FileServer(http.Dir("assets"))
	router.Handle("/assets/", http.StripPrefix("/assets/", files))

	router.HandleFunc("/", homePage)
	router.HandleFunc("/all", all)
	router.HandleFunc("/new-poll", newPoll)
	router.HandleFunc("/signup", signup)
	router.HandleFunc("/login", login)
	router.HandleFunc("/new-user", newUser)

	server := http.Server{
		Addr:    getPort(),
		Handler: router,
	}
	server.ListenAndServe()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "home", "home.content", "footer")
}

func all(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "home", "all.content", "footer")
}

func newPoll(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "home", "new.poll", "footer")
}

func signup(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "home", "signup", "footer")
}

func login(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "home", "login", "footer")
}

func newUser(w http.ResponseWriter, r *http.Request) {
	var name, email, password string
	r.ParseForm()
	name = r.FormValue("name")
	email = r.FormValue("email")
	password = r.FormValue("password")
	fmt.Println(name, email, password)
}

func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	if false {
		files = append(files, "templates/header-verified.html")
	} else {
		files = append(files, "templates/header.html")
	}
	templ := template.Must(template.ParseFiles(files...))
	templ.ExecuteTemplate(writer, "layout", "MANASSE")
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
