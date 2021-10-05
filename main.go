package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/random", random).Methods("GET")

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":1111", handler))
}

func index(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "index page")
}

func random(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Random!")
}
