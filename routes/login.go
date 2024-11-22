package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ntpu_gdg.org/blog/auth"
)

func AddLoginRoutes(rg *gin.RouterGroup) {
	login := rg.Group("/login")

	login.GET("", func(c *gin.Context) {
		// TODO state should be generated and stored in the session
		c.Redirect(http.StatusFound, auth.GoogleOauthConfig.AuthCodeURL("state"))
	})

	login.GET("/google/callback", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Google Callback",
		})
	})
}
