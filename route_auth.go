package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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

type PageData struct {
	CurrentUser         User
	AllPollsIdAndTitles map[int]string
	CurrentPoll         Poll
}

func homePage(w http.ResponseWriter, r *http.Request) {
	homeData.AllPollsIdAndTitles, err = getAllPollTitle(10, 0)
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
	var poll Poll
	var options []string
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/loginPage", 302)
		return
	}
	r.ParseForm()

	poll.Title = r.FormValue("poll-name")
	options = strings.Split(r.FormValue("poll-options"), "×")[1:]

	for _, val := range options {
		poll.Options = append(poll.Options, Option{Name: val})
	}

	err := poll.addPollToDB()
	if err != nil {
		log.Println("failed to add create poll", err)
		fmt.Fprintln(w, "Error occured please contact admistrator: behouba@gmail.com.\n Error:", err)
		return
	}
	generateHTML(w, r, "", "home", "current-poll", "footer")
}

func setCurrentData(w http.ResponseWriter, r *http.Request) {
	currentUser, _ = getUserDataFromDB("voting_app", r)
	currentPoll.ID, err = strconv.Atoi(mux.Vars(r)["id"])
	data.CurrentUser = currentUser
	data.CurrentPoll.ID = currentPoll.ID
	data.CurrentPoll.setTitle()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data.CurrentPoll.getAndSetTotalVotes()
	data.CurrentPoll.getPollOptions()

	for i, opt := range data.CurrentPoll.Options {
		data.CurrentPoll.Options[i].Votes, _ = opt.getAndSetTotalVotes()
		data.CurrentPoll.Options[i].Percentage = float32(data.CurrentPoll.Options[i].Votes * 100 / data.CurrentPoll.Votes)
	}
}

func current(w http.ResponseWriter, r *http.Request) {
	setCurrentData(w, r)
	generateHTML(w, r, data, "home", "current-poll", "footer")
	return
}

func sendJSON(w http.ResponseWriter, r *http.Request) {
	setCurrentData(w, r)
	cData.Data = cData.Data[:0]
	cData.Title = data.CurrentPoll.Title
	for _, opt := range data.CurrentPoll.Options {
		fmt.Println(opt)
		val := []interface{}{opt.Name, opt.Votes}
		cData.Data = append(cData.Data, val)
	}

	jsonByte, err := json.Marshal(cData)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprint(w, string(jsonByte))
}

func submitVote(w http.ResponseWriter, r *http.Request) {
	currentUser, err = getUserDataFromDB("voting_app", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	r.ParseForm()
	vote := r.FormValue("option")
	err = currentPoll.submitVote(vote)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintln(w, vote)
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
