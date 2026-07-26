package data

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// --- contains tests ---

func TestContains_Found(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}
	if !contains(slice, "banana") {
		t.Error("expected contains to return true for 'banana'")
	}
}

func TestContains_NotFound(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}
	if contains(slice, "grape") {
		t.Error("expected contains to return false for 'grape'")
	}
}

func TestContains_EmptySlice(t *testing.T) {
	if contains([]string{}, "apple") {
		t.Error("expected contains to return false for empty slice")
	}
}

func TestContains_EmptyString(t *testing.T) {
	slice := []string{"", "banana"}
	if !contains(slice, "") {
		t.Error("expected contains to return true for empty string in slice")
	}
}

func TestContains_SingleElement(t *testing.T) {
	if !contains([]string{"only"}, "only") {
		t.Error("expected contains to return true for single matching element")
	}
}

// --- sanitize tests ---

func validData() map[string]interface{} {
	return map[string]interface{}{
		"Description":  "A test recipe",
		"Image":        "test.png",
		"Ingredients":  []interface{}{"flour", "sugar"},
		"Instructions": []interface{}{"mix", "bake"},
		"Name":         "Test Recipe",
		"Url":          "https://example.com",
	}
}

func TestSanitize_ValidData(t *testing.T) {
	err := sanitize(validData())
	if err != nil {
		t.Errorf("expected nil error, got: %s", err.Error)
	}
}

func TestSanitize_MissingKey(t *testing.T) {
	data := validData()
	delete(data, "Name")

	err := sanitize(data)
	if err == nil {
		t.Fatal("expected error for missing 'Name' key, got nil")
	}
	if err.Error != "Post body data does not contain all required keys!" {
		t.Errorf("unexpected error message: %s", err.Error)
	}
}

func TestSanitize_ExtraKey(t *testing.T) {
	data := validData()
	data["Extra"] = "not allowed"

	err := sanitize(data)
	if err == nil {
		t.Fatal("expected error for extra key 'Extra', got nil")
	}
	if err.Error != "Post body data contains redundant/illegal keys!" {
		t.Errorf("unexpected error message: %s", err.Error)
	}
}

func TestSanitize_EmptyDescription(t *testing.T) {
	data := validData()
	data["Description"] = ""

	err := sanitize(data)
	if err == nil {
		t.Fatal("expected error for empty Description, got nil")
	}
	if err.Error != "Description must not be an empty string!" {
		t.Errorf("unexpected error message: %s", err.Error)
	}
}

func TestSanitize_EmptyImage(t *testing.T) {
	data := validData()
	data["Image"] = ""

	err := sanitize(data)
	if err == nil {
		t.Fatal("expected error for empty Image, got nil")
	}
	if err.Error != "Image must not be an empty string!" {
		t.Errorf("unexpected error message: %s", err.Error)
	}
}

func TestSanitize_EmptyName(t *testing.T) {
	data := validData()
	data["Name"] = ""

	err := sanitize(data)
	if err == nil {
		t.Fatal("expected error for empty Name, got nil")
	}
	if err.Error != "Name must not be an empty string!" {
		t.Errorf("unexpected error message: %s", err.Error)
	}
}

func TestSanitize_EmptyUrl(t *testing.T) {
	data := validData()
	data["Url"] = ""

	err := sanitize(data)
	if err == nil {
		t.Fatal("expected error for empty Url, got nil")
	}
	if err.Error != "Url must not be an empty string!" {
		t.Errorf("unexpected error message: %s", err.Error)
	}
}

func TestSanitize_EmptyIngredients(t *testing.T) {
	data := validData()
	data["Ingredients"] = []interface{}{}

	err := sanitize(data)
	if err == nil {
		t.Fatal("expected error for empty Ingredients, got nil")
	}
	if err.Error != "Ingredients must not be an empty array!" {
		t.Errorf("unexpected error message: %s", err.Error)
	}
}

func TestSanitize_EmptyInstructions(t *testing.T) {
	data := validData()
	data["Instructions"] = []interface{}{}

	err := sanitize(data)
	if err == nil {
		t.Fatal("expected error for empty Instructions, got nil")
	}
	if err.Error != "Instructions must not be an empty array!" {
		t.Errorf("unexpected error message: %s", err.Error)
	}
}

func TestSanitize_MissingMultipleKeys(t *testing.T) {
	data := map[string]interface{}{
		"Description": "desc",
		"Name":        "name",
	}

	err := sanitize(data)
	if err == nil {
		t.Fatal("expected error for missing multiple keys, got nil")
	}
	if err.Error != "Post body data does not contain all required keys!" {
		t.Errorf("unexpected error message: %s", err.Error)
	}
}

// --- recipeFromMap tests ---

func TestRecipeFromMap(t *testing.T) {
	data := map[string]interface{}{
		"Description":  "Delicious pasta",
		"Image":        "pasta.png",
		"Ingredients":  []interface{}{"pasta", "sauce", "cheese"},
		"Instructions": []interface{}{"boil water", "cook pasta", "add sauce"},
		"Name":         "Pasta",
		"Url":          "https://example.com/pasta",
	}

	r := recipeFromMap(data)

	if r.Name != "Pasta" {
		t.Errorf("expected Name 'Pasta', got '%s'", r.Name)
	}
	if r.Description != "Delicious pasta" {
		t.Errorf("expected Description 'Delicious pasta', got '%s'", r.Description)
	}
	if r.Image != "pasta.png" {
		t.Errorf("expected Image 'pasta.png', got '%s'", r.Image)
	}
	if r.Url != "https://example.com/pasta" {
		t.Errorf("expected Url 'https://example.com/pasta', got '%s'", r.Url)
	}
	if len(r.Ingredients) != 3 {
		t.Errorf("expected 3 ingredients, got %d", len(r.Ingredients))
	}
	if len(r.Instructions) != 3 {
		t.Errorf("expected 3 instructions, got %d", len(r.Instructions))
	}
	if r.Ingredients[0] != "pasta" || r.Ingredients[1] != "sauce" || r.Ingredients[2] != "cheese" {
		t.Errorf("unexpected ingredients: %v", r.Ingredients)
	}
	if r.Instructions[0] != "boil water" {
		t.Errorf("unexpected first instruction: %s", r.Instructions[0])
	}
}

func TestRecipeFromMap_FiltersEmptyStrings(t *testing.T) {
	data := map[string]interface{}{
		"Description":  "desc",
		"Image":        "img.png",
		"Ingredients":  []interface{}{"flour", "", "sugar"},
		"Instructions": []interface{}{"step one", "", "step three"},
		"Name":         "Recipe",
		"Url":          "https://example.com",
	}

	r := recipeFromMap(data)

	if len(r.Ingredients) != 2 {
		t.Errorf("expected 2 ingredients (empty filtered), got %d", len(r.Ingredients))
	}
	if len(r.Instructions) != 2 {
		t.Errorf("expected 2 instructions (empty filtered), got %d", len(r.Instructions))
	}
}

// --- CheckError tests ---

func TestCheckError_NilError(t *testing.T) {
	CheckError(nil)
}

func TestCheckError_NonNilError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected CheckError to panic on non-nil error, but did not panic")
		}
	}()

	CheckError(errors.New("test error"))
}

// --- GetRecipe tests (sqlmock) ---

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	db = mockDB
	return mockDB, mock
}

func TestGetRecipe_Success(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "image", "ingredients", "instructions", "url"}).
		AddRow(1, "Pasta", "Delicious pasta", "pasta.png",
			`{"flour","sauce"}`, `{"boil","serve"}`, "https://example.com")

	mock.ExpectPrepare("SELECT .+ FROM recipes").ExpectQuery().WithArgs(1).WillReturnRows(rows)

	recipe := new(Recipe)
	result, err := recipe.GetRecipe(1)

	if err != nil {
		t.Errorf("expected no error, got: %s", err.Error)
	}
	if result == nil {
		t.Fatal("expected recipe, got nil")
	}
	if result.Name != "Pasta" {
		t.Errorf("expected Name 'Pasta', got '%s'", result.Name)
	}
	if result.Id != 1 {
		t.Errorf("expected Id 1, got %d", result.Id)
	}
	if result.Url != "https://example.com" {
		t.Errorf("expected Url 'https://example.com', got '%s'", result.Url)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func TestGetRecipe_NotFound(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	mock.ExpectPrepare("SELECT .+ FROM recipes").ExpectQuery().WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	recipe := new(Recipe)
	result, err := recipe.GetRecipe(999)

	if result != nil {
		t.Error("expected nil result for non-existent recipe")
	}
	if err == nil {
		t.Fatal("expected error for non-existent recipe, got nil")
	}
	if err.Error != fmt.Sprintf("Unable to find recipe with id: %d", 999) {
		t.Errorf("unexpected error message: %s", err.Error)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

// --- DeleteRecipe tests (sqlmock) ---

func TestDeleteRecipe_Success(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	mock.ExpectPrepare("DELETE FROM recipes").ExpectExec().WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	recipe := new(Recipe)
	err := recipe.DeleteRecipe(1)

	if err != nil {
		t.Errorf("expected no error, got: %s", err.Error)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func TestDeleteRecipe_PrepareError(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	mock.ExpectPrepare("DELETE FROM recipes").
		WillReturnError(errors.New("prepare failed"))

	recipe := new(Recipe)
	err := recipe.DeleteRecipe(1)

	if err == nil {
		t.Fatal("expected error on prepare failure, got nil")
	}
	if err.Error != "There was an error preparing the delete recipe statement" {
		t.Errorf("unexpected error message: %s", err.Error)
	}
}

func TestDeleteRecipe_ExecError(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	mock.ExpectPrepare("DELETE FROM recipes").ExpectExec().WithArgs(1).
		WillReturnError(errors.New("exec failed"))

	recipe := new(Recipe)
	err := recipe.DeleteRecipe(1)

	if err == nil {
		t.Fatal("expected error on exec failure, got nil")
	}
	if err.Error != "There was an error executing the delete recipe statement" {
		t.Errorf("unexpected error message: %s", err.Error)
	}
}

// --- GetAllRecipes tests (sqlmock) ---

func TestGetAllRecipes_Success(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "image", "ingredients", "instructions", "url"}).
		AddRow(1, "Pasta", "Desc1", "img1.png", `{"flour"}`, `{"boil"}`, "https://a.com").
		AddRow(2, "Soup", "Desc2", "img2.png", `{"broth"}`, `{"simmer"}`, "https://b.com")

	mock.ExpectQuery("SELECT .+ FROM recipes").WillReturnRows(rows)

	recipe := new(Recipe)
	result := recipe.GetAllRecipes()

	if len(result) != 2 {
		t.Fatalf("expected 2 recipes, got %d", len(result))
	}
	if result[0].Name != "Pasta" {
		t.Errorf("expected first recipe Name 'Pasta', got '%s'", result[0].Name)
	}
	if result[1].Name != "Soup" {
		t.Errorf("expected second recipe Name 'Soup', got '%s'", result[1].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func TestGetAllRecipes_Empty(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "image", "ingredients", "instructions", "url"})

	mock.ExpectQuery("SELECT .+ FROM recipes").WillReturnRows(rows)

	recipe := new(Recipe)
	result := recipe.GetAllRecipes()

	if len(result) != 0 {
		t.Errorf("expected 0 recipes, got %d", len(result))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

// --- AddRecipe tests (sqlmock) ---

func TestAddRecipe_Success(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	data := validData()

	mock.ExpectPrepare("INSERT INTO recipes").ExpectExec().
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectPrepare("SELECT id FROM recipes").ExpectQuery().
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(11))

	mock.ExpectPrepare("INSERT INTO ingredients")
	mock.ExpectPrepare("INSERT INTO instructions")
	mock.ExpectExec("INSERT INTO ingredients").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO instructions").WillReturnResult(sqlmock.NewResult(0, 1))

	recipe := new(Recipe)
	id, err := recipe.AddRecipe(data)

	if err != nil {
		t.Errorf("expected no error, got: %s", err.Error)
	}
	if id != 11 {
		t.Errorf("expected id 11, got %d", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func TestAddRecipe_SanitizeError(t *testing.T) {
	mockDB, _ := setupMockDB(t)
	defer mockDB.Close()

	data := map[string]interface{}{
		"Name": "",
	}

	recipe := new(Recipe)
	id, err := recipe.AddRecipe(data)

	if err == nil {
		t.Fatal("expected sanitize error, got nil")
	}
	if id != 0 {
		t.Errorf("expected id 0 on error, got %d", id)
	}
}

func TestAddRecipe_PrepareInsertError(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	data := validData()

	mock.ExpectPrepare("INSERT INTO recipes").
		WillReturnError(errors.New("prepare failed"))

	recipe := new(Recipe)
	id, err := recipe.AddRecipe(data)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if id != 0 {
		t.Errorf("expected id 0 on error, got %d", id)
	}
	if err.Error != "System encountered an error preparing record to insert into the database" {
		t.Errorf("unexpected error message: %s", err.Error)
	}
}

func TestAddRecipe_ExecInsertError(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	data := validData()

	mock.ExpectPrepare("INSERT INTO recipes").ExpectExec().
		WillReturnError(errors.New("exec failed"))

	recipe := new(Recipe)
	id, err := recipe.AddRecipe(data)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if id != 0 {
		t.Errorf("expected id 0 on error, got %d", id)
	}
	if err.Error != "System encountered an error inserting record into the database" {
		t.Errorf("unexpected error message: %s", err.Error)
	}
}

func TestAddRecipe_SelectIdError(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	data := validData()

	mock.ExpectPrepare("INSERT INTO recipes").ExpectExec().
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectPrepare("SELECT id FROM recipes").
		WillReturnError(errors.New("prepare failed"))

	recipe := new(Recipe)
	id, err := recipe.AddRecipe(data)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if id != 0 {
		t.Errorf("expected id 0 on error, got %d", id)
	}
	if err.Error != "System encountered an error preparing the select recipe statement" {
		t.Errorf("unexpected error message: %s", err.Error)
	}
}

// --- UpdateRecipe tests (sqlmock) ---

func TestUpdateRecipe_Success(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	data := validData()

	mock.ExpectPrepare("UPDATE recipes")
	mock.ExpectPrepare("UPDATE ingredients")
	mock.ExpectPrepare("UPDATE instructions")
	mock.ExpectExec("UPDATE recipes").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE ingredients").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE instructions").WillReturnResult(sqlmock.NewResult(0, 1))

	recipe := new(Recipe)
	result, err := recipe.UpdateRecipe(1, data)

	if err != nil {
		t.Errorf("expected no error, got: %s", err.Error)
	}
	if result == nil {
		t.Fatal("expected recipe result, got nil")
	}
	if result.Name != "Test Recipe" {
		t.Errorf("expected Name 'Test Recipe', got '%s'", result.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func TestUpdateRecipe_SanitizeError(t *testing.T) {
	mockDB, _ := setupMockDB(t)
	defer mockDB.Close()

	data := map[string]interface{}{
		"Name": "",
	}

	recipe := new(Recipe)
	result, err := recipe.UpdateRecipe(1, data)

	if err == nil {
		t.Fatal("expected sanitize error, got nil")
	}
	if result != nil {
		t.Error("expected nil result on error")
	}
}

func TestUpdateRecipe_PrepareError(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	data := validData()

	mock.ExpectPrepare("UPDATE recipes").
		WillReturnError(errors.New("prepare failed"))

	recipe := new(Recipe)
	result, err := recipe.UpdateRecipe(1, data)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Error("expected nil result on error")
	}
}

func TestUpdateRecipe_ExecError(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	data := validData()

	mock.ExpectPrepare("UPDATE recipes").ExpectExec().
		WillReturnError(errors.New("exec failed"))

	recipe := new(Recipe)
	result, err := recipe.UpdateRecipe(1, data)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Error("expected nil result on error")
	}
}

// --- ErrorString tests ---

func TestErrorString_Struct(t *testing.T) {
	e := ErrorString{Error: "test error"}
	if e.Error != "test error" {
		t.Errorf("expected 'test error', got '%s'", e.Error)
	}
}
