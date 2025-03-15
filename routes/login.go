package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dchest/uniuri"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"ntpu_gdg.org/blog/auth"
	"ntpu_gdg.org/blog/database"
	"ntpu_gdg.org/blog/models"
	"ntpu_gdg.org/blog/routerRegister"
)

func init() {
	routerRegister.Register = append(routerRegister.Register, func(rg *gin.RouterGroup) {
		login := rg.Group("/login")

		login.GET("", func(c *gin.Context) {
			session := sessions.Default(c)

			redirect := c.Query("redirect")

			session.Set("Redirect", redirect)

			state := uniuri.NewLen(32)

			session.Set("OauthState", state)

			session.Save()

			c.Redirect(http.StatusFound, auth.GoogleOauthConfig.AuthCodeURL(state))
		})

		login.GET("/check", func(c *gin.Context) {
			session := sessions.Default(c)

			c.JSON(200, gin.H{
				"isLoggedIn": session.Get("Email") != nil,
			})
		})

		login.GET("/google/callback", func(c *gin.Context) {
			session := sessions.Default(c)

			if c.Query("state") != session.Get("OauthState") {
				fmt.Println("Invalid csrf token from IP:", c.ClientIP(), "=>", c.Query("state"), "!=", session.Get("OauthState"))
				c.JSON(401, gin.H{
					"error": "invalid csrf token",
				})
				return
			}

			session.Delete("OauthState")
			session.Save()

			code := c.Query("code")

			token, err := auth.GoogleOauthConfig.Exchange(context.Background(), code)

			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			client := auth.GoogleOauthConfig.Client(context.Background(), token)

			response, err := client.Get("https://people.googleapis.com/v1/people/me?personFields=emailAddresses")

			if err != nil {
				c.JSON(500, gin.H{
					"error": "fail to get userInfo from google",
				})
				return
			}

			responseData, err := io.ReadAll(response.Body)

			if err != nil {
				c.JSON(500, gin.H{
					"error": "fail to read response from google api",
				})
				return
			}

			var userInfo map[string]any

			if err := json.Unmarshal(responseData, &userInfo); err != nil {
				c.JSON(500, gin.H{
					"error": "userInfo from google parse failed : Invaild Json",
				})
				return
			}

			email, ok := userInfo["emailAddresses"].([]any)[0].(map[string]any)["value"].(string)

			if !ok {
				c.JSON(500, gin.H{
					"error": "userInfo from google parse failed : field email error",
				})
				return
			}

			db := database.GetDB(c)

			var user models.User

			if result := db.Model(&models.User{}).Where(&models.User{Email: email}).FirstOrCreate(&user); result.Error != nil {
				fmt.Println(result.Error)

				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}

			if result := db.Model(&models.User{}).Where(&models.User{Email: email}).Update("last_login", time.Now()); result.Error != nil {
				fmt.Println(result.Error)

				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				return
			}

			session.Set("Email", email)

			redirect := session.Get("Redirect").(string)

			if redirect == "" || redirect[0] != '/' {
				redirect = "/"
			}

			session.Delete("Redirect")
			session.Save()

			fmt.Println("Redirect to => ", redirect)

			c.Redirect(http.StatusFound, redirect)
		})
	})
}
