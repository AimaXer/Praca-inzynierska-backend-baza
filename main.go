package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	db *sql.DB
)

var DB_IP = "localhost"

type bodyMessageUserInfo struct {
	UserToken string `json:"userToken"`
}

type bodyMessageAddUser struct {
	UserToken string `json:"userToken"`
	Username  string `json:"username"`
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/createUser", createUser).Methods("POST")
	myRouter.HandleFunc("/getUserInfo", getUserInfo).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
func getUserInfo(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host="+DB_IP+" port=5432 user=postgres dbname=InzApp sslmode=disable password=Maciek0808")
	jsons := simplejson.New()

	defer db.Close()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("DB call get User info\n")
	}
	rows, _ := db.Query(fmt.Sprintf("SELECT * FROM users.users"))
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Printf("ioutil")
		panic(err.Error())
	}
	var msg bodyMessageUserInfo

	json.Unmarshal([]byte(body), &msg)

	userToken := msg.UserToken

	for rows.Next() {
		var (
			username  string
			token     string
			chatflow  string
			mainchats string
		)
		if err := rows.Scan(&username, &token, &chatflow, &mainchats); err != nil {
			log.Fatal(err)
		}
		if userToken == token {
			jsons.Set("ChatFlow", chatflow)
			jsons.Set("MainChats", mainchats)
		}
	}
	payload, err := jsons.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host="+DB_IP+" port=5432 user=postgres dbname=InzApp sslmode=disable password=Maciek0808")

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("DB connected")
	}

	defer db.Close()

	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("ioutil")
		panic(err.Error())
	}

	var msg bodyMessageAddUser

	json.Unmarshal([]byte(body), &msg)

	userToken := msg.UserToken
	username := msg.Username

	_, err = db.Exec(fmt.Sprintf("INSERT INTO users.users(username, token, chatflow, mainchats) VALUES ('%s', '%s', 'grupa-rowerowa,grupa-motocyklowa', 'grupa-rowerowa,grupa-motocyklowa,grupa-kwiatowa')", username, userToken))
	if err != nil {
		fmt.Printf("\nexec\n")
		panic(err.Error())
	}

	fmt.Fprintf(w, "New user was added")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage")
}

func main() {
	handleRequest()
}
