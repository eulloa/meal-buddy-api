package data

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Recipe struct {
	id   int
	Name string `json:"name"`
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

func GetRecipes() []Recipe {
	db := connect()

	stmt := fmt.Sprintf("SELECT * FROM %s", table)

	rows, err := db.Query(stmt)

	CheckError(err)

	rs := make([]Recipe, 0)

	for rows.Next() {
		var r Recipe
		e := rows.Scan(&r.Name)
		CheckError(e)
		rs = append(rs, r)
	}

	defer rows.Close()
	defer db.Close()

	return rs
}
