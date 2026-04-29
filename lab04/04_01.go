package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Celebrity struct {
	Id           int    `json:"id"`
	Fullname     string `json:"fullName"`
	Nationality  string `json:"nationality"`
	ReqPhotoPath string `json:"reqPhotoPath"`
}

type CelebrityDTO struct {
	Fullname     string `json:"fullName"`
	Nationality  string `json:"nationality"`
	ReqPhotoPath string `json:"reqPhotoPath"`
}

func newCelebrity(fullname string, nationality string, reqphotopath string) *Celebrity {
	newCel := new(Celebrity)
	lastId += 1

	newCel.Id = lastId
	newCel.Fullname = fullname
	newCel.Nationality = nationality
	newCel.ReqPhotoPath = reqphotopath

	return newCel
}

var Celebrities = make(map[int]Celebrity)
var lastId int = 0

func saveJSON() {
	var JSON, err = json.MarshalIndent(Celebrities, "", " ")
	if err != nil {
		log.Fatalln("[FATAL] JSON arr is corrupted")
	}

	err = os.WriteFile("celebrities.json", JSON, 0644)
	if err != nil {
		log.Fatalln("[FATAL] unable to write to celebrities.json")
	} else {
		log.Println("[LOG] saved to JSON at " + time.Now().String())
	}
}

func loadJSON() {
	var JSON, err = os.ReadFile("celebrities.json")
	if err != nil {
		log.Fatalln("[FATAL] unable to read celebrities.json")
	}

	err = json.Unmarshal(JSON, &Celebrities)
	if err != nil {
		log.Fatalln("[FATAL] unable to save celebrities from JSON")
	} else {
		log.Println("[LOG] saved celebritites to JSON")
	}
}

func createCelebrity(fullname string, nationality string, reqphotopath string) Celebrity {
	newCel := newCelebrity(fullname, nationality, reqphotopath)
	Celebrities[newCel.Id] = *newCel

	log.Println("[LOG] created id " + strconv.Itoa(newCel.Id))
	return *newCel
}

func getCelebrity(id int) Celebrity {
	return Celebrities[id]
}

func getCelebrities() map[int]Celebrity {
	return Celebrities
}

func updateCelebrity(id int, cel Celebrity) {
	Celebrities[id] = cel
	log.Println("[LOG] updated id " + strconv.Itoa(id))
	saveJSON()
}

func deleteCelebrity(id int) bool {
	var prevLength int = len(Celebrities)
	delete(Celebrities, id)
	var afterLength int = len(Celebrities)

	saveJSON()
	if prevLength > afterLength {
		log.Println("[LOG] deleted id " + strconv.Itoa(id))
		return true
	}

	log.Println("[LOG] unable to delete id " + strconv.Itoa(id))
	return false
}

func createNewCel(w http.ResponseWriter, r *http.Request) {
	var dto CelebrityDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	newCel := createCelebrity(dto.Fullname, dto.Nationality, dto.ReqPhotoPath)

	saveJSON()
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(newCel)
}

func getAllCels(w http.ResponseWriter, r *http.Request) {
	cels := getCelebrities()

	JSON, _ := json.MarshalIndent(cels, "", " ")
	w.WriteHeader(200)
	w.Write(JSON)
}

func getCel(w http.ResponseWriter, r *http.Request) {
	id := getId(r)
	cel := getCelebrity(id)

	JSON, _ := json.MarshalIndent(cel, "", " ")
	w.WriteHeader(200)
	w.Write(JSON)
}

func updateCel(w http.ResponseWriter, r *http.Request) {
	id := getId(r)
	var cel Celebrity
	json.NewDecoder(r.Body).Decode(&cel)
	saveJSON()

	updateCelebrity(id, cel)
	JSON, _ := json.MarshalIndent(cel, "", " ")
	w.WriteHeader(200)
	w.Write(JSON)
}

func deleteCel(w http.ResponseWriter, r *http.Request) {
	id := getId(r)
	var cel Celebrity = Celebrities[id]

	result := deleteCelebrity(id)
	saveJSON()

	if result {
		JSON, _ := json.MarshalIndent(cel, "", " ")
		w.Write(JSON)
	} else {
		w.WriteHeader(400)
	}
}

func getId(r *http.Request) int {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	return id
}

func main() {
	file, err := os.OpenFile("04_01.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Errorf("\n[FATAL] unable to open log file, aborting")
	}
	log.SetOutput(file)

	loadJSON()
	for id := range Celebrities {
		if id < lastId {
			lastId = id
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/celebrities", createNewCel).Methods("POST")
	r.HandleFunc("/celebrities", getAllCels).Methods("GET")
	r.HandleFunc("/celebrity/{id}", getCel).Methods("GET")
	r.HandleFunc("/celebrity/{id}", updateCel).Methods("PUT")
	r.HandleFunc("/celebrity/{id}", deleteCel).Methods("DELETE")

	fmt.Println("server starting at localhost:3000 ...")
	log.Fatalln(http.ListenAndServe(":8080", r))
}
