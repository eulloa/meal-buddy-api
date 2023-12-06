package main

import (
	"testing"
)

// func TestGetRecipe(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("Unexpeced error: %s", err)
// 	}
// 	defer db.Close()

// 	mock.ExpectBegin()
// 	mock.ExpectPrepare(
// 		"SELECT id, name, description, image, ingredients, instructions, url FROM recipes r INNER JOIN ingredients ing ON r.id = ing.recipe_id INNER JOIN instructions ins ON r.id = ins.recipe_id WHERE r.id = $1",
// 	).ExpectQuery().WithArgs(1).WillReturnRows(&sqlmock.Rows{})
// 	mock.ExpectCommit()

// 	r := new(data.Recipe)
// 	if _, recipeErr := r.GetRecipe(db, 1); recipeErr != nil {
// 		t.Errorf("GetRecipe unexpected error: %s", recipeErr)
// 	}

// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("Unfulfilled expectations: %s", err)
// 	}
// }

func TestFoo(t *testing.T) {
	var a = 2
	var b = 2

	if a != b {
		t.Fail()
	}
}
