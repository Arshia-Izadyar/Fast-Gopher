package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleConfig struct {
	GoogleLoginConfig oauth2.Config
}

var AppConfig GoogleConfig

func LoadGoogleConfig() {

	AppConfig.GoogleLoginConfig = oauth2.Config{
		RedirectURL:  "http://127.0.0.1:4000/auth/callback/google",
		ClientID:     "669774672842-eoia8gmm2bt3e5a98110b4j4hmnaq9n0.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}
}
