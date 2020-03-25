package main

import (
	"encoding/json"
	mux "github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Note struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

var noteStore = make(map[string]Note)
var id int

func GETnoteHandler(w http.ResponseWriter, r *http.Request) {
	var notes []Note
	for _, value := range noteStore {
		notes = append(notes, value)
	}

	w.Header().Set("Content-type", "application/json")
	j, err := json.Marshal(notes)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)
	w.Write(j)
}

func POSTnoteHandler(w http.ResponseWriter, r *http.Request) {
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)

	if err != nil {
		panic(err)
	}
	note.CreatedAt = time.Now()
	id++
	k := strconv.Itoa(id)
	noteStore[k] = note

	w.Header().Set("Content-type", "application/json")
	j, err := json.Marshal(note)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(j)

}

func PUTnoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]
	var noteUpdate Note
	err := json.NewDecoder(r.Body).Decode(&noteUpdate)
	if err != nil {
		panic(err)
	}

	if note, ok := noteStore[k]; ok {
		noteUpdate.CreatedAt = note.CreatedAt
		delete(noteStore, k)
		noteStore[k] = noteUpdate
	} else {
		log.Printf("No se encontro ese id %s", k)
	}
	w.WriteHeader(http.StatusNoContent)

}

func DELETEnoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]

	if _, ok := noteStore[k]; ok {

		delete(noteStore, k)

	} else {
		log.Printf("No se encontro ese id %s", k)
	}
	w.WriteHeader(http.StatusNoContent)

}

func main() {

	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/notes", GETnoteHandler).Methods("GET")
	r.HandleFunc("/api/notes/{id}", PUTnoteHandler).Methods("PUT")
	r.HandleFunc("/api/notes", POSTnoteHandler).Methods("POST")
	r.HandleFunc("/api/notes/{id}", DELETEnoteHandler).Methods("DELETE")

	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//fmt.Println("Servidor escuchando el puerto 8080")

	server.ListenAndServe()
}
