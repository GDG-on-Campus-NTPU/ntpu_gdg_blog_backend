package routes

import (
	"encoding/json" // 引入 encoding/json 處理圖片陣列
	"fmt"
	"strconv"
	"time"

	"blog/database"
	"blog/models"
	"blog/routerRegister"
	"blog/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func init() {
	routerRegister.Register = append(routerRegister.Register, func(rg *gin.RouterGroup) {
		project := rg.Group("/project")

		project.POST("", func(c *gin.Context) {
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
					"error": "internal server error - failed to get user",
				})
				return
			}

			if currentUser.Role < models.UserRoleUploader {
				c.JSON(403, gin.H{
					"error": "Permission denied. Only uploaders or admins can create projects.",
				})
				return
			}

			var body struct {
				Title       string    `json:"title"`
				Thumbnail   string    `json:"thumbnail"`
				Tags        []string  `json:"tags"`
				Description string    `json:"description"`
				Images      []string  `json:"images"` // urls
				StartDate   time.Time `json:"startDate"`
				EndDate     time.Time `json:"endDate"`
				Members     []uint    `json:"members"`
			}

			if err := c.ShouldBindJSON(&body); err != nil {
				fmt.Println("Error binding JSON for project POST:", err)
				c.JSON(400, gin.H{
					"error": fmt.Sprintf("Bad request: %v", err),
				})
				return
			}

			tagsJSON, err := json.Marshal(body.Tags)
			if err != nil {
				fmt.Println("Error marshalling tags for project POST:", err)
				c.JSON(500, gin.H{
					"error": "internal server error - marshalling tags",
				})
				return
			}

			imagesJSON, err := json.Marshal(body.Images)
			if err != nil {
				fmt.Println("Error marshalling images for project POST:", err)
				c.JSON(500, gin.H{
					"error": "internal server error - marshalling images",
				})
				return
			}

			project := models.Project{
				Title:       body.Title,
				Thumbnail:   body.Thumbnail,
				Tags:        datatypes.JSON(tagsJSON),
				Description: body.Description,
				Images:      datatypes.JSON(imagesJSON),
				StartDate:   body.StartDate,
				EndDate:     body.EndDate,
				UserId:      currentUser.Id,
			}

			tx := db.Begin()
			if tx.Error != nil {
				fmt.Println("Error starting transaction for project POST:", tx.Error)
				c.JSON(500, gin.H{"error": "internal server error - transaction start"})
				return
			}

			if result := tx.Create(&project); result.Error != nil {
				tx.Rollback()
				fmt.Println("Error creating project:", result.Error)
				c.JSON(500, gin.H{
					"error": "internal server error - creating project",
				})
				return
			}

			if len(body.Members) > 0 {
				var members []models.User
				if result := tx.Find(&members, body.Members); result.Error != nil {
					tx.Rollback()
					fmt.Println("Error finding members for project POST:", result.Error)
					c.JSON(500, gin.H{
						"error": "internal server error - finding members",
					})
					return
				}

				if len(members) != len(body.Members) {
					tx.Rollback()
					c.JSON(400, gin.H{
						"error": "Invalid member ID(s) provided",
					})
					return
				}

				if result := tx.Model(&project).Association("Members").Append(members); result != nil {
					tx.Rollback()
					fmt.Println("Error associating members with project:", result)
					c.JSON(500, gin.H{
						"error": "internal server error - associating members",
					})
					return
				}
			}

			tx.Commit()

			c.JSON(200, gin.H{
				"message": "success",
				"id":      project.Id,
			})
		})

		project.PATCH(":id", func(c *gin.Context) {
			session := sessions.Default(c)
			db := database.GetDB(c)

			email, ok := session.Get("Email").(string)
			if !ok {
				c.JSON(401, gin.H{"error": "Not logged in"})
				return
			}

			currentUser := models.User{}
			if result := db.Model(&models.User{}).Where(&models.User{Email: email}).First(&currentUser); result.Error != nil {
				c.JSON(500, gin.H{"error": "internal server error - failed to get user"})
				return
			}

			if currentUser.Role < models.UserRoleUploader {
				c.JSON(403, gin.H{"error": "Permission denied."})
				return
			}

			idStr := c.Param("id")
			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid project ID"})
				return
			}

			project := models.Project{}

			if result := db.Model(&models.Project{}).Where(&models.Project{Id: uint(id)}).First(&project); result.Error != nil {
				c.JSON(404, gin.H{"error": "Project not found"})
				return
			}

			if project.UserId != currentUser.Id && currentUser.Role < models.UserRoleAdmin {
				c.JSON(403, gin.H{"error": "Permission denied. Only project uploader or admin can update this project."})
				return
			}

			var body struct {
				Title       string    `json:"title"`
				Thumbnail   string    `json:"thumbnail"`
				Tags        []string  `json:"tags"`
				Description string    `json:"description"`
				Images      []string  `json:"images"`
				StartDate   time.Time `json:"startDate"`
				EndDate     time.Time `json:"endDate"`
				Members     []uint    `json:"members"`
			}

			if err := c.ShouldBindJSON(&body); err != nil {
				fmt.Println("Error binding JSON for project PUT:", err)
				c.JSON(400, gin.H{"error": fmt.Sprintf("Bad request: %v", err)})
				return
			}

			tags, err := util.ToDataTypeJSON(body.Tags)
			if err != nil {
				c.JSON(500, gin.H{"error": "internal server error - marshalling tags"})
				return
			}

			images, err := util.ToDataTypeJSON(body.Images)
			if err != nil {
				c.JSON(500, gin.H{"error": "internal server error - marshalling images"})
				return
			}

			tx := db.Begin()
			if tx.Error != nil {
				fmt.Println("Error starting transaction for project PUT:", tx.Error)
				c.JSON(500, gin.H{"error": "internal server error - transaction start"})
				return
			}

			if len(body.Members) > 0 {
				var members []models.User
				if result := tx.Find(&members, body.Members); result.Error != nil {
					tx.Rollback()
					fmt.Println("Error finding members for project PUT:", result.Error)
					c.JSON(500, gin.H{"error": "internal server error - finding members for update"})
					return
				}

				if len(members) != len(body.Members) {
					tx.Rollback()
					c.JSON(400, gin.H{"error": "Invalid member ID(s) provided for update"})
					return
				}

				if result := tx.Model(&project).Association("Members").Replace(members); result != nil {
					tx.Rollback()
					fmt.Println("Error replacing members for project PUT:", result)
					c.JSON(500, gin.H{"error": "internal server error - replacing members"})
					return
				}
			}
			if body.Title != "" {
				project.Title = body.Title
			}
			if body.Thumbnail != "" {
				project.Thumbnail = body.Thumbnail
			}
			if body.Tags != nil {
				project.Tags = tags
			}
			if body.Description != "" {
				project.Description = body.Description
			}
			if body.Images != nil {
				project.Images = images
			}
			if !body.StartDate.IsZero() {
				project.StartDate = body.StartDate
			}
			if !body.EndDate.IsZero() {
				project.EndDate = body.EndDate
			}

			if result := tx.Model(&models.Project{}).Where(&models.Project{Id: project.Id}).Updates(project); result.Error != nil {
				tx.Rollback()
				fmt.Println("Error saving project for PATCH:", result.Error)
				c.JSON(500, gin.H{"error": "internal server error - saving project"})
				return
			}

			if result := tx.Commit(); result.Error != nil {
				c.JSON(500, gin.H{"error": "internal server error - commit transaction"})
				return
			}

			c.JSON(200, gin.H{"message": "success"})
		})

		project.GET(":id", func(c *gin.Context) {
			db := database.GetDB(c)

			idStr := c.Param("id")
			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid project ID"})
				return
			}

			project := models.Project{}

			if result := db.Model(&models.Project{}).Preload("Members", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name", "profile_photo", "description", "avatar", "major")
			}).Where(&models.Project{Id: uint(id)}).First(&project); result.Error != nil {
				c.JSON(404, gin.H{"error": "Project not found"})
				return
			}

			//c.JSON(200, project)
			c.Render(200, util.JsonL(project))
		})

		project.DELETE(":id", func(c *gin.Context) {
			session := sessions.Default(c)
			db := database.GetDB(c)

			email, ok := session.Get("Email").(string)
			if !ok {
				c.JSON(401, gin.H{"error": "Not logged in"})
				return
			}

			currentUser := models.User{}
			if result := db.Model(&models.User{}).Where(&models.User{Email: email}).First(&currentUser); result.Error != nil {
				c.JSON(500, gin.H{"error": "internal server error - failed to get user"})
				return
			}

			if currentUser.Role < models.UserRoleUploader {
				c.JSON(403, gin.H{"error": "Permission denied. Only uploaders or admins can delete projects."})
				return
			}

			idStr := c.Param("id")
			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid project ID"})
				return
			}

			project := models.Project{}

			if result := db.Where(&models.Project{Id: uint(id)}).First(&project); result.Error != nil {
				c.JSON(404, gin.H{"error": "Project not found"})
				return
			}

			if project.UserId != currentUser.Id && currentUser.Role < models.UserRoleAdmin {
				c.JSON(403, gin.H{"error": "Permission denied. Only project uploader or admin can delete this project."})
				return
			}

			if result := db.Delete(&project); result.Error != nil {
				fmt.Println("Error deleting project:", result.Error)
				c.JSON(500, gin.H{"error": "internal server error - deleting project"})
				return
			}

			c.JSON(200, gin.H{"message": "success"})
		})

		project.GET("/all", func(c *gin.Context) {
			db := database.GetDB(c)

			var projects []struct {
				Id          uint      `json:"id"`
				Title       string    `json:"title"`
				Thumbnail   string    `json:"thumbnail"`
				Tags        string    `json:"tags"`
				Description string    `json:"description"`
				StartDate   time.Time `json:"startDate"`
				EndDate     time.Time `json:"endDate"`
			}

			if result := db.Model(&models.Project{}).Select("id", "title", "thumbnail", "tags", "description", "start_date", "end_date").Find(&projects); result.Error != nil {
				fmt.Println("Error getting all projects:", result.Error)
				c.JSON(500, gin.H{"error": "internal server error"})
				return
			}

			c.JSON(200, projects)
		})
	})
}
