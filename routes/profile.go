package routes

import (
	"blog/database"
	"blog/models"
	"blog/routerRegister"
	"blog/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func init() {
	routerRegister.Register = append(routerRegister.Register, func(rg *gin.RouterGroup) {
		profile := rg.Group("/profile")
		profile.GET("", func(c *gin.Context) {
			session := sessions.Default(c)
			db := database.GetDB(c)

			email, ok := session.Get("Email").(string)

			if !ok {
				c.JSON(401, gin.H{
					"error": "Not logged in",
				})
				return
			}

			user := models.User{Email: email}

			var findUser models.User

			if result := db.Model(&models.User{}).Where(&user).First(&findUser); result.Error != nil {
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}

			c.Render(200, util.JsonL(findUser))
		})

	})
}
