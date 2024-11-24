package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AddProfileRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/profile")

	user.GET("", func(c *gin.Context) {
		session := sessions.Default(c)

		name, ok := session.Get("Name").(string)

		if !ok {
			c.AbortWithStatus(401)
			return
		}

		email, ok := session.Get("Email").(string)

		if !ok {
			c.AbortWithStatus(401)
			return
		}

		c.JSON(200, gin.H{
			"name":  name,
			"email": email,
		})
	})
}
