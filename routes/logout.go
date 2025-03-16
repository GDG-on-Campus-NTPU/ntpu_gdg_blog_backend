package routes

import (
	"blog/routerRegister"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func init() {
	routerRegister.Register = append(routerRegister.Register, func(rg *gin.RouterGroup) {
		logout := rg.Group("/logout")

		logout.POST("", func(c *gin.Context) {
			session := sessions.Default(c)

			session.Delete("Email")

			session.Save()

			c.JSON(200, gin.H{
				"message": "Logged out",
			})
		})
	})
}
