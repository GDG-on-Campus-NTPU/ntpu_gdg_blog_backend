package blog

import (
	"ntpu_gdg.org/blog/env"
	"ntpu_gdg.org/blog/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {

	r := gin.Default()

	cookieStore := cookie.NewStore([]byte(env.Getenv("AUTH_SECRET")))

	r.Use(sessions.Sessions("ginSession", cookieStore))

	api := r.Group("/api")

	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API Index",
		})
	})

	routes.AddLoginRoutes(api)
	routes.AddProfileRoutes(api)
	routes.AddLogoutRoutes(api)

	if env.Getenv("LOG_EXECUTION_ID") != "" {
		//gcp cloud functions
		functions.HTTP("HelloHTTP", r.Handler().ServeHTTP)
	} else {
		//local
		r.Run(":8080")
	}

}
