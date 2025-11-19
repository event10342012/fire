package ioc

import (
	auth "fire/internal/service/oauth2/google"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func InitGoogleService() auth.AuthService {
	googleOauthConfig := &oauth2.Config{
		ClientID:     "223534559756-idf2u1r7uglqokpsedlevld2djk7qmeh.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/oauth2/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return auth.AuthService{Config: googleOauthConfig}
}
