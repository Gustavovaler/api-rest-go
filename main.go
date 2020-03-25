package main

import (
	"fmt"
	mux "github.com/gorilla/mux"
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
		notes = append(notes, v)
	}

	w.Header().Set("Content-type", "application/json")
	j, err := json.Marshal(notes)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOk)
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
	j, err := json.Marshal(notes)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(j)

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
	fmt.Println("Servidor escuchando el puerto 8080")

	server.ListenAndServe()
}
