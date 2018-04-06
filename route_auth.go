package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var currentUser User
var homeData HomePageData
var err error

type HomePageData struct {
	CurrentUser    User
	AllPollsTitles []string
}

func homePage(w http.ResponseWriter, r *http.Request) {
	homeData.AllPollsTitles, err = getAllPollTitle()
	if err != nil {
		log.Println(err)
	}
	cookie, err := r.Cookie("voting_app")
	if err != nil {
		log.Println(err)
	} else {
		currentUser, err = getUserDataFromDB(cookie.Value)
		homeData.CurrentUser = currentUser
	}
	generateHTML(w, r, homeData, "home", "home.content", "footer")
}

func all(w http.ResponseWriter, r *http.Request) {
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

func newUser(w http.ResponseWriter, r *http.Request) {
	u := getUserData(r)
	err := u.createUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.createSession(w)
	http.Redirect(w, r, "/", 302)
	fmt.Println(u)
}

func createPoll(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/loginPage", 302)
		return
	}
	r.ParseForm()
	poll := Poll{r.FormValue("poll-name"), strings.Split(r.FormValue("poll-options"), "×")[1:]}
	err := poll.addPollToDB()
	if err != nil {
		log.Println("failed to add create poll", err)
		fmt.Fprintln(w, "Error occured please contact admistrator: behouba@gmail.com.\n Error:", err)
		return
	}
	generateHTML(w, r, "", "home", "current-poll", "footer")
}

func getUserData(r *http.Request) (u User) {
	r.ParseForm()
	u.Name = r.FormValue("name")
	u.Email = r.FormValue("email")
	u.Password = r.FormValue("password")
	return u
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
		return
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
