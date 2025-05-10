package routes

import (
	"blog/database"
	"blog/models"
	"blog/routerRegister"
	"blog/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

func init() {
	routerRegister.Register = append(routerRegister.Register, func(rg *gin.RouterGroup) {
		user := rg.Group("/user")

		user.GET(":id", func(c *gin.Context) {
			session := sessions.Default(c)
			db := database.GetDB(c)

			email, ok := session.Get("Email").(string)

			if !ok {
				c.JSON(401, gin.H{
					"error": "Not logged in",
				})
				return
			}

			currentUser := models.User{}

			if result := db.Model(&models.User{}).Where(&models.User{Email: email}).First(&currentUser); result.Error != nil {
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}

			idStr := c.Param("id")

			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				c.JSON(400, gin.H{
					"error": "Invalid user ID",
				})
				return
			}

			if currentUser.Role < models.UserRoleAdmin {
				if uint64(currentUser.Id) != id {
					c.JSON(401, gin.H{
						"error": "Not authorized",
					})
					return
				}
			}

			c.Render(200, util.JsonL(currentUser))
		})

	})
}
