package database

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"ntpu_gdg.org/blog/env"
)

var ORMModles = []any{}

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

func getDB(c *gin.Context) *gorm.DB {
	db, ok := c.Get("db")
	if !ok {
		db, err := CreateClient()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		}
		c.Set("db", db)
	}
	return db.(*gorm.DB)
}
