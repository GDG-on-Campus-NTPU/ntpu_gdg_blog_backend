package database

import (
	"blog/env"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ORMModels = []any{}

func CreateClient() (*gorm.DB, error) {
	connectString := env.Getenv("DATABASE_URL")
	if connectString == "" {
		panic("DATABASE_URL is not set in env")
	}
	db, err := gorm.Open(postgres.Open(connectString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

var dbInstance *gorm.DB

func GetDB(c *gin.Context) *gorm.DB {
	db := dbInstance
	if db == nil {
		var err error
		db, err = CreateClient()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		}
		c.Set("db", db)
	}
	return db
}
