package data

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Recipe struct {
	Err  string `json:"error,omitempty"`
	id   int
	Name string `json:"name,omitempty"`
}

const (
	dbname   = "mealbuddy"
	host     = "localhost"
	password = "postgres"
	port     = 5432
	table    = "recipes"
	user     = "efrenulloa"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func connect() *sql.DB {
	conn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname,
	)

	db, connErr := sql.Open("postgres", conn)

	CheckError(connErr)

	pingErr := db.Ping()

	CheckError(pingErr)

	return db
}

func GetAllRecipes() []Recipe {
	db := connect()

	stmt := fmt.Sprintf("SELECT * FROM %s", table)

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
	wc := "%"
	name += wc
	stmt := fmt.Sprintf("SELECT * FROM %s WHERE name LIKE '%s'", table, name)

	rows, err := db.Query(stmt)

	CheckError(err)

	fmt.Println(rows)

	if rows == nil {
		fmt.Print("Rows are nil!")
		return Recipe{
			Err: fmt.Sprintf("No matching recipes with the name %s were found", name),
		}
	}

	var r Recipe

	for rows.Next() {
		e := rows.Scan(&r.id, &r.Name)
		CheckError(e)
	}

	defer rows.Close()
	defer db.Close()

	return r
}
