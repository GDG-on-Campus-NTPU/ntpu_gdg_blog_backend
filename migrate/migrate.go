package main

import (
	"fmt"

	"ntpu_gdg.org/blog/database"
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
