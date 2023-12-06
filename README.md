## What is the meal-buddy-api? 👨‍🍳

A RESTful CRUD app built with [Golang](https://go.dev/dl/) (go1.21.5 darwin/amd64), [PostgreSQL](https://www.postgresql.org/) (14.2) and [gorilla/mux](https://github.com/gorilla/mux) that provides a few personal favorite recipes to simplify the chore of creating a list of meals for a given week.

My wife and I 🧔💁🏻‍♀️ are busy people and often don't have enough bandwidth to create a list of meals we want to eat during the week. This simple API allows us to feed it a set of reliable and delicious meals that we've enjoyed in the past and have it generate a random list of meals, depending on how many we're looking for that week. If you want the meal buddy to return a list of 3 recipes, simply do a GET against the `/recipe/list/{number}` endpoint and make sure to include a path variable for how many total meals you'd like it to generate, it's that easy!

### How to run the meal-buddy-api locally

Once you have Golang and PostgreSQL installed, pull down the meal-buddy-api and run the `db-init.sql` script to create a database and the tables as well as seed it with some initial recipes.

`psql -h localhost -U {your postgres username} -f scripts/db-init.sql`

_note that the `db-init.sql` script may require you to add execute permissions; in the event this is the case, run `chmod u+wxr scripts/db-init.sql` to grant the file owner permissions to read, write or execute the file_

With the database created and seeded, navigate to the project root in a terminal window and use the following command to run the application.

`go run main.go`

### Generate a list of n recipes

```
curl -X GET http://localhost:1111/recipe/list/{number}
```

### Get a specific recipe by id

```
curl -X GET http://localhost:1111/recipe/{id}
```

### Add a new recipe

```
curl -X POST http://localhost:1111/recipe/add -d '{ "Description": "Recipe description", "Image": "recipe-name.png", "Ingredients": ["ingredient one", "ingredient two"], "Instructions": ["instruction numero uno", "instruction numero dos"], "Name": "Recipe name", "Url": "https://myrecipe.com/recipe" }' -H "Content-Type: application/json"
```

### Delete a recipe

```
curl -X DELETE http://localhost:1111/recipe/delete/{id}
```

_Coming Soon: a UI to consume, display and allow for the modification, creation, deletion and updating of recipes using the endpoints exposed by the meal-buddy-api._

Bon appetite! 😋
