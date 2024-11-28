package routes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/dchest/uniuri"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"ntpu_gdg.org/blog/auth"
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

			state := c.Query("state")

			if state != session.Get("OauthState") {
				//c.AbortWithError(http.StatusUnauthorized, errors.New("invalid csrf token"))
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
				//c.AbortWithError(http.StatusInternalServerError, err)
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			client := auth.GoogleOauthConfig.Client(context.Background(), token)

			response, err := client.Get("https://people.googleapis.com/v1/people/me?personFields=names,emailAddresses")

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

			err = json.Unmarshal(responseData, &userInfo)

			if err != nil {
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

			name, ok := userInfo["names"].([]any)[0].(map[string]any)["displayName"].(string)

			if !ok {
				c.JSON(500, gin.H{
					"error": "userInfo from google parse failed : field email error",
				})
				return
			}

			session.Set("Name", name)
			session.Set("Email", email)

			redirect := session.Get("Redirect").(string)

			if redirect != "" && redirect[0] != '/' {
				redirect = "/"
			}

			session.Delete("Redirect")
			session.Save()

			c.Redirect(http.StatusFound, redirect)
		})
	})
}
