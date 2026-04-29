package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Celebrity struct {
	Id           int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Fullname     string `json:"fullName"`
	Nationality  string `json:"nationality"`
	ReqPhotoPath string `json:"reqPhotoPath"`
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("./celebrities.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("[FATAL] failed to connect database:", err)
	}

	db.AutoMigrate(&Celebrity{})
}

func createNewCelebrity(w http.ResponseWriter, r *http.Request) {
	var c Celebrity
	json.NewDecoder(r.Body).Decode(&c)

	result := db.Create(&c)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
	log.Println("[LOG] created new id " + strconv.Itoa(c.Id))
}

func getAllCels(w http.ResponseWriter, r *http.Request) {
	var cels []Celebrity
	db.Find(&cels)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cels)
}

func getCel(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var c Celebrity

	if err := db.First(&c, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
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

	var existing Celebrity
	if err := db.First(&existing, id).Error; err != nil {
		w.WriteHeader(404)
		return
	}

	db.Model(&existing).Updates(c)

	w.WriteHeader(200)
	log.Println("[LOG] updated id " + strconv.Itoa(id))
}

func deleteCel(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if err := db.Delete(&Celebrity{}, id).Error; err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
	log.Println("[LOG] deleted id " + strconv.Itoa(id))
}

func main() {
	initDB()

	file, _ := os.OpenFile("06_01.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
