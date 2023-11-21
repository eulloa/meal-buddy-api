package data

import (
	"database/sql"
	"fmt"
	"log"
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

	rows, queryErr := db.Query("SELECT id, name, ingredients, instructions FROM recipes r INNER JOIN ingredients ing ON r.id = ing.recipe_id INNER JOIN instructions ins ON r.id = ins.recipe_id")

	CheckError(queryErr)

	rs := make([]Recipe, 0)

	for rows.Next() {
		var r Recipe
		e := rows.Scan(&r.id, &r.Name, (*pq.StringArray)(&r.Ingredients), (*pq.StringArray)(&r.Instructions))
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

func CreateRecipeList(recipesInList int) []Recipe {
	recipes := GetAllRecipes()
	length := len(recipes)
	var list []Recipe

	switch {
	case recipesInList > length:
		log.Println("Warning: The number of requested recipes in the list is greater than the total number of recipes!")
		break
	case recipesInList <= 0:
		fmt.Println("Warning: Number of recipes in the list may not be less than or equal to 0!")
		break
	case recipesInList == length:
		return recipes
	default:
		for recipesInList > 0 {
			rand := rand.Intn(len(recipes))
			list = append(list, recipes[rand])
			//remove random recipe from recipes so we don't have duplicates
			recipes = append(recipes[:rand], recipes[rand+1:]...)
			recipesInList -= 1
		}
	}

	return list
}

func AddRecipe(res map[string]interface{}) int {
	sanitize(res)

	ins := res["Instructions"].([]interface{})
	instructions := make([]string, 0)

	ings := res["Ingredients"].([]interface{})
	ingredients := make([]string, 0)

	for _, val := range ins {
		if len(val.(string)) > 0 {
			instructions = append(instructions, val.(string))
		}
	}

	for _, val := range ings {
		if len(val.(string)) > 0 {
			ingredients = append(ingredients, val.(string))
		}
	}

	r := Recipe{
		Name:         res["Name"].(string),
		Ingredients:  ingredients,
		Instructions: instructions,
	}

	db := connect()

	stmt, err := db.Prepare("INSERT INTO recipes (name) VALUES ($1)")

	CheckError(err)

	_, execErr := stmt.Exec(r.Name)

	CheckError(execErr)

	idQuery, idErr := db.Prepare("SELECT id FROM recipes WHERE name = $1")

	CheckError(idErr)

	row := idQuery.QueryRow(r.Name)

	var id int
	scanErr := row.Scan(&id)

	CheckError(scanErr)

	ingsStmt, ingsErr := db.Prepare("INSERT INTO ingredients (ingredients, recipe_id) VALUES ($1, $2)")
	insStmt, insErr := db.Prepare("INSERT INTO instructions (instructions, recipe_id) VALUES ($1, $2)")

	CheckError(ingsErr)
	CheckError(insErr)

	_, ingsExecErr := ingsStmt.Exec(pq.Array(r.Ingredients), id)
	CheckError(ingsExecErr)

	_, insExecErr := insStmt.Exec(pq.Array(r.Instructions), id)
	CheckError(insExecErr)

	return id
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if item == v {
			return true
		}
	}

	return false
}

func sanitize(res map[string]interface{}) {
	required := []string{"Name", "Ingredients", "Instructions"}
	validKeys := make([]string, 0)
	invalidKeys := make([]string, 0)

	for k := range res {
		c := contains(required, k)

		if !c {
			invalidKeys = append(invalidKeys, k)
		} else {
			validKeys = append(validKeys, k)
		}
	}

	if len(validKeys) != len(required) {
		log.Fatal("Error: Post body data does not contain all required keys!")
	}

	if len(invalidKeys) > 0 {
		log.Fatal("Error: Post body data contains redundant/illegal keys!")
	}
}
