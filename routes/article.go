package routes

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"ntpu_gdg.org/blog/database"
	"ntpu_gdg.org/blog/models"
	"ntpu_gdg.org/blog/routerRegister"
)

func init() {
	routerRegister.Register = append(routerRegister.Register, func(rg *gin.RouterGroup) {
		article := rg.Group("/article")

		article.POST("", func(c *gin.Context) {
			session := sessions.Default(c)
			db := database.GetDB(c)

			email, ok := session.Get("Email").(string)

			if !ok {
				c.JSON(401, gin.H{
					"error": "Not logged in",
				})
				return
			}

			var role int

			if result := db.Model(&models.User{}).Where(&models.User{Email: email}).Select("role").First(&role); result.Error != nil {
				fmt.Println(result.Error)
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}

			if role < models.UserRoleUploader {
				c.JSON(403, gin.H{
					"error": "沒有權限發布文章 Permission denied",
				})
				return
			}

			var body struct {
				Title      string    `json:"title"`
				Topic      int       `json:"topic"`
				Author     string    `json:"author"`
				AuthorInfo string    `json:"authorInfo"`
				Time       time.Time `json:"time"`
				Content    string    `json:"content"`
				Tags       []string  `json:"tags"`
			}

			if err := c.ShouldBindJSON(&body); err != nil {
				fmt.Println(err)
				c.JSON(400, gin.H{
					"error": "Bad request",
				})
				return
			}

			tags := "[" + strings.Join(body.Tags, "\",\"") + "]"

			article := models.Article{
				Title:      body.Title,
				Topic:      body.Topic,
				Author:     body.Author,
				AuthorInfo: body.AuthorInfo,
				Time:       body.Time,
				Content:    body.Content,
				Tags:       tags,
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

		article.GET("/all", func(c *gin.Context) {
			db := database.GetDB(c)

			var articles []struct {
				Id         uint
				Title      string
				Author     string
				AuthorInfo string
				Time       time.Time
				Tags       []string
			}

			if result := db.Model(&models.Article{}).Select("id", "title", "author", "author_info", "time", "tags").Find(&articles); result.Error != nil {
				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}

			c.JSON(200, articles)
		})

	})
}
