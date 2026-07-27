package data

import (
	"database/sql"
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var (
	vars = getVars()
	db   *sql.DB
)

type Recipe struct {
	Description  string   `json:"description"`
	Name         string   `json:"name,omitempty"`
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Url          string   `json:"url"`
}

type IRecipe interface {
	AddRecipe(vals map[string]interface{}) (int, *ErrorString)
	CreateRecipeList(n int) (*[]Recipe, *ErrorString)
	DeleteRecipe(id int) *ErrorString
	GetAllRecipes() []Recipe
	GetRecipe(id int) (*Recipe, *ErrorString)
	UpdateRecipe(id int) (int, *ErrorString)
}

type ErrorString struct {
	Error string `json:"error"`
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func getVars() map[string]string {
	envs, err := godotenv.Read(".env")
	if err != nil {
		envs, err = godotenv.Read("../.env")
		CheckError(err)
	}

	return envs
}

func InitDB() {
	port, err := strconv.Atoi(vars["PORT"])
	CheckError(err)

	conn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		vars["HOST"], port, vars["USER"], vars["PASSWORD"], vars["DBNAME"],
	)

	var connErr error
	db, connErr = sql.Open("postgres", conn)
	CheckError(connErr)

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	pingErr := db.Ping()
	CheckError(pingErr)
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func (r Recipe) GetAllRecipes() []Recipe {
	rows, queryErr := db.Query("SELECT id, name, description, image, ingredients, instructions, url FROM recipes r INNER JOIN ingredients ing ON r.id = ing.recipe_id INNER JOIN instructions ins ON r.id = ins.recipe_id")

	CheckError(queryErr)

	rs := make([]Recipe, 0)

	for rows.Next() {
		var r Recipe
		e := rows.Scan(&r.Id, &r.Name, &r.Description, &r.Image, (*pq.StringArray)(&r.Ingredients), (*pq.StringArray)(&r.Instructions), &r.Url)
		CheckError(e)
		rs = append(rs, r)
	}

	defer rows.Close()

	return rs
}

func (r Recipe) GetRecipe(id int) (*Recipe, *ErrorString) {
	stmt, prepareErr := db.Prepare("SELECT id, name, description, image, ingredients, instructions, url FROM recipes r INNER JOIN ingredients ing ON r.id = ing.recipe_id INNER JOIN instructions ins ON r.id = ins.recipe_id WHERE r.id = $1")

	if prepareErr != nil {
		return nil, &ErrorString{
			Error: "There was an error preparing the SELECT statement",
		}
	}

	err := stmt.QueryRow(id).Scan(&r.Id, &r.Name, &r.Description, &r.Image, (*pq.StringArray)(&r.Ingredients), (*pq.StringArray)(&r.Instructions), &r.Url)

	if err != nil {
		return nil, &ErrorString{
			Error: fmt.Sprintf("Unable to find recipe with id: %d", id),
		}
	}

	return &r, nil
}

func (r Recipe) CreateRecipeList(recipesInList int) (*[]Recipe, *ErrorString) {
	recipes := r.GetAllRecipes()
	length := len(recipes)
	var list []Recipe

	switch {
	case recipesInList > length:
		return nil, &ErrorString{
			Error: "The number of requested recipes in the list is greater than the total number of recipes!",
		}
	case recipesInList <= 0:
		return nil, &ErrorString{
			Error: "Number of recipes in the list may not be less than or equal to 0!",
		}
	case recipesInList == length:
		return &recipes, nil
	default:
		for recipesInList > 0 {
			rand := rand.Intn(len(recipes))
			list = append(list, recipes[rand])
			//remove random recipe from recipes so we don't have duplicates
			recipes = append(recipes[:rand], recipes[rand+1:]...)
			recipesInList -= 1
		}
	}

	return &list, nil
}

func (r Recipe) AddRecipe(data map[string]interface{}) (*Recipe, *ErrorString) {
	err := sanitize(data)

	if err != nil {
		return nil, err
	}

	r = recipeFromMap(data)

	stmt, prepareErr := db.Prepare("INSERT INTO recipes (name, description, image, url) VALUES ($1, $2, $3, $4)")

	if prepareErr != nil {
		return nil, &ErrorString{
			Error: "System encountered an error preparing record to insert into the database",
		}
	}

	defer stmt.Close()

	_, execErr := stmt.Exec(r.Name, r.Description, r.Image, r.Url)

	if execErr != nil {
		return nil, &ErrorString{
			Error: "System encountered an error inserting record into the database",
		}
	}

	idQuery, idErr := db.Prepare("SELECT id FROM recipes ORDER BY id DESC LIMIT 1")

	if idErr != nil {
		return nil, &ErrorString{
			Error: "System encountered an error preparing the select recipe statement",
		}
	}

	row := idQuery.QueryRow()

	var id int
	scanErr := row.Scan(&id)

	if scanErr != nil {
		return nil, &ErrorString{
			Error: fmt.Sprintf("System encountered an error scanning row with recipe id: %d", id),
		}
	}

	ingsStmt, ingsErr := db.Prepare("INSERT INTO ingredients (ingredients, recipe_id) VALUES ($1, $2)")
	insStmt, insErr := db.Prepare("INSERT INTO instructions (instructions, recipe_id) VALUES ($1, $2)")

	if ingsErr != nil {
		return nil, &ErrorString{
			Error: "System encountered an error preparing insert into ingredients table",
		}
	}

	if insErr != nil {
		return nil, &ErrorString{
			Error: "System encountered an error prepating insert into instructions table",
		}
	}

	_, ingsExecErr := ingsStmt.Exec(pq.Array(r.Ingredients), id)

	if ingsExecErr != nil {
		return nil, &ErrorString{
			Error: fmt.Sprintf("System encountered an error inserting ingredients associated with recipe id: %d", id),
		}
	}

	_, insExecErr := insStmt.Exec(pq.Array(r.Instructions), id)

	if insExecErr != nil {
		return nil, &ErrorString{
			Error: fmt.Sprintf("System encountered an error inserting instructions associated with recipe id: %d", id),
		}
	}

	return &r, nil
}

func (r Recipe) DeleteRecipe(id int) *ErrorString {
	stmt, err := db.Prepare("DELETE FROM recipes WHERE id = $1")

	if err != nil {
		return &ErrorString{
			Error: "There was an error preparing the delete recipe statement",
		}
	}

	_, stmtErr := stmt.Exec(id)

	if stmtErr != nil {
		return &ErrorString{
			Error: "There was an error executing the delete recipe statement",
		}
	}

	return nil
}

func (r Recipe) UpdateRecipe(id int, data map[string]interface{}) (*Recipe, *ErrorString) {
	sErr := sanitize(data)

	if sErr != nil {
		return nil, &ErrorString{
			Error: sErr.Error,
		}
	}

	r = recipeFromMap(data)

	stmt, updateErr := db.Prepare("UPDATE recipes SET (name, description, image, url) = ($1, $2, $3, $4) WHERE id = $5")
	ingStmt, ingErr := db.Prepare("UPDATE ingredients SET ingredients = $1 WHERE recipe_id = $2")
	insStmt, insErr := db.Prepare("UPDATE instructions SET instructions = $1 WHERE recipe_id = $2")

	if updateErr != nil {
		return nil, &ErrorString{
			Error: updateErr.Error(),
		}
	}

	if ingErr != nil {
		return nil, &ErrorString{
			Error: ingErr.Error(),
		}
	}

	if insErr != nil {
		return nil, &ErrorString{
			Error: insErr.Error(),
		}
	}

	_, execErr := stmt.Exec(&r.Name, &r.Description, &r.Image, &r.Url, id)
	_, ingExecErr := ingStmt.Exec(pq.Array(r.Ingredients), id)
	_, insExecErr := insStmt.Exec(pq.Array(r.Instructions), id)

	if execErr != nil {
		return nil, &ErrorString{
			Error: execErr.Error(),
		}
	}

	if ingExecErr != nil {
		return nil, &ErrorString{
			Error: ingExecErr.Error(),
		}
	}

	if insExecErr != nil {
		return nil, &ErrorString{
			Error: insExecErr.Error(),
		}
	}

	return &r, nil
}

func recipeFromMap(data map[string]interface{}) Recipe {
	ins := data["Instructions"].([]interface{})
	instructions := make([]string, 0)

	ings := data["Ingredients"].([]interface{})
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

	return Recipe{
		Description:  data["Description"].(string),
		Ingredients:  ingredients,
		Image:        data["Image"].(string),
		Instructions: instructions,
		Name:         data["Name"].(string),
		Url:          data["Url"].(string),
	}
}

func sanitize(res map[string]interface{}) *ErrorString {
	required := []string{"Description", "Image", "Ingredients", "Instructions", "Name", "Url"}
	validKeys := make([]string, 0)
	invalidKeys := make([]string, 0)

	for k, v := range res {
		if k == "Description" || k == "Image" || k == "Name" || k == "Url" {
			if v == "" {
				return &ErrorString{
					Error: fmt.Sprintf("%s must not be an empty string!", k),
				}
			}
		}

		if k == "Ingredients" {
			vals := res["Ingredients"].([]interface{})
			if len(vals) == 0 {
				return &ErrorString{
					Error: "Ingredients must not be an empty array!",
				}
			}
		}

		if k == "Instructions" {
			vals := res["Instructions"].([]interface{})
			if len(vals) == 0 {
				return &ErrorString{
					Error: "Instructions must not be an empty array!",
				}
			}
		}

		// check all required keys (and no redundant ones) are passed
		c := slices.Contains(required, k)

		if !c {
			invalidKeys = append(invalidKeys, k)
		} else {
			validKeys = append(validKeys, k)
		}
	}

	if len(validKeys) != len(required) {
		return &ErrorString{
			Error: "Post body data does not contain all required keys!",
		}
	}

	if len(invalidKeys) > 0 {
		return &ErrorString{
			Error: "Post body data contains redundant/illegal keys!",
		}
	}

	return nil
}
