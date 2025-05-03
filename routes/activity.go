package routes

import (
	"blog/database"
	"blog/models"
	"blog/routerRegister"
	"blog/util"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func init() {
	routerRegister.Register = append(routerRegister.Register, func(rg *gin.RouterGroup) {
		activity := rg.Group("/activity")

		activity.POST("", func(c *gin.Context) {
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
			if currentUser.Role < models.UserRoleUploader {
				c.JSON(403, gin.H{
					"error": "Permission denied",
				})
				return
			}

			var body struct {
				Thumbnail   string    `json:"thumbnail"`
				Title       string    `json:"title"`
				Date        time.Time `json:"date"`
				Description string    `json:"description"`
				Tags        []string  `json:"tags"`
			}

			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(400, gin.H{
					"error": "Invalid request body",
				})
				return
			}

			tagsJson, err := util.ToDataTypeJSON(body.Tags)

			if err != nil {
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}

			activity := models.Activity{
				Thumbnail:   body.Thumbnail,
				Title:       body.Title,
				Date:        body.Date,
				Description: body.Description,
				Tags:        tagsJson,
				UserId:      currentUser.Id,
			}

			if result := db.Model(&models.Activity{}).Create(&activity); result.Error != nil {
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}

			c.JSON(200, gin.H{
				"message":    "success",
				"activityId": activity.Id,
			})
		})

		activity.GET("", func(c *gin.Context) {
			db := database.GetDB(c)
			startStr := c.Query("startDate")
			endStr := c.Query("endDate")
			limitStr := c.Query("num")

			startDate, err := time.Parse(time.RFC3339, startStr)

			if err != nil || startDate.IsZero() {
				c.JSON(400, gin.H{
					"error": "Invalid startDate format",
				})
				return
			}

			endDate, err := time.Parse(time.RFC3339, endStr)

			if err != nil || endDate.IsZero() {
				c.JSON(400, gin.H{
					"error": "Invalid endDate format",
				})
				return
			}

			limit, err := strconv.Atoi(limitStr)

			if err != nil || limit <= 0 {
				c.JSON(400, gin.H{
					"error": "Invalid limit",
				})
				return
			}

			activitys := []models.Activity{}

			if result := db.Model(&models.Activity{}).Where("date >= ? AND date <= ?", startDate, endDate).Limit(limit).Order("date asc").Find(&activitys); result.Error != nil {
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}

			c.Render(200, util.JsonL(activitys))
		})

		activity.DELETE(":id", func(c *gin.Context) {
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
			if currentUser.Role < models.UserRoleUploader {
				c.JSON(403, gin.H{
					"error": "Permission denied",
				})
				return
			}

			idStr := c.Param("id")

			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				c.JSON(400, gin.H{
					"error": "Invalid activity ID",
				})
				return
			}

			activity := models.Activity{}
			if result := db.Model(&models.Activity{}).Where(&models.Activity{Id: uint(id)}).First(&activity); result.Error != nil {
				c.JSON(404, gin.H{
					"error": "not found activity",
				})
				return
			}

			if currentUser.Role < models.UserRoleAdmin && activity.UserId != currentUser.Id {
				c.JSON(403, gin.H{
					"error": "Permission denied",
				})
				return
			}

			if result := db.Model(&models.Activity{}).Where(&models.Activity{Id: uint(id)}).Delete(&activity); result.Error != nil {
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}

			c.JSON(200, gin.H{
				"message": "success",
			})
		})

		activity.PUT(":id", func(c *gin.Context) {
			db := database.GetDB(c)
			session := sessions.Default(c)
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

			if currentUser.Role < models.UserRoleUploader {
				c.JSON(403, gin.H{
					"error": "Permission denied",
				})
				return
			}

			idStr := c.Param("id")
			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				c.JSON(400, gin.H{
					"error": "Invalid activity ID",
				})
				return
			}

			activity := models.Activity{}

			if result := db.Model(&models.Activity{}).Where(&models.Activity{Id: uint(id)}).First(&activity); result.Error != nil {
				c.JSON(404, gin.H{
					"error": "not found activity",
				})
				return
			}

			if currentUser.Role < models.UserRoleAdmin && activity.UserId != currentUser.Id {
				c.JSON(403, gin.H{
					"error": "Permission denied",
				})
				return
			}

			var body struct {
				Thumbnail   string    `json:"thumbnail"`
				Title       string    `json:"title"`
				Date        time.Time `json:"date"`
				Description string    `json:"description"`
				Tags        []string  `json:"tags"`
			}

			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(400, gin.H{
					"error": "Invalid request body",
				})
				return
			}

			tagsJson, err := util.ToDataTypeJSON(body.Tags)
			if err != nil {
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}

			activity.Thumbnail = body.Thumbnail
			activity.Title = body.Title
			activity.Date = body.Date
			activity.Description = body.Description
			activity.Tags = tagsJson

			if result := db.Model(&models.Activity{}).Where(&models.Activity{Id: uint(id)}).Updates(&activity); result.Error != nil {
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}

			c.JSON(200, gin.H{
				"message": "success",
			})
		})

	})
}
