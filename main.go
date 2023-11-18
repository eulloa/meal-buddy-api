package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/eulloa/meal-buddy/data"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()
	port := ":1111"

	router.HandleFunc("/recipe/{id}", recipe).Methods("GET")
	router.HandleFunc("/recipe/add", add).Methods("GET")
	router.HandleFunc("/recipe/list/{number}", createList).Methods("GET")
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

func recipe(rw http.ResponseWriter, req *http.Request) {
	id, convErr := strconv.Atoi(mux.Vars(req)["id"])

	data.CheckError(convErr)

	r := data.GetRecipe(id)
	recipeJson, err := json.Marshal(r)

	data.CheckError(err)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(recipeJson)
}

func createList(rw http.ResponseWriter, req *http.Request) {
	numOfRecipes, convError := strconv.Atoi(mux.Vars(req)["number"])

	data.CheckError(convError)

	recipes := data.CreateRecipeList(numOfRecipes)
	rJson, err := json.Marshal(recipes)

	data.CheckError(err)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(rJson)
}
