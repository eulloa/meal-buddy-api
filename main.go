package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/eulloa/meal-buddy/data"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()
	port := ":1111"

	router.HandleFunc("/recipe/random", random).Methods("GET")
	router.HandleFunc("/recipe/{name}", recipe).Methods("GET")
	router.HandleFunc("/recipe/add", add).Methods("GET")
	router.HandleFunc("/", index).Methods("GET")

	handler := cors.Default().Handler(router)
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, handler))
}

func index(rw http.ResponseWriter, req *http.Request) {
	recipes := data.GetAllRecipes()
	recipesJson, err := json.Marshal(recipes)

	data.CheckError(err)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(recipesJson)
}

// TODO: implement method
func add(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte("add"))
}

func random(rw http.ResponseWriter, req *http.Request) {
	randomRecipe := data.GetRandomRecipe()
	randomJson, err := json.Marshal(randomRecipe)

	data.CheckError(err)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(randomJson)
}

func recipe(rw http.ResponseWriter, req *http.Request) {
	name := mux.Vars(req)["name"]
	r := data.GetRecipe(name)
	recipeJson, err := json.Marshal(r)

	data.CheckError(err)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(recipeJson)
}
