package auth

import (
	"fmt"

	"ntpu_gdg.org/blog/env"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func newGoogleOauthConfig() *oauth2.Config {
	if env.Getenv("GOOGLE_CLIENT_ID") == "" || env.Getenv("GOOGLE_CLIENT_SECRET") == "" {
		fmt.Println("missing GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET")
	}
	return &oauth2.Config{
		ClientID:     env.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: env.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  env.Getenv("BASE_URL") + "/api/login/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

var GoogleOauthConfig = newGoogleOauthConfig()