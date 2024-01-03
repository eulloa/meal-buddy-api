package main

import (
	"encoding/json"
	"fmt"
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
	router.HandleFunc("/recipe/update/{id}", update).Methods("PUT")
	router.HandleFunc("/", index).Methods("GET")

	handler := cors.Default().Handler(router)
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, handler))
}

func index(rw http.ResponseWriter, req *http.Request) {
	db := data.Connect()

	r := new(data.Recipe)

	recipes := r.GetAllRecipes(db)
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
	r := new(data.Recipe)

	id, rErr := r.AddRecipe(db, res)

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

	if convErr != nil {
		e := data.ErrorString{
			Error: fmt.Sprintf(
				"Unable to parse path param '%s'. Pass a valid unit instead", mux.Vars(req)["id"],
			),
		}
		j, _ := json.Marshal(e)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(j)
		return
	}

	db := data.Connect()

	r := new(data.Recipe)
	recipe, err := r.GetRecipe(db, id)

	if err != nil {
		j, _ := json.Marshal(err)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusNotFound)
		rw.Write(j)
		return
	}

	recipeJson, _ := json.Marshal(recipe)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(recipeJson)
}

func createList(rw http.ResponseWriter, req *http.Request) {
	numOfRecipes, convError := strconv.Atoi(mux.Vars(req)["number"])

	if convError != nil {
		e := data.ErrorString{
			Error: fmt.Sprintf(
				"Unable to parse path param '%s'. Pass a valid uint instead", mux.Vars(req)["number"],
			),
		}
		j, _ := json.Marshal(e)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(j)
		return
	}

	db := data.Connect()

	r := new(data.Recipe)

	recipes, err := r.CreateRecipeList(db, numOfRecipes)

	if err != nil {
		j, _ := json.Marshal(err)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(j)
		return
	}

	rJson, _ := json.Marshal(recipes)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(rJson)
}

func delete(rw http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])

	if err != nil {
		e := data.ErrorString{
			Error: fmt.Sprintf(
				"Unable to parse path param '%s'. Pass a valid uint instead", mux.Vars(req)["id"],
			),
		}
		j, _ := json.Marshal(e)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(j)
		return
	}

	db := data.Connect()
	r := new(data.Recipe)

	dErr := r.DeleteRecipe(db, id)

	if dErr != nil {
		j, _ := json.Marshal(dErr)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(j)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func update(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNotImplemented)

	e := struct {
		Message string
	}{
		Message: "Method not implemented",
	}

	m, _ := json.Marshal(e)

	rw.Write(m)
}
