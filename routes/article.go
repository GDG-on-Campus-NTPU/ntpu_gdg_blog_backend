package routes

import (
	"fmt"
	"strconv"
	"time"

	"blog/database"
	"blog/models"
	"blog/routerRegister"
	"blog/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func init() {
	routerRegister.Register = append(routerRegister.Register, func(rg *gin.RouterGroup) {
		article := rg.Group("/article")

		article.POST("", func(c *gin.Context) {
			session := sessions.Default(c)
			db := database.GetDB(c)

			email, ok := session.Get("Email").(string)

			currentUser := models.User{}

			if result := db.Model(&models.User{}).Where(&models.User{Email: email}).First(&currentUser); result.Error != nil {
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}
			if currentUser.Role < models.UserRoleUploader {
				c.JSON(403, gin.H{
					"error": "Permission denied 沒有權限發布文章 ",
				})
				return
			}

			if !ok {
				c.JSON(401, gin.H{
					"error": "Not logged in",
				})
				return
			}

			var body struct {
				Title       string    `json:"title"`
				Topic       int       `json:"topic"`
				Author      string    `json:"author"`
				AuthorInfo  *string   `json:"authorInfo"`
				Time        time.Time `json:"time"`
				Content     string    `json:"content"`
				Tags        []string  `json:"tags"`
				Type        int       `json:"type"`
				Description string    `json:"description"`
			}

			if err := c.ShouldBindJSON(&body); err != nil {
				fmt.Println(err)
				c.JSON(400, gin.H{
					"error": "Bad request",
				})
				return
			}

			tags, err := util.ToDataTypeJSON(body.Tags)

			if err != nil {
				c.JSON(400, gin.H{
					"error": "Bad request",
				})
				return
			}

			authorInfo := ""
			if body.AuthorInfo != nil {
				authorInfo = *body.AuthorInfo
			}

			article := models.Article{
				Title:       body.Title,
				Topic:       body.Topic,
				Author:      body.Author,
				AuthorInfo:  authorInfo,
				Time:        body.Time,
				Content:     body.Content,
				Tags:        tags,
				UserId:      currentUser.Id,
				Type:        body.Type,
				Description: body.Description,
			}

			result := db.Model(&models.Article{}).Create(&article)

			if result.Error != nil {
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}

			c.JSON(200, gin.H{
				"message": "success",
				"id":      article.Id,
			})
		})

		article.PUT(":id", func(c *gin.Context) {
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
				c.JSON(401, gin.H{
					"error": "Invalid article ID",
				})
				return
			}
			article := models.Article{}
			if result := db.Where(&models.Article{Id: uint(id)}).First(&article); result.Error != nil {
				c.JSON(404, gin.H{
					"error": "not found article",
				})
				return
			}
			if currentUser.Role < models.UserRoleAdmin {
				if article.UserId != currentUser.Id {
					c.JSON(403, gin.H{
						"error": "Permission denied; Only admin or article owner can update this article",
					})
					return
				}
			}

			var body struct {
				Title       string    `json:"title"`
				Topic       int       `json:"topic"`
				Author      string    `json:"author"`
				AuthorInfo  *string   `json:"authorInfo"`
				Time        time.Time `json:"time"`
				Content     string    `json:"content"`
				Tags        []string  `json:"tags"`
				Type        int       `json:"type"`
				Description string    `json:"description"`
			}

			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(400, gin.H{
					"error": "Bad request",
				})
				fmt.Println(err)
				return
			}

			tags, err := util.ToDataTypeJSON(body.Tags)

			if err != nil {
				c.JSON(400, gin.H{
					"error": "Bad request",
				})
				return
			}

			authorInfo := ""
			if body.AuthorInfo != nil {
				authorInfo = *body.AuthorInfo
			}

			newArticle := models.Article{
				Title:       body.Title,
				Topic:       body.Topic,
				Author:      body.Author,
				AuthorInfo:  authorInfo,
				Time:        body.Time,
				Content:     body.Content,
				Tags:        tags,
				Type:        body.Type,
				Description: body.Description,
			}

			if result := db.Model(&models.Article{}).Where(&models.Article{Id: uint(id)}).Updates(&newArticle); result.Error != nil {
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "success",
			})
		})

		article.GET(":id", func(c *gin.Context) {
			idStr := c.Param("id")

			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				c.JSON(401, gin.H{
					"error": "Invalid article ID",
				})
				return
			}

			db := database.GetDB(c)

			article := models.Article{}

			if result := db.Where(&models.Article{Id: uint(id)}).First(&article); result.Error != nil {
				c.JSON(404, gin.H{
					"error": "not found article",
				})
				return
			}

			c.JSON(200, article)
		})

		article.DELETE(":id", func(c *gin.Context) {
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
				c.JSON(401, gin.H{
					"error": "Invalid article ID",
				})
				return
			}
			article := models.Article{}
			if result := db.Where(&models.Article{Id: uint(id)}).First(&article); result.Error != nil {
				c.JSON(404, gin.H{
					"error": "not found article",
				})
				return
			}

			if currentUser.Role < models.UserRoleAdmin {
				if article.UserId != currentUser.Id {
					c.JSON(403, gin.H{
						"error": "Permission denied; Only admin or article owner can delete this article",
					})
					return
				}
			}

			if result := db.Delete(&article); result.Error != nil {
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "success",
			})
		})

		article.GET("/all", func(c *gin.Context) {
			db := database.GetDB(c)

			var articles []struct {
				Id          uint
				Title       string
				Author      string
				AuthorInfo  string
				Time        time.Time
				Tags        []string
				Type        int
				Description string
			}

			if result := db.Model(&models.Article{}).Select("id", "title", "author", "author_info", "time", "tags", "type", "description").Find(&articles); result.Error != nil {
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}

			c.JSON(200, articles)
		})

	})
}
