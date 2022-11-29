package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Entry struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt int64  `json:"created_at"`
	Views     int    `json:"views"`
	Tags      string `json:"tags"`
}

const (
	hostname = "localhost"
	port     = 5432
	user     = "blogo"
	password = "blogo123"
	dbname   = "blogo"
)

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", hostname, port, user, password, dbname)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}

	return db
}

func Home(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	rows, err := db.Query("SELECT * FROM entries")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var entries []Entry
	for rows.Next() {
		var id int
		var title string
		var body string
		var createdAt time.Time
		var views int
		var tags string

		err = rows.Scan(&id, &createdAt, &tags, &views, &title, &body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		entries = append(entries, Entry{Id: id, Title: title, Body: body, CreatedAt: createdAt.Unix(), Views: views, Tags: tags})
	}
	json.NewEncoder(w).Encode(entries)
}

func NewEntry(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	r.Body.Close()
	entry := Entry{}
	json.Unmarshal([]byte(body), &entry)
	log.Println(entry)

	var lastInsertID int
	err = db.QueryRow("INSERT INTO entries(title, body, created_at, views, tags) VALUES($1, $2, $3, $4, $5) returning id;", entry.Title, entry.Body, time.Now(), entry.Views, entry.Tags).Scan(&lastInsertID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", Home).Methods("GET")
	router.HandleFunc("/entry", NewEntry).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
