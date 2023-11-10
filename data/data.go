package data

import (
	"database/sql"
	"fmt"
	"math/rand"

	"strconv"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var vars = getVars()

type Recipe struct {
	id int
	// Ingredients Ingredients `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Name         string   `json:"name,omitempty"`
}

type Instructions struct {
	instructions_id int
	Instructions    []byte
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
	var i Instructions

	// TODO: separate instructions query logic into reusable func
	for rows.Next() {
		// TODO: add additional meal data (ingredients, instructions, etc)
		// e := rows.Scan(&r.id, &r.Ingredients, &r.Instructions, &r.Name)
		e := rows.Scan(&r.id, &r.Name)
		CheckError(e)

		// query for associated instructions
		instructionQuery := fmt.Sprintf("SELECT * FROM instructions WHERE instructions_id = %d", r.id)
		instructionRow, instructionsErr := db.Query(instructionQuery)
		CheckError(instructionsErr)

		for instructionRow.Next() {
			iErr := instructionRow.Scan(&i.instructions_id, &i.Instructions)
			CheckError(iErr)

			// convert byte array to string
			instructions := string(i.Instructions[:])
			ia := strings.Split(instructions, ", ")

			r.Instructions = ia
		}

		instructionRow.Close()
	}

	defer rows.Close()
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
