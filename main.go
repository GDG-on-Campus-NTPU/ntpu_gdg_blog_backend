package blog

import (
	"ntpu_gdg.org/blog/env"
	"ntpu_gdg.org/blog/routes"

	"github.com/gin-gonic/gin"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {

	r := gin.Default()

	api := r.Group("/api")

	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API Index",
		})
	})

	routes.AddLoginRoutes(api)

	routes.AddIndexRoutes(r.Group(""))

	if env.Getenv("LOG_EXECUTION_ID") != "" {
		//gcp cloud functions
		functions.HTTP("HelloHTTP", r.Handler().ServeHTTP)
	} else {
		//local
		r.Run()
	}

}

func main() {

}
