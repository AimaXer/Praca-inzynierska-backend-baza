package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
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

type bodyMessageUpdateUser struct {
	UserToken string `json:"userToken"`
	Username  string `json:"username"`
	Preferences  string `json:"preferences"`
}

func getChatGrups(preferences string) []string {
	var groups []string
	kwiaty := strings.Contains(preferences, "kwiaty")
	kolory := strings.Contains(preferences, "kolory")
	rozrywka := strings.Contains(preferences, "rozrywka")
	architektura := strings.Contains(preferences, "architektura")
	jednoslady := strings.Contains(preferences, "jednoslady")
	wyscigi := strings.Contains(preferences, "wyscigi")
	if kwiaty || kolory {
		groups = append(groups, "grupa-kwiatowa")
	}
	if rozrywka || architektura {
		groups = append(groups, "grupa-miastowa")
	}
	if jednoslady || wyscigi {
		groups = append(groups, "grupa-motocyklowa")
	}
	if kwiaty || architektura {
		groups = append(groups, "grupa-ogrodowa")
	}
	if jednoslady || rozrywka {
		groups = append(groups, "grupa-rowerowa")
	}
	return groups
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/createUser", createUser).Methods("POST")
	myRouter.HandleFunc("/getUserInfo", getUserInfo).Methods("POST")
	myRouter.HandleFunc("/updateUserInfo", updateUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func getUserInfo(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host="+DB_IP+" port=5432 user=postgres dbname=postgres sslmode=disable password=postgres")
	jsons := simplejson.New()

	defer db.Close()
	if err != nil {
		panic(err)
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
	db, err := sql.Open("postgres", "host="+DB_IP+" port=5432 user=postgres dbname=postgres sslmode=disable password=postgres")

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

	_, err = db.Exec(fmt.Sprintf("INSERT INTO users.users(username, token, chatflow, mainchats) VALUES ('%s', '%s', '', '')", username, userToken))
	if err != nil {
		fmt.Printf("\nexec\n")
		panic(err.Error())
	}

	fmt.Fprintf(w, "New user was added")
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host="+DB_IP+" port=5432 user=postgres dbname=postgres sslmode=disable password=postgres")

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

	var msg bodyMessageUpdateUser

	json.Unmarshal([]byte(body), &msg)
	
	userToken := msg.UserToken
	chats := strings.Join(getChatGrups(msg.Preferences), ",")
	
	_, err = db.Exec(fmt.Sprintf("UPDATE users.users SET chatflow='%s', mainchats='%s' WHERE token='%s';", chats, chats,userToken))
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
