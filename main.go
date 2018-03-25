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

	server := http.Server{
		Addr:    getPort(),
		Handler: router,
	}
	server.ListenAndServe()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "home", "home.content", "footer", "header")
}

func all(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "home", "all.content", "footer", "header")
}

func newPoll(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "", "home", "new.poll", "header", "footer")
}

func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	templ := template.Must(template.ParseFiles(files...))
	templ.ExecuteTemplate(writer, "layout", data)
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
