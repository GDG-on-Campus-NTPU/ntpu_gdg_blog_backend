package main

import (
	"fmt"

	"ntpu_gdg.org/blog/database"

	_ "ntpu_gdg.org/blog/models"
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
