package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var currentUser User
var homeData PageData
var err error
var currentPoll Poll
var data PageData
var cData chartData

type chartData struct {
	Title string          `json:"title"`
	Data  [][]interface{} `json:"data"`
}

type Vote struct {
	Name    string
	Poll_ID int
}

type PageData struct {
	CurrentUser         User
	AllPollsIdAndTitles map[int]string
	CurrentPoll         Poll
	Vote                Vote
	Error               error
}

func homePage(w http.ResponseWriter, r *http.Request) {
	homeData.AllPollsIdAndTitles, err = getAllPollTitle(12, 0)
	if err != nil {
		log.Println(err)
	}

	currentUser, err = getUserDataFromDB("voting_app", r)
	if err != nil {
		log.Println(err)
	}
	homeData.CurrentUser = currentUser
	generateHTML(w, r, homeData, "home", "home.content", "footer")
}

func all(w http.ResponseWriter, r *http.Request) {
	homeData.AllPollsIdAndTitles, err = getAllPollTitle(100, 0)
	generateHTML(w, r, homeData, "home", "all.content", "footer")
}

func newPoll(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/loginPage", 302)
		return
	}
	generateHTML(w, r, homeData, "home", "new.poll", "footer")
}

func signupPage(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, r, "", "home", "signup", "footer")
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", 302)
		return
	}
	generateHTML(w, r, "", "home", "login", "footer")
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("voting_app")
	if err != nil {
		http.Redirect(w, r, "/", 302)
		log.Println("cookie don't exist!")
		return
	}
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 302)
}

func generateHTML(writer http.ResponseWriter, r *http.Request, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	if alreadyLoggedIn(r) {
		files = append(files, "templates/header-verified.html")
	} else {
		files = append(files, "templates/header.html")
	}
	templ, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err)
	}
	templ.ExecuteTemplate(writer, "layout", data)
}

func alreadyLoggedIn(r *http.Request) bool {
	//var cook *http.Cookie
	_, err := r.Cookie("voting_app")
	if err != nil {
		return false
	}
	return true
}

func login(w http.ResponseWriter, r *http.Request) {
	user := getUserData(r)
	if user.authentification() {
		user.createSession(w)
		http.Redirect(w, r, "/", 302)
		return
	}
	log.Println("invalid login!")
	fmt.Fprintln(w, "Email or Password is wrong... or maybe you are not registred yet...")
	http.Redirect(w, r, "/loginPage", 302)
}
