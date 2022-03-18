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

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/recipe/{name}", recipe).Methods("GET")

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":1111", handler))
}

func index(rw http.ResponseWriter, req *http.Request) {
	recipes := data.GetAllRecipes()
	recipesJson, err := json.Marshal(recipes)
	if err != nil {
		panic(err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(recipesJson)
}

func recipe(rw http.ResponseWriter, req *http.Request) {
	name := mux.Vars(req)["name"]
	r := data.GetRecipe(name)
	recipeJson, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(recipeJson)
}
