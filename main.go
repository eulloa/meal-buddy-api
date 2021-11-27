package main

import (
	"io"
	"log"
	"net/http"

	"github.com/eulloa/meal-buddy/data"
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
	r, err := data.GetRecipes()
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	io.WriteString(w, string(r))
}

func random(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Random!")
}
