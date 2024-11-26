package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"ntpu_gdg.org/blog/routerRegister"
)

func init() {
	routerRegister.Register = append(routerRegister.Register, func(rg *gin.RouterGroup) {
		logout := rg.Group("/logout")

		logout.POST("", func(c *gin.Context) {
			session := sessions.Default(c)

			session.Delete("Name")
			session.Delete("Email")

			session.Save()

			c.JSON(200, gin.H{
				"message": "Logged out",
			})
		})
	})
}
