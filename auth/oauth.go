package auth

import (
	"ntpu_gdg.org/blog/env"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GoogleOauthConfig(baseUrl string) *oauth2.Config {
	if env.Getenv("GOOGLE_CLIENT_ID") == "" || env.Getenv("GOOGLE_CLIENT_SECRET") == "" {
		panic("missing GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET for google oauth")
	}
	return &oauth2.Config{
		ClientID:     env.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: env.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  baseUrl + "/api/login/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
