package main

import (
	"blog/env"
	"blog/routerRegister"
	_ "blog/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	cookieStore := cookie.NewStore([]byte(env.Getenv("AUTH_SECRET")))

	r.Use(sessions.Sessions("ginSession", cookieStore))

	api := r.Group("/api")

	for _, route := range routerRegister.Register {
		route(api)
	}

	// if env.Getenv("LOG_EXECUTION_ID") != "" {
	// 	//gcp cloud functions
	// 	functions.HTTP("HelloHTTP", r.Handler().ServeHTTP)
	// } else {
	// 	//local
	// 	r.Run(":8080")
	// }

	port := "8080"

	if env.Getenv("PORT") != "" {
		port = env.Getenv("PORT")
	}

	r.Run(":" + port)

}
