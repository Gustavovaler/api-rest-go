package main

import (
	"encoding/json"
	"fmt"
	mux "github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

// definimos una estructura de tipo Note
type Note struct {
	Title       string    `json:"title"` //se formatea para que se mas legible el Json
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

//Creamos un map que va a  almacenas las notas que vamos a crear
// Tendra una key de tipo string como ID y recibe un dato del typo Note
var noteStore = make(map[string]Note)

//Creamos la variable de tipo int para el indice
var id int

// Creamos el primer router que recibira las petidiones de tipo GET
// Recibe como parametros un request y un responseWriter del paquete http
func GETnoteHandler(w http.ResponseWriter, r *http.Request) {
	// creamos un array que va a almacenar las notas desde la base de datos ficticia( noteStore )
	var notes []Note
	// iteramos notestore y cargamos las notas en el array notes
	for _, value := range noteStore {
		notes = append(notes, value)
	}
	// Configuramos la respuesta
	// con el metodo header.set configuramos el header del response
	w.Header().Set("Content-type", "application/json")
	// convertimos las notas a JSON
	j, err := json.Marshal(notes)
	// manejamos el error de forma basica para el ejemplo
	if err != nil {
		panic(err)
	}
	// configuramos el status de respuesta
	w.WriteHeader(200)
	//configuramos el body de la respuesta
	//y le pasamos el JSON  "j " que son nuestras notas
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
	// Creamso nuestro mux que será nuestro router maestro
	// El método StrictSlash se usa para reconocer o no la ruta que termina en un slash
	// si ponemos false soo tomará la ruta sin slash ej: /users/4
	// si ponemos  true tomará tambien como valida la ruta que termina en ej: /users/4/
	r := mux.NewRouter().StrictSlash(false)
	// definimos los controladores
	r.HandleFunc("/api/notes", GETnoteHandler).Methods("GET")
	r.HandleFunc("/api/notes/{id}", PUTnoteHandler).Methods("PUT")
	r.HandleFunc("/api/notes", POSTnoteHandler).Methods("POST")
	r.HandleFunc("/api/notes/{id}", DELETEnoteHandler).Methods("DELETE")

	// personalizamos nuestro servidor
	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("Servidor escuchando el puerto 8080")
	// arrancamos el servidor
	server.ListenAndServe()
}
