package migrate

import (
	"ntpu_gdg.org/blog/database"
)

func init() {
	db, err := database.CreateClient()
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(database.ORMModles...)

	if err != nil {
		panic(err)
	}
}
