package main

import (
	"log"
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"database/sql"
	_ "github.com/lib/pq"
)

type Post struct {
	Id 		int 		`json:"id"`
	Title string 	`json:"title"`
	Body 	string 	`json:"body`
}

const (
	hostname = "localhost"
	port = 5455
	user = "blogo"
	password = "blogo123"
	dbname = "blogo"
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
	
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var posts []Post
	for rows.Next() {
		var id int
		var title string
		var body string

		err = rows.Scan(&id, &title, &body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		posts = append(posts, Post{Id: id, Title: title, Body: body})
	}
	json.NewEncoder(w).Encode(posts)
}

func NewPost(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
 	}
 	r.Body.Close()
	post := Post{}
	json.Unmarshal([]byte(body), &post)
	log.Println(post)

	var lastInsertID int
	err = db.QueryRow("INSERT INTO posts(id, title, body) VALUES($1, $2, $3) returning id;", post.Id, post.Title, post.Body).Scan(&lastInsertID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}	
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", Home).Methods("GET")
	router.HandleFunc("/post", NewPost).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
	log.Println("Server running...")
}


