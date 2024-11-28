package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"ntpu_gdg.org/blog/routerRegister"
)

func init() {
	routerRegister.Register = append(routerRegister.Register, func(rg *gin.RouterGroup) {
		user := rg.Group("/profile")
		user.GET("", func(c *gin.Context) {
			session := sessions.Default(c)

			name, ok := session.Get("Name").(string)

			if !ok {
				c.JSON(401, gin.H{
					"error": "Not logged in",
				})
				return
			}

			email, ok := session.Get("Email").(string)

			if !ok {
				c.JSON(401, gin.H{
					"error": "Not logged in",
				})
				return
			}

			c.JSON(200, gin.H{
				"name":  name,
				"email": email,
			})
		})
	})
}
