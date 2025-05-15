package routes

import (
	"blog/database"
	"blog/models"
	"blog/routerRegister"
	"blog/util"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

			findUser := models.User{}

			if result := db.Model(&models.User{}).Where(&models.User{Id: uint(id)}).First(&findUser); result.Error != nil {
				c.JSON(404, gin.H{
					"error": "not found user",
				})
				return
			}
			c.Render(200, util.JsonL(findUser))
		})

		user.POST("", func(c *gin.Context) {
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

			if currentUser.Role < models.UserRoleAdmin {
				c.JSON(403, gin.H{
					"error": "Permission denied",
				})
				return
			}

			var body struct {
				Email        string `json:"email"`
				Name         string `json:"name" binding:"required"`
				ProfilePhoto string `json:"profilePhoto"`
				Description  string `json:"description"`
				Avatar       string `json:"avatar"`
				Major        string `json:"major"`
			}

			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(400, gin.H{
					"error": "Invalid request body",
				})
				return
			}

			if body.Email != "" {
				var existingUser models.User
				if result := db.Model(&models.User{}).Where(&models.User{Email: body.Email}).First(&existingUser); result.Error == nil {
					c.JSON(400, gin.H{
						"error": "Email already exists",
					})
					return
				}
			}

			user := models.User{
				Email:        body.Email,
				Name:         &body.Name,
				ProfilePhoto: &body.ProfilePhoto,
				Description:  &body.Description,
				Avatar:       &body.Avatar,
				Major:        &body.Major,
			}
			if result := db.Create(&user); result.Error != nil {
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "User success",
				"user":    user,
			})
		})
	})
}
