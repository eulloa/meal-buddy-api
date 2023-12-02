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
	Description  string `json:"description"`
	Name         string `json:"name,omitempty"`
	id           int
	Image        string   `json:"image"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Url          string   `json:"url"`
}

type ErrorString struct {
	Error string
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

func Connect() *sql.DB {
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
	db := Connect()

	rows, queryErr := db.Query("SELECT id, name, description, image, ingredients, instructions, url FROM recipes r INNER JOIN ingredients ing ON r.id = ing.recipe_id INNER JOIN instructions ins ON r.id = ins.recipe_id")

	CheckError(queryErr)

	rs := make([]Recipe, 0)

	for rows.Next() {
		var r Recipe
		e := rows.Scan(&r.id, &r.Name, &r.Description, &r.Image, (*pq.StringArray)(&r.Ingredients), (*pq.StringArray)(&r.Instructions), &r.Url)
		CheckError(e)
		rs = append(rs, r)
	}

	defer rows.Close()
	defer db.Close()

	return rs
}

func GetRecipe(db *sql.DB, id int) (*Recipe, *ErrorString) {
	stmt, prepareErr := db.Prepare("SELECT id, name, description, image, ingredients, instructions, url FROM recipes r INNER JOIN ingredients ing ON r.id = ing.recipe_id INNER JOIN instructions ins ON r.id = ins.recipe_id WHERE r.id = $1")

	if prepareErr != nil {
		e := ErrorString{
			Error: fmt.Sprint("There was an error preparing the SELECT statement"),
		}
		return nil, &e
	}

	var r Recipe
	err := stmt.QueryRow(id).Scan(&r.id, &r.Name, &r.Description, &r.Image, (*pq.StringArray)(&r.Ingredients), (*pq.StringArray)(&r.Instructions), &r.Url)

	if err != nil {
		e := ErrorString{
			Error: fmt.Sprintf("Unable to find recipe with id: %d", id),
		}
		return nil, &e
	}

	defer db.Close()
	return &r, nil
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
		Description:  res["Description"].(string),
		Image:        res["Image"].(string),
		Ingredients:  ingredients,
		Instructions: instructions,
		Name:         res["Name"].(string),
		Url:          res["Url"].(string),
	}

	db := Connect()

	stmt, err := db.Prepare("INSERT INTO recipes (name, description, image, url) VALUES ($1, $2, $3, $4)")

	CheckError(err)

	_, execErr := stmt.Exec(r.Name, r.Description, r.Image, r.Url)

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

	defer db.Close()
	return id
}

func DeleteRecipe(id int) {
	db := Connect()

	stmt, err := db.Prepare("DELETE FROM recipes WHERE id = $1")

	CheckError(err)

	_, stmtErr := stmt.Exec(id)

	CheckError(stmtErr)

	defer db.Close()
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
	required := []string{"Description", "Image", "Ingredients", "Instructions", "Name", "Url"}
	validKeys := make([]string, 0)
	invalidKeys := make([]string, 0)

	for k, v := range res {
		if k == "Description" || k == "Image" || k == "Name" || k == "Url" {
			if v == "" {
				log.Fatalf("Error: %s must not be an empty string!", k)
			}
		}

		if k == "Ingredients" {
			vals := res["Ingredients"].([]interface{})
			if len(vals) == 0 {
				log.Fatal("Error: Ingredients must not be an empty array!")
			}
		}

		if k == "Instructions" {
			vals := res["Instructions"].([]interface{})
			if len(vals) == 0 {
				log.Fatal("Error: Instructions must not be an empty array!")
			}
		}

		// check all required keys (and no redundant ones) are passed
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
