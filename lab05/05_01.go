package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Celebrity struct {
	Id           int    `json:"id"`
	Fullname     string `json:"fullName"`
	Nationality  string `json:"nationality"`
	ReqPhotoPath string `json:"reqPhotoPath"`
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./celebrities.db")
	if err != nil {
		log.Fatal(err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS celebrities (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fullName TEXT NOT NULL,
		nationality TEXT NOT NULL,
		reqPhotoPath TEXT
	);`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("[FATAL] не удалось создать таблицу:", err)
	}
}

func createNewCelebrity(w http.ResponseWriter, r *http.Request) {
	var c Celebrity
	json.NewDecoder(r.Body).Decode(&c)

	res, err := db.Exec("INSERT INTO celebrities (fullName, nationality, reqPhotoPath) VALUES (?, ?, ?)",
		c.Fullname, c.Nationality, c.ReqPhotoPath)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	id, _ := res.LastInsertId()
	c.Id = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
	log.Println("[LOG] created new id " + strconv.Itoa(c.Id))
}

func getAllCels(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, fullName, nationality, reqPhotoPath FROM celebrities")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var cels []Celebrity
	for rows.Next() {
		var c Celebrity
		rows.Scan(&c.Id, &c.Fullname, &c.Nationality, &c.ReqPhotoPath)
		cels = append(cels, c)
	}
	json.NewEncoder(w).Encode(cels)
}

func getCel(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var c Celebrity

	err := db.QueryRow("SELECT id, fullName, nationality, reqPhotoPath FROM celebrities WHERE id = ?", id).
		Scan(&c.Id, &c.Fullname, &c.Nationality, &c.ReqPhotoPath)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func updateCel(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var c Celebrity
	json.NewDecoder(r.Body).Decode(&c)

	_, err := db.Exec("UPDATE celebrities SET fullName=?, nationality=?, reqPhotoPath=? WHERE id=?",
		c.Fullname, c.Nationality, c.ReqPhotoPath, id)
	if err != nil {
		w.WriteHeader(500)
		log.Println("[LOG] unable to update id " + strconv.Itoa(id))
		return
	}
	w.WriteHeader(200)
	log.Println("[LOG] updated id " + strconv.Itoa(c.Id))

}

func deleteCel(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	_, err := db.Exec("DELETE FROM celebrities WHERE id = ?", id)
	if err != nil {
		w.WriteHeader(500)
		log.Println("[LOG] unable to delete id " + strconv.Itoa(id))
		return
	}
	w.WriteHeader(204)
	log.Println("[LOG] deleted id " + strconv.Itoa(id))
}

func main() {
	initDB()
	defer db.Close()

	file, err := os.OpenFile("05_01.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Errorf("\n[FATAL] unable to open log file, aborting")
	}
	log.SetOutput(file)

	r := mux.NewRouter()
	r.HandleFunc("/celebrities", createNewCelebrity).Methods("POST")
	r.HandleFunc("/celebrities", getAllCels).Methods("GET")
	r.HandleFunc("/celebrity/{id}", getCel).Methods("GET")
	r.HandleFunc("/celebrity/{id}", updateCel).Methods("PUT")
	r.HandleFunc("/celebrity/{id}", deleteCel).Methods("DELETE")

	fmt.Println("server started on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
