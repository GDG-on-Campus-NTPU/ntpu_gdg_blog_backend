package routes

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/dchest/uniuri"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"ntpu_gdg.org/blog/auth"
	"ntpu_gdg.org/blog/env"
)

func AddLoginRoutes(rg *gin.RouterGroup) {
	login := rg.Group("/login")

	login.GET("", func(c *gin.Context) {
		session := sessions.Default(c)
		defer session.Save()

		redirect := c.Query("redirect")

		session.Set("Redirect", redirect)

		state := uniuri.NewLen(32)

		session.Set("OauthState", state)

		c.Redirect(http.StatusFound, auth.GoogleOauthConfig.AuthCodeURL(state))
	})

	login.GET("/check", func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("Email") == nil {
			c.JSON(401, gin.H{
				"isLoggedIn": false,
			})
		}

		c.JSON(200, gin.H{
			"isLoggedIn": true,
		})
	})

	login.GET("/google/callback", func(c *gin.Context) {
		session := sessions.Default(c)
		defer session.Save()

		state := c.Query("state")
		if state != session.Get("OauthState") {
			c.AbortWithError(http.StatusUnauthorized, errors.New("invalid csrf token"))
			return
		}

		session.Delete("OauthState")

		code := c.Query("code")

		token, err := auth.GoogleOauthConfig.Exchange(context.Background(), code)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		client := auth.GoogleOauthConfig.Client(context.Background(), token)

		response, err := client.Get("https://people.googleapis.com/v1/people/me?personFields=names,emailAddresses")

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		responseData, err := io.ReadAll(response.Body)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var userInfo map[string]interface{}

		err = json.Unmarshal(responseData, &userInfo)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		email, ok := userInfo["emailAddresses"].([]interface{})[0].(map[string]interface{})["value"].(string)

		if !ok {
			c.AbortWithStatus(500)
			return
		}

		name, ok := userInfo["names"].([]interface{})[0].(map[string]interface{})["displayName"].(string)

		if !ok {
			c.AbortWithStatus(500)
			return
		}

		session.Set("Name", name)
		session.Set("Email", email)

		redirect := session.Get("Redirect").(string)

		if redirect != "" && redirect[0] != '/' {
			redirect = env.Getenv("BASE_URL")
		}

		session.Delete("Redirect")

		c.Redirect(http.StatusFound, redirect)
	})
}
