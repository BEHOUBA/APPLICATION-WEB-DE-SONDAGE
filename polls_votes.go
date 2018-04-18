package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func submitVote(w http.ResponseWriter, r *http.Request) {
	setCurrentData(w, r)
	data.Vote.Poll_ID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
	}
	currentUser, err = getUserDataFromDB("voting_app", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	r.ParseForm()
	data.Vote.Name = r.FormValue(data.CurrentPoll.Title)
	err = currentPoll.submitVote(data.Vote.Name)
	if err != nil {
		log.Println(err)
		data.Error = err
		generateHTML(w, r, data, "home", "error", "footer")
		return
	}
	generateHTML(w, r, data, "home", "success.vote", "footer")
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
	fmt.Println(options, "this is options")
	for _, val := range options {
		poll.Options = append(poll.Options, Option{Name: val})
	}

	id, err := poll.addPollToDB()
	if err != nil {
		log.Println("failed to add create poll", err)
		fmt.Fprintln(w, "Error occured please contact admistrator: behouba@gmail.com.\n Error:", err)
		return
	}
	url := "/current/" + strconv.Itoa(id)
	log.Println(url)
	http.Redirect(w, r, url, 302)
}

func userPolls(w http.ResponseWriter, r *http.Request) {
	setCurrentData(w, r)
	data.CurrentUser.OwnPolls, err = currentUser.getOwnPoll()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(data)
	generateHTML(w, r, data, "home", "user.polls", "footer")
}

func setCurrentData(w http.ResponseWriter, r *http.Request) {
	currentUser, _ = getUserDataFromDB("voting_app", r)
	currentPoll.ID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
	}
	data.CurrentUser = currentUser
	data.CurrentPoll.ID = currentPoll.ID
	data.CurrentPoll.setTitle()
	data.CurrentPoll.getAndSetTotalVotes()
	data.CurrentPoll.getPollOptions()

	if data.CurrentUser.ID == 0 {
		data.CurrentUser.CanVote = false
	} else {
		data.CurrentUser.CanVote = data.CurrentPoll.canVote(data.CurrentUser)
	}

	for i, opt := range data.CurrentPoll.Options {
		data.CurrentPoll.Options[i].Votes, _ = opt.getAndSetTotalVotes()
		if data.CurrentPoll.Votes == 0 {
			data.CurrentPoll.Options[i].Percentage = float32(data.CurrentPoll.Options[i].Votes * 100 / 1)
		} else {
			data.CurrentPoll.Options[i].Percentage = float32(data.CurrentPoll.Options[i].Votes * 100 / data.CurrentPoll.Votes)
		}

	}
}

func deletePoll(w http.ResponseWriter, r *http.Request) {
	data.CurrentPoll.ID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		return
	}
	err = data.CurrentPoll.deletePoll()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	http.Redirect(w, r, "/user_polls", 302)
}

func current(w http.ResponseWriter, r *http.Request) {
	setCurrentData(w, r)
	fmt.Println(data)
	generateHTML(w, r, data, "home", "current-poll", "footer")
}

func sendJSON(w http.ResponseWriter, r *http.Request) {
	setCurrentData(w, r)
	cData.Data = cData.Data[:0]
	cData.Title = data.CurrentPoll.Title
	for _, opt := range data.CurrentPoll.Options {
		val := []interface{}{opt.Name, opt.Percentage}
		cData.Data = append(cData.Data, val)
	}

	jsonByte, err := json.Marshal(cData)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprint(w, string(jsonByte))
}
