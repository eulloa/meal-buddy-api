package data

import (
	"database/sql"
	"fmt"
	"math/rand"

	"strconv"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var vars = getVars()

type Recipe struct {
	Name         string `json:"name,omitempty"`
	id           int
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
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

	// TODO: INNER JOIN
	rows, queryErr := db.Query("SELECT * FROM recipes")

	CheckError(queryErr)

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

// TODO: handle empty result set
func GetRecipe(id int) Recipe {
	db := connect()
	stmt, prepareErr := db.Prepare("SELECT id, name, ingredients, instructions FROM recipes r INNER JOIN ingredients ing ON r.id = ing.recipe_id INNER JOIN instructions ins ON r.id = ins.recipe_id WHERE r.id = $1")

	CheckError(prepareErr)

	var r Recipe
	err := stmt.QueryRow(id).Scan(&r.id, &r.Name, (*pq.StringArray)(&r.Ingredients), (*pq.StringArray)(&r.Instructions))

	CheckError(err)

	defer db.Close()
	return r
}

func GetRandomRecipe() Recipe {
	recipes := GetAllRecipes()
	randomInt := rand.Intn(len(recipes))
	randomRecipe := recipes[randomInt]

	return Recipe{
		id:   randomRecipe.id,
		Name: randomRecipe.Name,
	}
}
