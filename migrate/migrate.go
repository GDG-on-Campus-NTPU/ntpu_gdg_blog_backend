package main

import (
	"fmt"

	"blog/database"

	_ "blog/models"
)

func main() {
	db, err := database.CreateClient()
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(database.ORMModels...)

	if err != nil {
		panic(err)
	}

	fmt.Print("Migration complete")
}
