package data

import (
	"database/sql"
	"fmt"

	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var vars = getVars()

type Recipe struct {
	id int
	// Ingredients  []string `json:"ingredients,omitempty"`
	// Instructions []string `json:"instructions,omitempty"`
	Name string `json:"name,omitempty"`
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func getVars() map[string]string {
	var envs map[string]string

	envs, err := godotenv.Read(".env")

	CheckError(err)

	return envs
}

func connect() *sql.DB {
	port, err := strconv.Atoi(vars["PORT"])

	CheckError(err)

	conn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		vars["HOST"], port, vars["USER"], vars["PASSWORD"], vars["DBNAME"],
	)

	db, connErr := sql.Open("postgres", conn)

	CheckError(connErr)

	pingErr := db.Ping()

	CheckError(pingErr)

	return db
}

func GetAllRecipes() []Recipe {
	db := connect()

	stmt := fmt.Sprintf("SELECT * FROM %s", vars["TABLE"])

	rows, err := db.Query(stmt)

	CheckError(err)

	rs := make([]Recipe, 0)

	for rows.Next() {
		var r Recipe
		e := rows.Scan(&r.id, &r.Name)
		CheckError(e)
		rs = append(rs, r)
	}

	defer rows.Close()
	defer db.Close()

	return rs
}

func GetRecipe(name string) Recipe {
	db := connect()
	stmt := fmt.Sprintf("SELECT * FROM %s WHERE name = '%s'", vars["TABLE"], name)

	rows, err := db.Query(stmt)

	CheckError(err)

	var r Recipe

	for rows.Next() {
		// TODO: add additional meal data (ingredients, instructions, etc)
		// e := rows.Scan(&r.id, &r.Ingredients, &r.Instructions, &r.Name)
		e := rows.Scan(&r.id, &r.Name)
		CheckError(e)
	}

	defer rows.Close()
	defer db.Close()

	return r
}
