package data

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	dbname   = "mealbuddy"
	host     = "localhost"
	password = "postgres"
	port     = 5432
	user     = "efrenulloa"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func Connect() {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, connErr := sql.Open("postgres", conn)

	CheckError(connErr)

	pingErr := db.Ping()

	CheckError(pingErr)

	defer db.Close()
}

// func ListRecipes() {
// 	Connect()

// }
