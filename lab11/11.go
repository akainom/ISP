package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "lab11/docs"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title lab11
// @description GO11 swagger
// @BasePath /

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

// createNewCelebrity godoc
// @Summary      create celebrity
// @Description  adds new celebrity
// @Tags         celebrities
// @Accept       json
// @Produce      json
// @Param        celebrity body Celebrity true "data"
// @Success      201 {object} Celebrity
// @Failure      400 {string} string "Bad Request"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /celebrities [post]
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
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(c)
	log.Println("[LOG] created new id " + strconv.Itoa(c.Id))
}

// getAllCels godoc
// @Summary      get celebrities
// @Description  returns array of celebrities
// @Tags         celebrities
// @Produce      json
// @Success      200 {object} Celebrity
// @Failure      500 {string} string "Internal Server Error"
// @Router       /celebrities [get]
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

// getCel godoc
// @Summary      get celebrity
// @Description  returns celebrity
// @Tags         celebrities
// @Produce      json
// @Success      200 {object} Celebrity
// @Failure      500 {string} string "Internal Server Error"
// @Router       /celebritiy/{id} [get]
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

// updateCel godoc
// @Summary      updates celebrity
// @Description  returns updated celebrity
// @Tags         celebrities
// @Produce      json
// @Param        celebrity body Celebrity true "data"
// @Success      200 {object} Celebrity
// @Failure      500 {string} string "Internal Server Error"
// @Router       /celebritiy/{id} [put]
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

// deleteCel godoc
// @Summary      delete celebrity
// @Description  returns deleted celebrity
// @Tags         celebrities
// @Produce      json
// @Success      204 {object} Celebrity
// @Failure      500 {string} string "Internal Server Error"
// @Router       /celebritiy/{id} [delete]
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

	file, err := os.OpenFile("11.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/swagger/doc.json"),
	))

	fmt.Println("server started on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
