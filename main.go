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
	router.HandleFunc("/recipe/add", add).Methods("POST")
	router.HandleFunc("/recipe/list/{number}", createList).Methods("GET")
	router.HandleFunc("/recipe/delete/{id}", delete).Methods("DELETE")
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

func add(rw http.ResponseWriter, req *http.Request) {
	var res map[string]interface{}
	json.NewDecoder(req.Body).Decode(&res)

	db := data.Connect()

	id, rErr := data.AddRecipe(db, res)

	if rErr != nil {
		j, _ := json.Marshal(rErr)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(j)
		return
	}

	idJson, _ := json.Marshal(id)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(idJson)
}

func recipe(rw http.ResponseWriter, req *http.Request) {
	id, convErr := strconv.Atoi(mux.Vars(req)["id"])

	data.CheckError(convErr)
	db := data.Connect()

	r, err := data.GetRecipe(db, id)

	if err != nil {
		j, _ := json.Marshal(err)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusNotFound)
		rw.Write(j)
		return
	}

	recipeJson, _ := json.Marshal(r)

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

func delete(rw http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])

	data.CheckError(err)

	data.DeleteRecipe(id)

	rw.WriteHeader(http.StatusNoContent)
}
